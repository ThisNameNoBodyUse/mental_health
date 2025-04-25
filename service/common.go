package service

import (
	"mental/config"
	"mental/dao"
	"mental/utils"
	"path/filepath"
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
