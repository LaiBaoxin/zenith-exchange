package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义秘钥（实际应从环境变量获取）
var jwtSecret = []byte("zenith_exchange_secret_2026")

// CustomClaims 自定义载荷
type CustomClaims struct {
	Address string `json:"address"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(address string) (string, error) {
	claims := CustomClaims{
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置 24 小时过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "zenith-exchange",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析并验证 Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
