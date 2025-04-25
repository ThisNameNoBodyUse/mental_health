package utils

// 分布式锁工具类
import (
	"github.com/bsm/redislock"
	"mental/config"
	"time"
)

// 获取分布式锁
func TryLock(key string, ttl time.Duration) (*redislock.Lock, error) {
	lock, err := config.Locker.Obtain(ctx, key, ttl, nil)
	if err != nil {
		return nil, err
	}
	return lock, nil
}

// 释放分布式锁
func Unlock(lock *redislock.Lock) error {
	return lock.Release(ctx)
}
