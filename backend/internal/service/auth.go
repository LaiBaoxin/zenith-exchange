package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	dao "github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/pkg/utils"
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

	signedToken, err := utils.GenerateToken(user.ID, user.WalletAddress)
	if err != nil {
		return "", "", fmt.Errorf("token generation failed: %v", err)
	}

	return signedToken, user.WalletAddress, nil
}
