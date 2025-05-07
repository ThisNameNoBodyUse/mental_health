package models

// UserRole 用户角色关联表
type UserRole struct {
	UserID int `gorm:"column:user_id"` // 用户ID
	RoleID int `gorm:"column:role_id"` // 角色ID
}

// TableName 返回表名
func (UserRole) TableName() string {
	return "user_roles"
}
