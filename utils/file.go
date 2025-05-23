package utils

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"io"
	"mental/config"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

// GetFileMD5 用于获取文件的 MD5 值
func GetFileMD5(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()
	// 创建 MD5 hash 计算器
	hash := md5.New()
	// 计算文件的 MD5 值
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("计算文件 MD5 时出错: %v", err)
	}
	// 返回 MD5 值的十六进制表示
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// UploadFile 上传文件并保存至 storage 文件夹，命名规则为当前日期的目录加上 MD5 值作为文件名，扩展名不变
func UploadFile(filePath string) (string, error) {
	// 获取当前日期，格式为：年-月-日
	currentDate := time.Now().Format("2006-01-02")

	// 获取文件的扩展名
	ext := filepath.Ext(filePath)

	// 获取文件的 MD5 值作为文件名
	md5Value, err := GetFileMD5(filePath)
	if err != nil {
		return "", fmt.Errorf("获取文件 MD5 时出错: %v", err)
	}
	// 生成存储路径
	storageDir := filepath.Join("storage", currentDate)
	fileName := fmt.Sprintf("%s%s", md5Value, ext)
	savePath := filepath.Join(storageDir, fileName)
	// 检查目录是否存在，不存在则创建
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		err := os.MkdirAll(storageDir, 0755)
		if err != nil {
			return "", fmt.Errorf("创建目录失败: %v", err)
		}
	}
	// 打开上传的文件
	srcFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer destFile.Close()
	// 将文件内容复制到目标文件
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return "", fmt.Errorf("文件复制失败: %v", err)
	}
	// 返回相对路径
	return savePath, nil
}

// GenerateFileURL 从 gin.Context 获取主机信息，并生成完整的文件 URL
func GenerateFileURL(c *gin.Context, path string) string {
	// 获取当前请求的协议 (http 或 https)
	protocol := "http" // 默认使用 http
	if c.Request.TLS != nil {
		protocol = "https" // 如果请求使用了 TLS，则使用 https
	}
	// 获取当前请求的主机地址（例如 localhost:8080）
	serverHost := c.Request.Host
	// 拼接协议、主机和路径，返回完整的文件 URL
	return protocol + "://" + serverHost + path[1:]
}

// UploadFileToMinio 上传文件到 MinIO，并返回文件在桶中的路径
func UploadFileToMinio(filePath string) (string, error) {
	// 获取当前日期和扩展名
	currentDate := time.Now().Format("2006-01-02")
	ext := filepath.Ext(filePath)

	// 获取文件的 MD5 值作为文件名
	md5Value, err := GetFileMD5(filePath)
	if err != nil {
		return "", fmt.Errorf("获取文件 MD5 时出错: %v", err)
	}

	// 构造对象名（类似路径）：例如 2025-05-23/abc123.jpg
	objectName := path.Join(currentDate, md5Value+ext)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileStat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 上传到 MinIO
	_, err = config.MinioClient.PutObject(context.Background(), config.MinioSettings.Bucket, objectName, file, fileStat.Size(), minio.PutObjectOptions{
		ContentType: GetContentType(file, filePath),
	})
	if err != nil {
		return "", fmt.Errorf("上传到 MinIO 失败: %v", err)
	}

	// 返回文件路径
	return GenerateMinioFileURL(objectName), nil
}

// GenerateMinioFileURL 生成Minio的文件路径
func GenerateMinioFileURL(objectPath string) string {
	scheme := "http"
	if config.MinioSettings.Secure {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, config.MinioSettings.Endpoint, config.MinioSettings.Bucket, objectPath)
}

// GetContentType 根据文件内容或扩展名判断 Content-Type
func GetContentType(file *os.File, filePath string) string {
	// 读取前 512 字节来尝试判断 MIME 类型
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err == nil && n > 0 {
		// 重置文件指针，避免影响后续读取
		file.Seek(0, io.SeekStart)
		return http.DetectContentType(buffer[:n])
	}

	// 如果读取失败或文件太小，则根据扩展名判断
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}
