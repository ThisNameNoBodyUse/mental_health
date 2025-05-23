package models

import "time"

// File 文件表结构体
type File struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	FileId     string    `json:"file_id" gorm:"column:file_id;"`
	Path       string    `json:"path" gorm:"column:path;"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoCreateTime"`
	Status     int       `json:"status" gorm:"column:status;"` // 是否已解析
}

func (File) TableName() string {
	return "file"
}
