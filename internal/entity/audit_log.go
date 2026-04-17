package entity

import "time"

type AuditLog struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ActorUserID  *string   `gorm:"type:uuid;index"` // nil for system actions
	Action       string    `gorm:"type:varchar(50);not null;index"` // login | withdraw | admin_approve | ...
	ResourceType *string   `gorm:"type:varchar(50)"`
	ResourceID   *string   `gorm:"type:uuid"`
	Metadata     string    `gorm:"type:jsonb;default:'{}'"`
	IP           string    `gorm:"type:varchar(45)"`
	OccurredAt   time.Time `gorm:"not null;index"`
}

func (AuditLog) TableName() string { return "audit_logs" }
