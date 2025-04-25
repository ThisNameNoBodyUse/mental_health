package models

import "time"

// User 数据库表user结构体
type User struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Account    string    `json:"account"`
	Password   string    `json:"password"`
	Username   string    `json:"username,omitempty"`
	Email      string    `json:"email,omitempty"`
	Avatar     string    `json:"avatar" gorm:"default:'./storage/default_avatar.jpg'"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoCreateTime"`
}

// TableName 手动指定表名，防止gorm自动转换错误
func (User) TableName() string {
	return "user"
}
