package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/golang-jwt/jwt/v5"
	dao "github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"gorm.io/gorm"
)

type AuthService struct{}

func (s *AuthService) LoginByAddress(address string) (string, string, error) {
	if dao.DB == nil {
		return "", "", errors.New("后端数据库连接对象为 nil，请检查初始化顺序")
	}

	if address == "" {
		return "", "", errors.New("invalid address")
	}

	var user model.User

	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 查找或创建用户
		err := tx.Where("wallet_address = ?", address).First(&user).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建记录
			user = model.User{
				WalletAddress: address,
				ApiKey:        uuid.New().String(),
				CreatedAt:     time.Now(),
			}
			if err = tx.Create(&user).Error; err != nil {
				return fmt.Errorf("failed to create user: %v", err)
			}

			// 初始化资产账户 (可选：默认给新用户开启常用币种账户)
			defaultCurrencies := []string{"USDT", "ETH", "BTC"}
			for _, cur := range defaultCurrencies {
				account := model.Account{
					UserID:    user.ID,
					Currency:  cur,
					Available: 0,
					Frozen:    0,
					Version:   0,
				}
				if err = tx.Create(&account).Error; err != nil {
					return fmt.Errorf("failed to init account for %s: %v", cur, err)
				}
			}
		} else if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}

	// 生成 JWT Token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"address": user.WalletAddress,
		"exp":     time.Now().Add(time.Hour * time.Duration(config.GlobalConfig.JWT.ExpireHour)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, user.WalletAddress, nil
}
