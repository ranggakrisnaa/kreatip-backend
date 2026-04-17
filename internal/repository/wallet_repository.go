package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	FindByUserID(ctx context.Context, userID string) (*entity.Wallet, error)
	FindByUserIDForUpdate(ctx context.Context, userID string, tx *gorm.DB) (*entity.Wallet, error)
	Create(ctx context.Context, wallet *entity.Wallet) error
	Save(ctx context.Context, wallet *entity.Wallet) error
	IncrementBalance(ctx context.Context, userID string, amount int64, tx *gorm.DB) error
}

type walletRepository struct {
	Repository[entity.Wallet]
	Log *logrus.Logger
}

func NewWalletRepository(db *gorm.DB, log *logrus.Logger) WalletRepository {
	return &walletRepository{
		Repository: Repository[entity.Wallet]{DB: db},
		Log:        log,
	}
}

func (r *walletRepository) FindByUserID(ctx context.Context, userID string) (*entity.Wallet, error) {
	var w entity.Wallet
	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *walletRepository) FindByUserIDForUpdate(ctx context.Context, userID string, tx *gorm.DB) (*entity.Wallet, error) {
	db := tx
	if db == nil {
		db = r.DB
	}
	var w entity.Wallet
	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *walletRepository) Create(ctx context.Context, wallet *entity.Wallet) error {
	return r.DB.WithContext(ctx).Create(wallet).Error
}

func (r *walletRepository) IncrementBalance(ctx context.Context, userID string, amount int64, tx *gorm.DB) error {
	db := tx
	if db == nil {
		db = r.DB
	}
	return db.WithContext(ctx).Model(&entity.Wallet{}).
		Where("user_id = ?", userID).
		UpdateColumn("balance_available", gorm.Expr("balance_available + ?", amount)).Error
}
