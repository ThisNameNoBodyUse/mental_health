package routers

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/admin"
	"mental/middleware"
)

// InitAdminRouter 管理员基础路由
func InitAdminRouter(r *gin.Engine) {
	adminRouter := r.Group("/admin")
	{
		// 不需要鉴权
		adminRouter.POST("/login", admin.UserController{}.AdminLogin)       // 管理员登录
		adminRouter.POST("/register", admin.UserController{}.AdminRegister) // 管理员注册

		// 需要鉴权中间件
		adminRouter.Use(middleware.JWTMiddleWare())
		adminRouter.GET("", admin.UserController{}.GetAdminInfo)                // 获取基本信息
		adminRouter.POST("/logout", admin.UserController{}.AdminLogout)         // 退出登录
		adminRouter.POST("/refresh-token", admin.UserController{}.RefreshToken) // 刷新令牌
	}

}
