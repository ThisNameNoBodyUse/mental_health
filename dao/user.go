package dao

import (
	"errors"
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

	// "记录未找到"
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil // 不返回错误，而是返回 nil，表示用户不存在
	}
	// 其他错误
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
