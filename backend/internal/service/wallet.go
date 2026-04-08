package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
)

type WalletService struct{}

func NewWalletService() *WalletService {
	return &WalletService{}
}

// CheckUserAsset 检查用户的余额
func (s *WalletService) CheckUserAsset(userID uint64, asset string) error {
	// 查询当前账户状态
	var account model.Account
	if err := db.DB.Where("user_id = ? AND currency = ?", userID, asset).First(&account).Error; err != nil {
		return err
	}

	// 从流水表聚合该用户该币种的所有变动
	var sumResult struct {
		TotalChange float64
	}
	err := db.DB.Model(&model.BalanceLog{}).
		Where("user_id = ? AND currency = ?", userID, asset).
		Select("SUM(amount) as total_change").
		Scan(&sumResult).Error

	if err != nil {
		return fmt.Errorf("聚合流水失败: %v", err)
	}

	currentAvailable, _ := decimal.NewFromString(account.Available)
	sumAmount := decimal.NewFromFloat(sumResult.TotalChange)

	// 保证精度不丢失
	if !currentAvailable.Equal(sumAmount) {
		return fmt.Errorf("警告：资产不一致！账户余额: %s, 流水统计: %s, 差额: %s",
			currentAvailable.String(),
			sumAmount.String(),
			currentAvailable.Sub(sumAmount).String(),
		)
	}

	return nil
}
