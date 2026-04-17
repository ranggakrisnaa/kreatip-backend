package entity

import "time"

type RefreshToken struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string     `gorm:"type:uuid;not null;index"`
	TokenHash string     `gorm:"type:varchar(64);uniqueIndex;not null"` // SHA-256 of raw token
	UserAgent string     `gorm:"type:varchar(500)"`
	IP        string     `gorm:"type:varchar(45)"`
	ExpiresAt time.Time  `gorm:"not null;index"`
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (RefreshToken) TableName() string { return "refresh_tokens" }
