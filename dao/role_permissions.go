package dao

import (
	"mental/config"
	"mental/models"
)

// GetRolePermissionsFromDB 从数据库查询角色权限列表
func GetRolePermissionsFromDB(roleID string) ([]string, error) {
	var permissions []string

	// 查询角色权限
	err := config.DB.Raw(`
		SELECT p.id
		FROM role_permissions rp
		JOIN permissions p ON rp.permission_id = p.id
		WHERE rp.role_id = ?
	`, roleID).Scan(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetPermissionAPIsFromDB 查找指定的权限id对应的接口列表
func GetPermissionAPIsFromDB(permissionID string) ([]string, error) {
	var apis []models.API
	err := config.DB.Where("permission_id = ?", permissionID).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	var apiSet []string
	for _, api := range apis {
		apiSet = append(apiSet, api.Path+":"+api.Method)
	}
	return apiSet, nil
}
