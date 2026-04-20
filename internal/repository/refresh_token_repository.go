package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken, tx *gorm.DB) error
	FindByTokenHash(ctx context.Context, hash string) (*entity.RefreshToken, error)
	RevokeByID(ctx context.Context, id string) error
	RevokeAllByUserID(ctx context.Context, userID string) error
}

type refreshTokenRepository struct {
	Repository[entity.RefreshToken]
	Log *logrus.Logger
}

func NewRefreshTokenRepository(db *gorm.DB, log *logrus.Logger) RefreshTokenRepository {
	return &refreshTokenRepository{
		Repository: Repository[entity.RefreshToken]{DB: db},
		Log:        log,
	}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken, tx *gorm.DB) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	if err := r.DB.WithContext(ctx).Where("token_hash = ?", hash).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) RevokeByID(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Model(&entity.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	return r.DB.WithContext(ctx).Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}
