package model

import (
	"time"
)

// Account 资产账户表
type Account struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;uniqueIndex:uk_user_currency" json:"user_id"`
	Currency  string    `gorm:"type:varchar(20);not null;uniqueIndex:uk_user_currency" json:"currency"`
	Available float64   `gorm:"type:decimal(36,18);not null;default:0" json:"available"` // 注意：高精度金融计算建议用 shopspring/decimal
	Frozen    float64   `gorm:"type:decimal(36,18);not null;default:0" json:"frozen"`
	Version   uint32    `gorm:"not null;default:0" json:"version"` // 乐观锁
	UpdatedAt time.Time `gorm:"type:datetime(3);default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updated_at"`
}
