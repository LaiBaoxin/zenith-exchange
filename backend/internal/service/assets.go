package service

import (
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
)

type AssetsService struct{}

func NewAssetsService() *AssetsService {
	return &AssetsService{}
}

// GetUserBalances 获取用户所有币种的余额
func (s *AssetsService) GetUserBalances(userID uint64) ([]model.Account, error) {
	var accounts []model.Account
	err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}
