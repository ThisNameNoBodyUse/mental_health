package routers

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/user"
	"mental/middleware"
)

// InitSCLRouter 初始化公用路由
func InitSCLRouter(r *gin.Engine) {
	commonRouter := r.Group("/scl")
	{
		commonRouter.Use(middleware.JWTMiddleWare())
		commonRouter.POST("", user.SCLController{}.InsertSCL)        // 用户导入scl数据
		commonRouter.GET("", user.SCLController{}.SelectAllByUserId) // 查询用户scl数据
		commonRouter.DELETE("", user.SCLController{}.DeleteSCL)      // 删除指定id的scl数据
		commonRouter.POST("/update", user.SCLController{}.UpdateSCL) // 更新指定id的scl数据
		commonRouter.GET("/all", user.SCLController{}.SelectSCLs)    // 查询所有用户的scl数据
	}
}
