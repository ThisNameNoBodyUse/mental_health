package utils

import (
	"mental/config"
	"mental/constant"
	"mental/models"
	"strconv"
	"time"
)

// CacheAPIPermissions 缓存接口权限，启动的时候调用
func CacheAPIPermissions() error {
	// 查询所有接口的权限
	var apis []models.API
	err := config.DB.Find(&apis).Error
	if err != nil {
		return err
	}

	// 将接口权限缓存到 Redis
	// key 为 权限id，value为对应的接口集合
	for _, api := range apis {
		if api.PermissionID != nil {
			permissionID := *api.PermissionID // 解引用指针
			key := constant.APIPermissionPrefix + strconv.Itoa(permissionID)
			err := SAdd(key, api.Path+":"+api.Method)
			if err != nil {
				return err
			}
			// 设置缓存过期时间（24 小时）
			err = Expire(key, 24*time.Hour)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
