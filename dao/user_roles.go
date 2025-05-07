package dao

import (
	"gorm.io/gorm"
	"mental/models"
)

// UserRoleDao 用户角色关联表 DAO
type UserRoleDao struct {
	DB *gorm.DB
}

// NewUserRoleDao 创建 UserRoleDao
func NewUserRoleDao(db *gorm.DB) *UserRoleDao {
	return &UserRoleDao{DB: db}
}

// Save 插入用户角色关联记录
func (dao *UserRoleDao) Save(userRole *models.UserRole) *gorm.DB {
	return dao.DB.Create(userRole)
}
