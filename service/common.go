package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"mental/config"
	"mental/dao"
	"mental/models"
	"mental/utils"
	"path/filepath"
	"strconv"
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
func (FileService *FileService) SaveFile(filePath string) (string, error) {
	// 计算文件的 MD5 值
	fileMD5, _ := utils.GetFileMD5(filePath)

	// 保存文件到磁盘，并获取文件路径
	path, err := utils.UploadFile(filePath)
	if err != nil {
		return "", err
	}

	// 使用 filepath.ToSlash 将路径转换为适用于 URL 或相对路径的格式（统一使用 /）
	// 并确保路径以 ./ 开头，表示相对路径
	relativePath := "./" + filepath.ToSlash(path)

	// 保存文件到数据库
	fileDao := dao.NewFileDao(config.DB)
	err = fileDao.SaveFile(fileMD5, relativePath)
	if err != nil {
		return "", err
	}

	return relativePath, nil
}

// ImportFromFileId 根据 file_id 获取路径并导入 Excel 数据
func (s *FileService) ImportFromFileId() (int, []string) {
	fileId := s.FileId
	// 1. 查找文件路径
	path, exists := s.CheckFileIsExist(fileId)
	if !exists {
		return 0, []string{fmt.Sprintf("文件ID %s 不存在", fileId)}
	}

	// 2. 打开 Excel 文件
	f, err := excelize.OpenFile(path)
	if err != nil {
		return 0, []string{fmt.Sprintf("打开文件失败: %v", err)}
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return 0, []string{fmt.Sprintf("读取工作表失败: %v", err)}
	}

	sclDao := dao.NewSCLDao(config.DB)
	successCount := 0
	var errorRows []string

	// 3. 解析数据行，跳过表头
	for i, row := range rows[1:] {
		rowNum := i + 2 // Excel 的真实行号
		if len(row) < 15 {
			errorRows = append(errorRows, fmt.Sprintf("第 %d 行数据列数不足（当前列数：%d）", rowNum, len(row)))
			continue
		}

		defer func() {
			if r := recover(); r != nil {
				errorRows = append(errorRows, fmt.Sprintf("第 %d 行解析异常: %v", rowNum, r))
			}
		}()

		age, err := strconv.Atoi(row[3])
		if err != nil {
			errorRows = append(errorRows, fmt.Sprintf("第 %d 行年龄格式错误: %v", rowNum, err))
			continue
		}

		gender, err := strconv.Atoi(row[2])
		if err != nil {
			errorRows = append(errorRows, fmt.Sprintf("第 %d 行性别格式错误: %v", rowNum, err))
			continue
		}

		testDate, err := time.Parse("2006-01-02", row[4])
		if err != nil {
			errorRows = append(errorRows, fmt.Sprintf("第 %d 行测评日期格式错误: %v", rowNum, err))
			continue
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
			continue
		}

		successCount++
	}

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
