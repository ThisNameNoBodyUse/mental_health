package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func (BaseController) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "",
		"data": data,
	})
}
func (BaseController) Error(c *gin.Context, data interface{}, err string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  err,
		"data": data,
	})
}
