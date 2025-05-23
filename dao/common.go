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

// UpdateStatusAnalyzed 根据文件id，将文件状态设置为已解析
func (dao *FileDao) UpdateStatusAnalyzed(fileID string) error {
	return dao.DB.Model(&models.File{}).
		Where("file_id = ?", fileID).
		Update("status", 1).Error
}

// IsFileAnalyzed 根据文件id判断是否已经被解析
func (dao *FileDao) IsFileAnalyzed(fileID string) (bool, error) {
	var file models.File
	err := dao.DB.Select("status").Where("file_id = ?", fileID).First(&file).Error
	if err != nil {
		return false, err // 查询出错
	}
	return file.Status == 1, nil
}
