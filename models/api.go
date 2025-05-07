package models

import "time"

type API struct {
	ID           int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Path         string    `json:"path"`
	Method       string    `json:"method"`
	Description  string    `json:"description"`
	PermissionID *int      `json:"permission_id"` // *int 指针表示可为 null
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
}

// TableName 手动指定表名，防止gorm自动转换错误
func (API) TableName() string {
	return "api"
}
