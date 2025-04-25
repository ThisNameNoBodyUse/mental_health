package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"mental/constant"
	"mental/utils"
	"net/http"
)

// JWTMiddleWare 是 JWT鉴权中间件 对于某些需要进行登录校验的路由进行令牌鉴权
func JWTMiddleWare() gin.HandlerFunc { // gin.HandlerFunc用于定义中间件
	return func(c *gin.Context) {
		// 从请求头中获取访问令牌
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		_, claims, err := utils.ParseJWT(token, true) // 访问令牌解析
		if err != nil {
			// 无效的令牌
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// Redis黑名单检查是否已经退出登录 4-25
		jti := claims["jti"].(string)
		var tokenKey = constant.BlackListPrefix + jti
		_, err = utils.Get(tokenKey)
		if err == nil {
			// 此时说明已经在黑名单中，则放回令牌错误的信息
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		} else if !errors.Is(err, redis.Nil) { // 防止redis本身错误导致误判
			// 处理其他可能的 Redis 错误
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		// 获取 "id" 字段，并将其转换为 float64
		id, ok := claims["id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// 令牌校验成功，将必要信息存入gin上下文中
		// 将 float64 类型的 id 转换为 int64
		userId := int64(id)
		c.Set("id", userId)
		c.Set("account", claims["account"].(string))
		// 放行
		c.Next()
	}
}
