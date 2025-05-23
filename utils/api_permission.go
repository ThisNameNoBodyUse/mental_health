package utils

import (
	"fmt"
	"mental/config"
	"mental/constant"
	"mental/models"
	"strconv"
	"time"
)

// CacheAPIPermissions 缓存接口权限，启动时调用
func CacheAPIPermissions() error {
	// 查询所有接口的权限
	var apis []models.API
	err := config.DB.Find(&apis).Error
	if err != nil {
		return fmt.Errorf("查询接口权限失败: %v", err)
	}

	// 将接口权限缓存到 Redis
	for _, api := range apis {
		permissionID := api.PermissionID
		key := constant.APIPermissionPrefix + strconv.Itoa(permissionID)

		// 添加接口信息到 Redis
		if err := SAdd(key, api.Path+":"+api.Method); err != nil {
			return fmt.Errorf("添加权限缓存失败: %v", err)
		}

		// 设置缓存过期时间（24 小时）
		if err := Expire(key, 24*time.Hour); err != nil {
			return fmt.Errorf("设置缓存过期时间失败: %v", err)
		}
	}

	return nil
}
