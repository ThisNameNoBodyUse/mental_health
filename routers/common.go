package routers

import (
	"github.com/gin-gonic/gin"
	"mental/controllers/admin"
)

func InitCommonRouter(r *gin.Engine) {
	commonRouter := r.Group("/common")
	{

		// commonRouter.Use(middleware.JWTMiddleWare())
		commonRouter.POST("/check-file", admin.FileController{}.Check)
		commonRouter.POST("/upload", admin.FileController{}.Upload)
	}
}
