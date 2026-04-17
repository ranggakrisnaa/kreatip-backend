package entity

import "time"

type CreatorProfile struct {
	ID              string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          string  `gorm:"type:uuid;uniqueIndex;not null"`
	Username        string  `gorm:"type:varchar(30);uniqueIndex;not null"`
	DisplayName     string  `gorm:"type:varchar(100);not null"`
	AvatarURL       *string `gorm:"type:varchar(500)"`
	Bio             *string `gorm:"type:text"`
	SocialLinks     string  `gorm:"type:jsonb;default:'{}'"`  // {"youtube":"","twitch":"","instagram":""}
	ThemeColor      string  `gorm:"type:varchar(7);default:'#6366f1'"`
	MinDonation     int64   `gorm:"not null;default:1000"`
	MaxDonation     int64   `gorm:"not null;default:10000000"`
	ThankYouMessage string  `gorm:"type:text;default:'Terima kasih atas dukungannya!'"`
	IsActive        bool    `gorm:"not null;default:true"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (CreatorProfile) TableName() string { return "creator_profiles" }
