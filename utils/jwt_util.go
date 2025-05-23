package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"mental/config"
	"mental/constant"
	"mental/models"
	"strconv"
	"time"
)

// GenerateJWT 生成 JWT
func GenerateJWT(user *models.User, isAccessToken bool) (string, error) {
	// 选择使用访问令牌密钥还是刷新令牌密钥
	var secretKey string
	var ttl int64

	if isAccessToken {
		secretKey = config.JWTSettings.SecretKey
		ttl = config.JWTSettings.TTL
	} else {
		secretKey = config.JWTSettings.RefreshSecretKey
		ttl = config.JWTSettings.RefreshTTL
	}

	// 计算过期时间
	expirationTime := time.Now().Add(time.Duration(ttl) * time.Millisecond)

	// 获取用户的角色列表
	roles, err := GetUserRoles(user.Id)
	if err != nil {
		return "", fmt.Errorf("获取用户角色列表失败: %v", err)
	}

	// 生成 JWT Claims
	claims := jwt.MapClaims{
		"id":       user.Id,               // 用户ID
		"account":  user.Account,          // 用户账号
		"username": user.Username,         // 用户名
		"roles":    roles,                 // 用户角色列表
		"exp":      expirationTime.Unix(), // 过期时间
		"jti":      GenerateJTI(),         // JWT 唯一 ID
	}

	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseJWT 解析 JWT，根据布尔值判断是访问令牌还是刷新令牌，选择对应的密钥
func ParseJWT(tokenString string, isAccessToken bool) (*jwt.Token, jwt.MapClaims, error) {
	// 选择正确的密钥
	var secretKey string
	if isAccessToken {
		secretKey = config.JWTSettings.SecretKey
	} else {
		secretKey = config.JWTSettings.RefreshSecretKey
	}

	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, nil, err
	}

	// 获取 Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, err
	}

	return token, claims, nil
}

// GetExpireTime 获取令牌过期时间，根据布尔值判断是访问令牌还是刷新令牌，选择对应的密钥
func GetExpireTime(tokenString string, isAccessToken bool) (int64, error) {
	_, claims, err := ParseJWT(tokenString, isAccessToken)
	if err != nil {
		return 0, err
	}

	exp := int64(claims["exp"].(float64))
	return exp - time.Now().Unix(), nil
}

// GenerateJTI 生成 JTI（UUID）
func GenerateJTI() string {
	return uuid.New().String()
}

// GetUserRoles 查询用户的角色列表
func GetUserRoles(userID int) ([]string, error) {
	var roles []string

	// 1. 从 Redis 中查询用户角色列表
	key := constant.UserRolePrefix + strconv.Itoa(userID)
	roles, err := SMembers(key)
	if err != nil {
		return nil, fmt.Errorf("从 Redis 查询用户角色失败: %v", err)
	}

	// 2. 如果 Redis 中有数据，直接返回
	if len(roles) > 0 {
		return roles, nil
	}

	// 3. 如果 Redis 中没有数据，从数据库中查询
	err = config.DB.Raw(`
		SELECT r.id
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`, userID).Scan(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("从数据库查询用户角色失败: %v", err)
	}

	// 4. 将查询结果写入 Redis
	if len(roles) > 0 {
		// 将角色列表写入 Redis 的 Set 中
		if err := SAdd(key, roles); err != nil {
			return nil, fmt.Errorf("角色列表写入 Redis 失败: %v", err)
		}

		// 设置过期时间（12小时）
		if err := Expire(key, 12*time.Hour); err != nil {
			return nil, fmt.Errorf("角色列表设置 Redis 过期时间失败: %v", err)
		}
	}

	return roles, nil
}
