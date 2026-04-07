package model

import "time"

// User 用户表
type User struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	WalletAddress string    `gorm:"type:char(42);not null;uniqueIndex:idx_wallet" json:"wallet_address"`
	ApiKey        string    `gorm:"type:varchar(64);unique" json:"api_key"`
	CreatedAt     time.Time `gorm:"type:datetime(3);default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	Accounts      []Account `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
}
