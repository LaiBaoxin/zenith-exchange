package service

import (
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"gorm.io/gorm"
)

type AssetsService struct {
	db *gorm.DB
}

// GetUserBalances 获取用户所有币种余额
func (s *AssetsService) GetUserBalances(userID uint64) ([]model.Account, error) {
	var accounts []model.Account
	err := s.db.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}
