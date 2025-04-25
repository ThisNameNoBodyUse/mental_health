package dao

import (
	"gorm.io/gorm"
	"mental/models"
)

type FileDao struct {
	*gorm.DB
}

// NewFileDao 依赖注入
func NewFileDao(db *gorm.DB) *FileDao {
	return &FileDao{db}
}

// CheckFileIfExit 根据文件id查找数据库是否有对应记录
func (dao *FileDao) CheckFileIfExist(file_id string) (*models.File, error) {
	file := new(models.File)
	res := dao.Where("file_id = ?", file_id).First(file)
	return file, res.Error
}

// SaveFile 根据文件id，文件相对路径保存到数据库
func (dao *FileDao) SaveFile(file_id string, path string) error {
	file := new(models.File)
	file.FileId = file_id
	file.Path = path
	// 保存文件
	res := dao.Save(file)
	return res.Error
}
