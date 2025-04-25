package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
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
