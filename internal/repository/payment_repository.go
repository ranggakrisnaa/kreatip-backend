package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	FindByDonationID(ctx context.Context, donationID string) (*entity.Payment, error)
	FindByExternalID(ctx context.Context, externalID string) (*entity.Payment, error)
	Create(ctx context.Context, payment *entity.Payment) error
	Save(ctx context.Context, payment *entity.Payment) error
}

type paymentRepository struct {
	Repository[entity.Payment]
	Log *logrus.Logger
}

func NewPaymentRepository(db *gorm.DB, log *logrus.Logger) PaymentRepository {
	return &paymentRepository{
		Repository: Repository[entity.Payment]{DB: db},
		Log:        log,
	}
}

func (r *paymentRepository) FindByDonationID(ctx context.Context, donationID string) (*entity.Payment, error) {
	var p entity.Payment
	if err := r.DB.WithContext(ctx).Where("donation_id = ?", donationID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) FindByExternalID(ctx context.Context, externalID string) (*entity.Payment, error) {
	var p entity.Payment
	if err := r.DB.WithContext(ctx).Where("external_id = ?", externalID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	return r.DB.WithContext(ctx).Create(payment).Error
}
