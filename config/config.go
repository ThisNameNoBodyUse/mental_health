package config

// InitAll 聚合所有的初始化配置，一起初始化
func InitAll() {
	InitDB()
	LoadJWTConfig()
	InitRedis()
}
