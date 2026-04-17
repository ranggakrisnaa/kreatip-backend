package entity

import (
	"time"

	"gorm.io/gorm"
)

type WithdrawalStatus string

const (
	WithdrawalStatusRequested  WithdrawalStatus = "requested"
	WithdrawalStatusApproved   WithdrawalStatus = "approved"
	WithdrawalStatusProcessing WithdrawalStatus = "processing"
	WithdrawalStatusSuccess    WithdrawalStatus = "success"
	WithdrawalStatusFailed     WithdrawalStatus = "failed"
	WithdrawalStatusRejected   WithdrawalStatus = "rejected"
	WithdrawalStatusCancelled  WithdrawalStatus = "cancelled"
)

type Withdrawal struct {
	ID             string           `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         string           `gorm:"type:uuid;not null;index"`
	BankAccountID  string           `gorm:"type:uuid;not null"`
	Amount         int64            `gorm:"not null"`
	Fee            int64            `gorm:"not null;default:0"`
	NetAmount      int64            `gorm:"not null"`
	Status         WithdrawalStatus `gorm:"type:varchar(20);not null;default:'requested';index"`
	ExternalID     *string          `gorm:"type:varchar(255)"` // xendit disbursement id
	FailureReason  *string          `gorm:"type:text"`
	ApprovedBy     *string          `gorm:"type:uuid"`         // admin user id
	ApprovedAt     *time.Time
	CompletedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt   `gorm:"index"`

	BankAccount *BankAccount `gorm:"foreignKey:BankAccountID"`
}

func (Withdrawal) TableName() string { return "withdrawals" }
