package config

// JWT属性配置
import (
	"gopkg.in/ini.v1"
	"log"
)

// JWTConfig 存储 JWT 配置
type JWTConfig struct {
	SecretKey        string
	TTL              int64
	RefreshSecretKey string
	RefreshTTL       int64
}

// JWTSettings 作为全局变量存储 JWT 配置
var JWTSettings JWTConfig

// 读取 JWT 配置
func LoadJWTConfig() {
	// 读取配置文件
	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	JWTSettings = JWTConfig{
		SecretKey:        cfg.Section("jwt").Key("secretKey").String(),
		TTL:              cfg.Section("jwt").Key("ttl").MustInt64(0),
		RefreshSecretKey: cfg.Section("jwt").Key("refreshSecretKey").String(),
		RefreshTTL:       cfg.Section("jwt").Key("refresh-ttl").MustInt64(0),
	}
}
