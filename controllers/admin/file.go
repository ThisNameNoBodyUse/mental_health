package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mental/controllers/common"
	"mental/service"
	"mental/utils"
	"os"
	"path/filepath"
)

type FileController struct {
	common.BaseController
}

// Check 检查文件是否已经上传过
func (con FileController) Check(c *gin.Context) {
	fileId := c.Query("file_id") // 文件的md5值即为文件id
	var fileService service.FileService
	path, isExist := fileService.CheckFileIsExist(fileId)
	// var定义临时结构体，返回路径和是否存在
	var upload struct {
		Path    string
		IsExist bool
	}
	upload.IsExist = isExist
	if isExist {
		// 获取当前服务器的协议和主机地址
		serverHost := c.Request.Host
		upload.Path = "http://" + serverHost + path[1:]
	}
	con.Success(c, upload)
}

// Upload 上传文件接口
func (con FileController) Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file") // "file" 是前端表单中文件上传字段的名字
	if err != nil {
		con.Error(c, nil, fmt.Sprintf("文件上传失败: %v", err))
		return
	}

	// 生成唯一的 UUID 作为文件名
	uniqueFileName := uuid.New().String()

	// 获取文件扩展名
	ext := filepath.Ext(file.Filename)

	// 临时文件保存路径，使用 UUID 作为文件名
	tempFilePath := "./storage/temp/" + uniqueFileName + ext

	// 保存临时文件到 storage/temp 目录
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		con.Error(c, nil, fmt.Sprintf("保存文件失败: %v", err))
		return
	}

	// 获取文件的绝对路径
	absPath, err := filepath.Abs(tempFilePath)
	if err != nil {
		con.Error(c, nil, fmt.Sprintf("获取文件绝对路径失败: %v", err))
		return
	}

	// 使用 defer 确保在函数结束时删除临时文件
	defer func() {
		if err := os.Remove(absPath); err != nil {
			// 如果删除临时文件失败，记录日志但不影响返回结果
			fmt.Printf("删除临时文件失败: %v\n", err)
		}
	}()

	// 计算文件的 MD5 值
	fileMD5, err := utils.GetFileMD5(absPath) // 使用绝对路径计算 MD5
	if err != nil {
		con.Error(c, nil, fmt.Sprintf("计算文件 MD5 时出错: %v", err))
		return
	}

	// 调用 FileService 检查文件是否已经上传过
	var fileService service.FileService
	filePath, exists := fileService.CheckFileIsExist(fileMD5)
	if exists {
		// 如果已经上传过，直接返回文件路径
		con.Success(c, utils.GenerateFileURL(c, filePath))
		return
	}

	// 没有上传过，调用service的文件上传方法进行上传
	path, err := fileService.SaveFile(absPath) // 使用绝对路径进行保存
	if err != nil {
		con.Error(c, nil, "文件保存失败")
		return
	}

	// 文件保存成功
	con.Success(c, utils.GenerateFileURL(c, path))
}
