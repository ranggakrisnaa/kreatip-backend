package entity

import (
	"time"

	"gorm.io/gorm"
)

type BankAccount struct {
	ID                string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID            string         `gorm:"type:uuid;not null;index"`
	Type              string         `gorm:"type:varchar(10);not null"` // bank | ewallet
	BankCode          string         `gorm:"type:varchar(20);not null"` // BCA | BNI | GOPAY | OVO | ...
	AccountNumber     string         `gorm:"type:varchar(30);not null"`
	AccountHolderName string         `gorm:"type:varchar(100);not null"`
	IsVerified        bool           `gorm:"not null;default:false"`
	IsDefault         bool           `gorm:"not null;default:false"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (BankAccount) TableName() string { return "bank_accounts" }
