package config

// Redis连接配置
import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
	"time"
)

// RedisConfig 存储 Redis 相关配置
type RedisConfig struct {
	Host      string
	Port      int
	Password  string
	DB        int
	MaxActive int
	MaxIdle   int
	MinIdle   int
	MaxWait   time.Duration
}

// Redis 客户端
var RDB *redis.Client

// RedisSettings 全局 Redis 配置
var RedisSettings RedisConfig

var Locker *redislock.Client

// InitRedis 读取配置并初始化 Redis 连接
func InitRedis() error {
	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		return fmt.Errorf("加载 Redis 配置失败: %v", err)
	}

	// 读取配置，不使用 MustString()
	RedisSettings.Host = cfg.Section("redis").Key("host").String()
	RedisSettings.Port, _ = cfg.Section("redis").Key("port").Int()
	RedisSettings.Password = cfg.Section("redis").Key("password").String()
	RedisSettings.DB, _ = cfg.Section("redis").Key("db").Int()
	RedisSettings.MaxActive, _ = cfg.Section("redis").Key("max_active").Int()
	RedisSettings.MaxIdle, _ = cfg.Section("redis").Key("max_idle").Int()
	RedisSettings.MinIdle, _ = cfg.Section("redis").Key("min_idle").Int()
	// 最大等待时间，将纳秒转换为毫秒，因为Duration单位是纳秒，将纳秒单位换成毫秒单位
	RedisSettings.MaxWait = time.Duration(cfg.Section("redis").Key("max_wait").MustInt(1000)) * time.Millisecond // 这里保持默认值

	// 初始化 Redis 连接
	RedisAddr := fmt.Sprintf("%s:%d", RedisSettings.Host, RedisSettings.Port)
	RDB = redis.NewClient(&redis.Options{
		Addr:         RedisAddr,
		Password:     RedisSettings.Password,
		DB:           RedisSettings.DB,
		PoolSize:     RedisSettings.MaxActive,
		MinIdleConns: RedisSettings.MinIdle,
	})

	// 初始化 `redislock`
	Locker = redislock.New(RDB)

	// 测试连接
	ctx := context.Background()
	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis 连接失败: %v", err)
	}

	fmt.Println("Redis 连接成功!")
	return nil
}
