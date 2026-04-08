package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	dao "github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"gorm.io/gorm"
)

type AuthService struct{}

// 只保留纯粹的业务逻辑，不要把 Handler 写在这里
func (s *AuthService) LoginByAddress(address string) (string, string, error) {
	if dao.DB == nil {
		return "", "", errors.New("后端数据库连接对象为 nil，请检查初始化顺序")
	}

	if address == "" {
		return "", "", errors.New("invalid address")
	}

	var user model.User

	// 查找或创建用户
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("wallet_address = ?", address).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				user = model.User{
					WalletAddress: address,
					ApiKey:        uuid.New().String(),
				}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}

				// 初始化账户
				defaultCurrencies := []string{"USDT", "ETH", "BTC"}
				accounts := make([]model.Account, len(defaultCurrencies))
				for i, cur := range defaultCurrencies {
					accounts[i] = model.Account{
						UserID:    user.ID,
						Currency:  cur,
						Available: "0",
						Frozen:    "0",
						Version:   0,
					}
				}
				// 批量创建执行
				if err := tx.CreateInBatches(accounts, len(accounts)).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}

	expireHour := config.GlobalConfig.JWT.ExpireHour
	if expireHour <= 0 {
		expireHour = 24 // 默认 24 小时
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"address": user.WalletAddress,
		"exp":     time.Now().Add(time.Hour * time.Duration(expireHour)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, user.WalletAddress, nil
}
