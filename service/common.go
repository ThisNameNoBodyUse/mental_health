package service

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/xuri/excelize/v2"
	"mental/config"
	"mental/dao"
	"mental/models"
	"mental/utils"
	"strconv"
	"strings"
	"time"
)

type FileService struct {
	FileId string `json:"file_id"`
}

// CheckFileIsExit 检查文件是否已经存在
func (fileService *FileService) CheckFileIsExist(file_id string) (string, bool) {
	fileDao := dao.NewFileDao(config.DB)
	file, err := fileDao.CheckFileIfExist(file_id)
	if err != nil {
		return "", false // 文件不存在
	}
	return file.Path, true // 文件存在，返回文件路径和true
}

// SaveFile 使用工具类上传文件，并将文件信息保存到数据库中
func (fileService *FileService) SaveFile(filePath string) (string, error) {
	// 先计算文件 MD5
	fileMD5, err := utils.GetFileMD5(filePath)
	if err != nil {
		return "", err
	}

	// 上传到 MinIO，返回路径
	url, err := utils.UploadFileToMinio(filePath)
	if err != nil {
		return "", err
	}

	// 这里存数据库时，保存路径
	fileDao := dao.NewFileDao(config.DB)
	err = fileDao.SaveFile(fileMD5, url)
	if err != nil {
		return "", err
	}

	// 返回完整 URL
	return url, nil
}

// ImportFromFileId 根据 file_id 获取路径并导入 Excel 数据
func (s *FileService) ImportFromFileId() (int, []string) {
	fileId := s.FileId

	// 1. 根据 fileId 查数据库拿到 MinIO 对象路径（完整 URL）
	objectPath, exists := s.CheckFileIsExist(fileId)
	if !exists {
		return 0, []string{fmt.Sprintf("文件ID %s 不存在", fileId)}
	}

	// 判断文件是否已解析
	fileDao := dao.NewFileDao(config.DB)
	hasAnalyzed, err := fileDao.IsFileAnalyzed(fileId)
	if err != nil {
		return 0, []string{fmt.Sprintf("判断文件状态错误！")}
	}
	if hasAnalyzed { // 如果解析过，防止重复解析
		return 0, []string{fmt.Sprintf("该文件已经解析过！")}
	}

	ctx := context.Background()

	// 2. 处理 objectPath，去除 URL 前缀和桶名，得到 MinIO 的 objectKey
	endpoint := config.MinioSettings.Endpoint // "8.130.77.225:9000"
	prefix := "http://" + endpoint + "/"      // "http://8.130.77.225:9000/"
	bucket := config.MinioSettings.Bucket     // "mental"

	// 去掉 URL 前缀
	objectKey := strings.TrimPrefix(objectPath, prefix) // "mental/2025-05-23/xxxx.xlsx"

	// 去掉桶名和斜杠
	objectKey = strings.TrimPrefix(objectKey, bucket+"/") // "2025-05-23/xxxx.xlsx"

	// 3. 从 MinIO 下载文件流
	object, err := config.MinioClient.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return 0, []string{fmt.Sprintf("从MinIO获取文件失败: %v", err)}
	}
	defer object.Close()

	// 4. 读取 Excel 文件
	f, err := excelize.OpenReader(object)
	if err != nil {
		return 0, []string{fmt.Sprintf("打开Excel失败: %v", err)}
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return 0, []string{fmt.Sprintf("读取工作表失败: %v", err)}
	}

	sclDao := dao.NewSCLDao(config.DB)
	successCount := 0
	var errorRows []string

	for i, row := range rows[1:] {
		rowNum := i + 2
		if len(row) < 15 {
			errorRows = append(errorRows, fmt.Sprintf("第 %d 行数据列数不足（当前列数：%d）", rowNum, len(row)))
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					errorRows = append(errorRows, fmt.Sprintf("第 %d 行解析异常: %v", rowNum, r))
				}
			}()

			age, err := strconv.Atoi(row[3])
			if err != nil {
				errorRows = append(errorRows, fmt.Sprintf("第 %d 行年龄格式错误: %v", rowNum, err))
				return
			}

			gender, err := strconv.Atoi(row[2])
			if err != nil {
				errorRows = append(errorRows, fmt.Sprintf("第 %d 行性别格式错误: %v", rowNum, err))
				return
			}

			testDate, err := time.Parse("2006-01-02", row[4])
			if err != nil {
				errorRows = append(errorRows, fmt.Sprintf("第 %d 行测评日期格式错误: %v", rowNum, err))
				return
			}

			scl := &models.SCL{
				StudentID:     parseInt64Ptr(row[0]),
				Name:          row[1],
				Gender:        gender,
				Age:           age,
				TestDate:      models.CustomTime(testDate),
				Somatization:  parseFloat(row[5]),
				Obsession:     parseFloat(row[6]),
				Interpersonal: parseFloat(row[7]),
				Depression:    parseFloat(row[8]),
				Anxiety:       parseFloat(row[9]),
				Hostility:     parseFloat(row[10]),
				Phobia:        parseFloat(row[11]),
				Paranoia:      parseFloat(row[12]),
				Psychoticism:  parseFloat(row[13]),
				Other:         parseFloat(row[14]),
			}

			if err := sclDao.Save(scl); err != nil {
				errorRows = append(errorRows, fmt.Sprintf("第 %d 行插入失败: %v", rowNum, err))
				return
			}

			successCount++
		}()
	}
	fileDao.UpdateStatusAnalyzed(fileId) // 将文件设置为已解析
	return successCount, errorRows
}

// parseFloat 安全解析 float32
func parseFloat(s string) float32 {
	if s == "" {
		return 0.0
	}
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0.0
	}
	return float32(val)
}

// parseInt64Ptr 安全解析为 *int64，如果为空或非法返回 nil
func parseInt64Ptr(s string) *int64 {
	if s == "" {
		return nil
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &val
}
