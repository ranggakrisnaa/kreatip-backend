package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WithdrawalRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Withdrawal, error)
	Create(ctx context.Context, withdrawal *entity.Withdrawal) error
	Save(ctx context.Context, withdrawal *entity.Withdrawal) error
	ListByUser(ctx context.Context, userID string, status entity.WithdrawalStatus) ([]*entity.Withdrawal, error)
	ListPending(ctx context.Context) ([]*entity.Withdrawal, error)
}

type withdrawalRepository struct {
	Repository[entity.Withdrawal]
	Log *logrus.Logger
}

func NewWithdrawalRepository(db *gorm.DB, log *logrus.Logger) WithdrawalRepository {
	return &withdrawalRepository{
		Repository: Repository[entity.Withdrawal]{DB: db},
		Log:        log,
	}
}

func (r *withdrawalRepository) Create(ctx context.Context, withdrawal *entity.Withdrawal) error {
	return r.DB.WithContext(ctx).Create(withdrawal).Error
}

func (r *withdrawalRepository) ListByUser(ctx context.Context, userID string, status entity.WithdrawalStatus) ([]*entity.Withdrawal, error) {
	var items []*entity.Withdrawal
	q := r.DB.WithContext(ctx).Where("user_id = ?", userID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *withdrawalRepository) ListPending(ctx context.Context) ([]*entity.Withdrawal, error) {
	var items []*entity.Withdrawal
	if err := r.DB.WithContext(ctx).
		Where("status = ?", entity.WithdrawalStatusRequested).
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
