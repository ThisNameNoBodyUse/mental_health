package dao

import (
	"errors"
	"gorm.io/gorm"
	"mental/models"
	"time"
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

// UpdateAvatar 根据用户id修改用户头像
func (dao *UserDao) UpdateAvatar(userId int64, url string) error {
	// Updates(updates) 传入map，批量更新字段
	res := dao.DB.Model(models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"avatar":      url,
		"update_time": time.Now(),
	})
	return res.Error
}

// UpdateUsernameOrEmail 根据用户id，修改用户名/邮箱
func (dao *UserDao) UpdateUsernameOrEmail(id int64, username string, email string) error {
	updates := map[string]interface{}{
		"update_time": time.Now(), // 一定要放进去，不然为空时不会更新
	}
	if username != "" {
		updates["username"] = username
	}
	if email != "" {
		updates["email"] = email
	}

	res := dao.DB.Model(models.User{}).Where("id = ?", id).Updates(updates)
	return res.Error
}

// UpdatePassword 更新用户密码
func (dao *UserDao) UpdatePassword(userId int64, hashedPwd string) error {
	return dao.DB.Model(models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"password":    hashedPwd,
		"update_time": time.Now(),
	}).Error
}
