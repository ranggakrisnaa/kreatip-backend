package entity

import "time"

type Payment struct {
	ID             string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DonationID     string  `gorm:"type:uuid;uniqueIndex;not null"`
	Provider       string  `gorm:"type:varchar(30);not null"` // xendit | midtrans
	Method         string  `gorm:"type:varchar(30);not null"` // qris | va_bca | gopay | ovo | ...
	ExternalID     string  `gorm:"type:varchar(255);index"`   // provider's reference id
	ExternalStatus string  `gorm:"type:varchar(50)"`
	CheckoutURL    *string `gorm:"type:varchar(500)"`
	RawRequest     string  `gorm:"type:jsonb;default:'{}'"`
	RawResponse    string  `gorm:"type:jsonb;default:'{}'"`
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Payment) TableName() string { return "payments" }
