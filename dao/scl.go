package dao

import (
	"gorm.io/gorm"
	"mental/models"
)

// SCLDao 负责操作 scl 表（心理测评数据）
type SCLDao struct {
	DB *gorm.DB
}

// NewSCLDao 创建 SCLDao 实例
func NewSCLDao(db *gorm.DB) *SCLDao {
	return &SCLDao{DB: db}
}

// Save 插入一条 SCL 记录 *
func (dao *SCLDao) Save(scl *models.SCL) error {
	result := dao.DB.Create(scl)
	return result.Error
}

// FindByID 根据 ID 查询一条记录
func (dao *SCLDao) FindByID(id int64) (*models.SCL, error) {
	var scl models.SCL
	if err := dao.DB.First(&scl, id).Error; err != nil {
		return nil, err
	}
	return &scl, nil
}

// List 列出所有记录（可分页后续扩展）
func (dao *SCLDao) List() ([]models.SCL, error) {
	var list []models.SCL
	if err := dao.DB.Order("test_date desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// DeleteByID 根据 ID 删除
func (dao *SCLDao) DeleteByID(id int64) error {
	return dao.DB.Delete(&models.SCL{}, id).Error
}

// SelectAllByUserId 根据用户id，和时间排序，查找该用户的历史测量数据 TODO 2025-5-22
func (dao *SCLDao) SelectAllByUserId(userId int64) ([]models.SCL, error) {
	var list []models.SCL
	if err := dao.DB.Where("student_id = ?", userId).Order("test_date desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateByID 根据 ID 更新指定字段
func (dao *SCLDao) UpdateByID(id int64, updated map[string]interface{}) error {
	return dao.DB.Model(&models.SCL{}).Where("id = ?", id).Updates(updated).Error
}
