package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CreatorProfileRepository interface {
	FindByUserID(ctx context.Context, userID string) (*entity.CreatorProfile, error)
	FindByUsername(ctx context.Context, username string) (*entity.CreatorProfile, error)
	Create(ctx context.Context, profile *entity.CreatorProfile) error
	Save(ctx context.Context, profile *entity.CreatorProfile) error
	ExistsUsername(ctx context.Context, username string) (bool, error)
}

type creatorProfileRepository struct {
	Repository[entity.CreatorProfile]
	Log *logrus.Logger
}

func NewCreatorProfileRepository(db *gorm.DB, log *logrus.Logger) CreatorProfileRepository {
	return &creatorProfileRepository{
		Repository: Repository[entity.CreatorProfile]{DB: db},
		Log:        log,
	}
}

func (r *creatorProfileRepository) FindByUserID(ctx context.Context, userID string) (*entity.CreatorProfile, error) {
	var p entity.CreatorProfile
	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *creatorProfileRepository) FindByUsername(ctx context.Context, username string) (*entity.CreatorProfile, error) {
	var p entity.CreatorProfile
	if err := r.DB.WithContext(ctx).Where("username = ?", username).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *creatorProfileRepository) Create(ctx context.Context, profile *entity.CreatorProfile) error {
	return r.DB.WithContext(ctx).Create(profile).Error
}

func (r *creatorProfileRepository) ExistsUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&entity.CreatorProfile{}).
		Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
