package routers

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/user"
	"mental/middleware"
)

// InitCommonRouter 初始化公用路由
func InitCommonRouter(r *gin.Engine) {
	commonRouter := r.Group("/common")
	{
		commonRouter.Use(middleware.JWTMiddleWare())
		commonRouter.POST("/check-file", user.FileController{}.Check)
		commonRouter.POST("/upload", user.FileController{}.Upload)
		commonRouter.POST("/parse-file", user.FileController{}.ParseFile) // 解析文件
	}
}
