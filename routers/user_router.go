package routers

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/user"
	"mental/middleware"
)

// InitUserRouter 基础路由 Accept表示接受的请求参数类型 Produce表示返回结果类型
func InitUserRouter(r *gin.Engine) {
	adminRouter := r.Group("/")
	{

		adminRouter.POST("/login", user.UserController{}.Login) // 登录

		adminRouter.POST("/register", user.UserController{}.Register) // 注册

		adminRouter.Use(middleware.JWTMiddleWare()) // 需要鉴权中间件

		adminRouter.GET("", user.UserController{}.GetUserInfo) // 获取基本信息

		adminRouter.POST("/logout", user.UserController{}.Logout) // 退出登录

		adminRouter.POST("/refresh-token", user.UserController{}.RefreshToken) // 刷新令牌

		adminRouter.POST("/avatar", user.UserController{}.UpdateUserAvatar) // 修改头像

		adminRouter.POST("/base", user.UserController{}.UpdateUsernameOrEmail) // 修改用户名、邮箱

		adminRouter.POST("/password", user.UserController{}.UpdatePassword) // 修改密码

	}
}
