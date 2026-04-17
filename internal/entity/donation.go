package entity

import (
	"time"

	"gorm.io/gorm"
)

// DonationStatus represents the lifecycle of a donation.
type DonationStatus string

const (
	DonationStatusPending  DonationStatus = "pending"
	DonationStatusPaid     DonationStatus = "paid"
	DonationStatusFailed   DonationStatus = "failed"
	DonationStatusExpired  DonationStatus = "expired"
	DonationStatusRefunded DonationStatus = "refunded"
)

type Donation struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatorID   string         `gorm:"type:uuid;not null;index"`    // refs creator_profiles.id
	DonorName   *string        `gorm:"type:varchar(100)"`           // nil = anonymous
	DonorEmail  *string        `gorm:"type:varchar(255)"`
	Amount      int64          `gorm:"not null"`                    // gross, in IDR
	PlatformFee int64          `gorm:"not null;default:0"`
	PaymentFee  int64          `gorm:"not null;default:0"`
	NetAmount   int64          `gorm:"not null"`                    // Amount - PlatformFee - PaymentFee
	Message     *string        `gorm:"type:text"`
	Status      DonationStatus `gorm:"type:varchar(20);not null;default:'pending';index"`
	Currency    string         `gorm:"type:varchar(3);not null;default:'IDR'"`
	PaidAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Payment *Payment `gorm:"foreignKey:DonationID"`
}

func (Donation) TableName() string { return "donations" }
