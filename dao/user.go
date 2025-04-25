package dao

import (
	"gorm.io/gorm"
	"mental/models"
)

// UserDao 操作用户表的dao结构体
type UserDao struct {
	*gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// GetUserByAccount 根据账号查找数据库中是否有该用户，用于登录/注册
func (dao *UserDao) GetUserByAccount(account string) (*models.User, error) {
	user := new(models.User)
	res := dao.Where("account = ?", account).First(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

// GetUserById 根据用户Id查找用户基本信息
func (dao *UserDao) GetUserById(id int64) (*models.User, error) {
	user := new(models.User)
	res := dao.Where("id = ?", id).First(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
