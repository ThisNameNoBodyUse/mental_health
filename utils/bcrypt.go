package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 生成密码哈希
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 校验密码是否正确
func CheckPassword(hashedPassword, inputPassword string) bool { // 加密密码，明文密码比对
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
