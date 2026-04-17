package entity

import "time"

type WebhookEvent struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Provider       string     `gorm:"type:varchar(30);not null;index"`
	EventType      string     `gorm:"type:varchar(50);not null"`
	ExternalID     string     `gorm:"type:varchar(255);index"`
	IdempotencyKey string     `gorm:"type:varchar(255);uniqueIndex;not null"` // provider+external_id+event_type
	Payload        string     `gorm:"type:jsonb;not null;default:'{}'"`
	Processed      bool       `gorm:"not null;default:false"`
	Error          *string    `gorm:"type:text"`
	ReceivedAt     time.Time  `gorm:"not null"`
	ProcessedAt    *time.Time
}

func (WebhookEvent) TableName() string { return "webhook_events" }
