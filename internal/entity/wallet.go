package entity

import "time"

type Wallet struct {
	ID               string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           string    `gorm:"type:uuid;uniqueIndex;not null"`
	Currency         string    `gorm:"type:varchar(3);not null;default:'IDR'"`
	BalanceAvailable int64     `gorm:"not null;default:0"` // cache; source of truth = ledger
	BalancePending   int64     `gorm:"not null;default:0"` // pending withdrawal hold
	LastReconciledAt *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (Wallet) TableName() string { return "wallets" }
