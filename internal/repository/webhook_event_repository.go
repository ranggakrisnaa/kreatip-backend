package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WebhookEventRepository interface {
	// UpsertByIdempotencyKey inserts or does nothing if already exists.
	// Returns (isNew, error).
	UpsertByIdempotencyKey(ctx context.Context, event *entity.WebhookEvent) (bool, error)
	MarkProcessed(ctx context.Context, id string, errMsg *string) error
}

type webhookEventRepository struct {
	Repository[entity.WebhookEvent]
	Log *logrus.Logger
}

func NewWebhookEventRepository(db *gorm.DB, log *logrus.Logger) WebhookEventRepository {
	return &webhookEventRepository{
		Repository: Repository[entity.WebhookEvent]{DB: db},
		Log:        log,
	}
}

func (r *webhookEventRepository) UpsertByIdempotencyKey(ctx context.Context, event *entity.WebhookEvent) (bool, error) {
	result := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "idempotency_key"}},
			DoNothing: true,
		}).Create(event)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *webhookEventRepository) MarkProcessed(ctx context.Context, id string, errMsg *string) error {
	updates := map[string]interface{}{
		"processed": true,
	}
	if errMsg != nil {
		updates["error"] = errMsg
	}
	return r.DB.WithContext(ctx).Model(&entity.WebhookEvent{}).
		Where("id = ?", id).Updates(updates).Error
}
