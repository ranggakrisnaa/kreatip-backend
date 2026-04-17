package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string         `gorm:"type:varchar(255);not null"`
	Role          string         `gorm:"type:varchar(20);not null;default:'user'"` // user | creator | admin
	EmailVerified bool           `gorm:"not null;default:false"`
	TOTPSecret    *string        `gorm:"type:varchar(64)"`
	TOTPEnabled   bool           `gorm:"not null;default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// associations (preload only, no FK constraint in struct)
	Profile      *CreatorProfile
	RefreshTokens []RefreshToken
}

func (User) TableName() string { return "users" }
