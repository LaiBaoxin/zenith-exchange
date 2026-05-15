package model

import "time"

// BalanceLog 资产变更流水
type BalanceLog struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"index:idx_user_currency" json:"user_id"`
	Currency   string    `gorm:"size:20;index:idx_user_currency" json:"currency"`
	ChangeType string    `gorm:"size:20" json:"change_type"` // deposit, withdraw, trade, freeze
	Amount     float64   `gorm:"type:decimal(36,18)" json:"amount"`
	Balance    float64   `gorm:"type:decimal(36,18)" json:"balance"`
	LogTime    time.Time `gorm:"type:datetime(3);default:CURRENT_TIMESTAMP(3)" json:"log_time"`
}
