package utils

import (
	"context"
	"mental/config"
	"time"
)

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

// SAdd 向 Set 集合中添加元素
func SAdd(key string, members ...interface{}) error {
	return config.RDB.SAdd(ctx, key, members...).Err()
}

// SMembers 获取 Set 集合中的所有元素
func SMembers(key string) ([]string, error) {
	return config.RDB.SMembers(ctx, key).Result()
}

// SRem 从 Set 集合中移除元素
func SRem(key string, members ...interface{}) error {
	return config.RDB.SRem(ctx, key, members...).Err()
}

// SIsMember 检查元素是否在 Set 集合中
func SIsMember(key string, member interface{}) (bool, error) {
	return config.RDB.SIsMember(ctx, key, member).Result()
}

// Expire 设置键的过期时间
func Expire(key string, expiration time.Duration) error {
	return config.RDB.Expire(ctx, key, expiration).Err()
}

// Exists 判断 key 是否存在（适用于所有类型）
func Exists(key string) (bool, error) {
	count, err := config.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
