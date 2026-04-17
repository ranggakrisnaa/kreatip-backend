package entity

import "time"

type AlertToken struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string     `gorm:"type:uuid;not null;index"`
	Token     string     `gorm:"type:varchar(64);uniqueIndex;not null"` // opaque, rotatable
	Settings  string     `gorm:"type:jsonb;default:'{}'"`               // duration, sound, template, animation
	RevokedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (AlertToken) TableName() string { return "alert_tokens" }
