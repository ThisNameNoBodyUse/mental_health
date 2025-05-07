package routers

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/admin"
	"mental/middleware"
)

// InitAdminRouter 管理员基础路由 Accept表示接受的请求参数类型 Produce表示返回结果类型
func InitAdminRouter(r *gin.Engine) {
	adminRouter := r.Group("/admin")
	{

		adminRouter.POST("/login", admin.UserController{}.AdminLogin) // 管理员登录

		adminRouter.POST("/register", admin.UserController{}.AdminRegister) // 管理员注册

		adminRouter.Use(middleware.JWTMiddleWare()) // 需要鉴权中间件

		adminRouter.GET("", admin.UserController{}.GetAdminInfo) // 获取基本信息

		adminRouter.POST("/logout", admin.UserController{}.AdminLogout) // 退出登录

		adminRouter.POST("/refresh-token", admin.UserController{}.RefreshToken) // 刷新令牌
	}
}
