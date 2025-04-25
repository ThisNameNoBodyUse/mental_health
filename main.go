package main

import (
	"github.com/gin-gonic/gin"
	"mental/config"
	"mental/routers"
)

func main() {
	config.InitAll() // 初始化所有配置
	// gin.New() + 默认中间件
	r := gin.Default()
	// 设置最大上传文件大小为 20MB
	r.MaxMultipartMemory = 20 << 20 // 20MB = 20 * 1024 * 1024 = 20 << 20
	// 设置静态文件路由，暴露/storage文件夹
	r.Static("/storage", "./storage")
	// 路由初始化
	routers.InitAdminRouter(r)
	routers.InitCommonRouter(r)
	r.Run()
}
