package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"time"
)

type CustomClaims struct {
	UserID  uint64 `json:"user_id"`
	Address string `json:"address"`
	jwt.RegisteredClaims
}

// GenerateToken 生成函数
func GenerateToken(userID uint64, address string) (string, error) {
	claims := CustomClaims{
		UserID:  userID,
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "zenith-exchange",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

// ParseToken 解析函数
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func getSecret() []byte {
	secret := config.GlobalConfig.JWT.Secret
	if secret == "" {
		return []byte("zenith_exchange_secure_key_2026_@#!$")
	}
	return []byte(secret)
}
