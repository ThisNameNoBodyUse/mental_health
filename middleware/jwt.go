package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"mental/constant"
	"mental/dao"
	"mental/utils"
	"net/http"
	"strconv"
	"time"
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

		// 访问令牌解析
		_, claims, err := utils.ParseJWT(token, true)
		if err != nil {
			// 无效的令牌
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Redis黑名单检查是否已经退出登录
		jti, ok := claims["jti"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
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

		// 获取 "id" 字段，并将其转换为 int64（由于用字符串存储，需先断言为 string 再转换）
		idStr, ok := claims["id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		userId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}

		// 获取用户角色列表
		roles, ok := claims["roles"].([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 获取用户权限列表
		var permissions []string
		for _, role := range roles {
			roleID := fmt.Sprintf("%v", role)

			// 先查询 Redis 是否有缓存的权限列表
			rolePermissions, err := utils.SMembers(constant.RolePermissionPrefix + roleID)

			// 如果 Redis 返回的是空集合，说明没有缓存，去数据库查询
			if err != nil || len(rolePermissions) == 0 {
				rolePermissions, err = dao.GetRolePermissionsFromDB(roleID)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get permissions"})
					return
				}

				// 将数据库查询结果缓存到 Redis
				if len(rolePermissions) > 0 {
					err = utils.SAdd(constant.RolePermissionPrefix+roleID, rolePermissions)
					if err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache permissions"})
						return
					}

					// 设置缓存过期时间（12 小时）
					err = utils.Expire(constant.RolePermissionPrefix+roleID, 12*time.Hour)
					if err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache expiration"})
						return
					}
				}
			}

			// 将 rolePermissions 合并到 permissions 中
			permissions = append(permissions, rolePermissions...)
		}

		// 判断当前用户的权限列表对应的接口集合中 是否存在当前需要访问的接口
		requestKey := c.FullPath() + ":" + c.Request.Method
		fmt.Println("requestKey : " + requestKey)

		hasPermission := false
		for _, permission := range permissions {
			permID := fmt.Sprintf("%v", permission)
			key := constant.APIPermissionPrefix + permID
			fmt.Println("key : " + key)

			// 判断 Redis 中是否存在该权限对应的接口集合
			exists, err := utils.Exists(key)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
				return
			}

			if !exists {
				// 如果 Redis 中没有，查数据库并缓存
				apiList, err := dao.GetPermissionAPIsFromDB(permID)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
					return
				}

				// 缓存到 Redis
				if len(apiList) > 0 {
					err = utils.SAdd(key, stringSliceToInterfaceSlice(apiList)...)
					if err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache permissions"})
						return
					}
					err = utils.Expire(key, 24*time.Hour)
					if err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to set expiration"})
						return
					}
				}
			}

			// 检查该权限对应的接口集合中是否包含当前请求
			inSet, err := utils.SIsMember(key, requestKey)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
				return
			}
			if inSet {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		// 令牌校验成功，将必要信息存入gin上下文中
		c.Set("id", userId)
		c.Set("account", claims["account"].(string))

		// 放行
		c.Next()
	}
}

// 把每个 string 转为 interface{} 存到新切片中
func stringSliceToInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, v := range strs {
		result[i] = v
	}
	return result
}
