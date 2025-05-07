package main

import (
	"fmt"
	"github.com/gin-contrib/cors" // 引入 CORS 中间件
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"mental/config"
	"mental/routers"
	"mental/utils"
	"time"
)

// @title 某系统
// @version 1.0
// @description API文档
// @host localhost:8080
// @BasePath /
func main() {
	config.InitAll() // 初始化所有配置

	// 创建 Gin 实例
	r := gin.Default()

	// 自定义 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 明确指定前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许携带凭证
		MaxAge:           12 * time.Hour,
	}))

	// 处理预检请求

	r.OPTIONS("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(200)
	})

	// 设置最大上传文件大小为 20MB
	r.MaxMultipartMemory = 20 << 20 // 20MB = 20 * 1024 * 1024 = 20 << 20

	// 设置静态文件路由，暴露 /storage 文件夹
	r.Static("/storage", "./storage")

	// 添加 Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 路由初始化
	routers.InitAdminRouter(r)
	routers.InitCommonRouter(r)

	// 自动生成 swagger 文档（如果 docs 不存在或首次运行）
	err := utils.RunSwagInit()
	if err != nil {
		fmt.Printf("Swagger 文档生成失败: %v\n", err)
		return
	}

	// 再插入 API 到数据库
	err = utils.InsertSwaggerAPIs("docs/swagger.json")
	if err != nil {
		fmt.Printf("插入 Swagger 接口失败: %v\n", err)
	} else {
		fmt.Println("Swagger 接口信息已成功插入数据库。")
	}

	// 缓存 权限对应的接口列表
	err = utils.CacheAPIPermissions()
	if err != nil {
		fmt.Printf("缓存接口权限失败: %v\n", err)
	}

	// 启动服务
	r.Run()
}
