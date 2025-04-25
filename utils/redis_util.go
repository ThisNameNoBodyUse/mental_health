package utils

import (
	"context"
	"mental/config"
	"time"
)

// Redis基本操作工具类
var ctx = context.Background()

// Set 设置键值对和过期时间
func Set(key string, value interface{}, expiration time.Duration) error {
	return config.RDB.Set(ctx, key, value, expiration).Err()
}

// Get 取value值
func Get(key string) (interface{}, error) {
	return config.RDB.Get(ctx, key).Result()
}

// Delete 删除键
func Delete(key string) error {
	return config.RDB.Del(ctx, key).Err()
}
