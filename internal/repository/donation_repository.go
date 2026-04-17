package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DonationRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Donation, error)
	FindByIDWithPayment(ctx context.Context, id string) (*entity.Donation, error)
	Create(ctx context.Context, donation *entity.Donation) error
	Save(ctx context.Context, donation *entity.Donation) error
	ListByCreator(ctx context.Context, creatorID string, filter ListDonationFilter) ([]*entity.Donation, int64, error)
}

type ListDonationFilter struct {
	Status string
	From   string
	To     string
	Cursor string
	Limit  int
}

type donationRepository struct {
	Repository[entity.Donation]
	Log *logrus.Logger
}

func NewDonationRepository(db *gorm.DB, log *logrus.Logger) DonationRepository {
	return &donationRepository{
		Repository: Repository[entity.Donation]{DB: db},
		Log:        log,
	}
}

func (r *donationRepository) FindByIDWithPayment(ctx context.Context, id string) (*entity.Donation, error) {
	var d entity.Donation
	if err := r.DB.WithContext(ctx).Preload("Payment").First(&d, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *donationRepository) Create(ctx context.Context, donation *entity.Donation) error {
	return r.DB.WithContext(ctx).Create(donation).Error
}

func (r *donationRepository) ListByCreator(ctx context.Context, creatorID string, filter ListDonationFilter) ([]*entity.Donation, int64, error) {
	// TODO: implement cursor-based pagination + filters
	var donations []*entity.Donation
	var total int64
	q := r.DB.WithContext(ctx).Where("creator_id = ?", creatorID)
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	q.Count(&total)
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if err := q.Order("created_at DESC").Limit(limit).Find(&donations).Error; err != nil {
		return nil, 0, err
	}
	return donations, total, nil
}
