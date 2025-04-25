package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"mental/config"
	"mental/models"
	"time"
)

// GenerateJWT 生成 JWT
func GenerateJWT(user *models.User, isAccessToken bool) (string, error) {
	// 选择使用访问令牌密钥还是刷新令牌密钥 如果是true，就生成访问令牌，否则，生成刷新令牌
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

	// 生成 JWT Claims
	claims := jwt.MapClaims{
		"id":       user.Id,
		"account":  user.Account,
		"username": user.Username,
		"exp":      expirationTime.Unix(),
		"jti":      GenerateJTI(), // JWT 唯一 ID
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
