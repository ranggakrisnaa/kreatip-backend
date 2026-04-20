package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User, tx *gorm.DB) error
	MarkEmailVerified(ctx context.Context, userID string) error
}

type userRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(db *gorm.DB, log *logrus.Logger) UserRepository {
	return &userRepository{
		Repository: Repository[entity.User]{DB: db},
		Log:        log,
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User, tx *gorm.DB) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) MarkEmailVerified(ctx context.Context, userID string) error {
	return r.DB.WithContext(ctx).Model(&entity.User{}).
		Where("id = ?", userID).
		Update("email_verified", true).Error
}
