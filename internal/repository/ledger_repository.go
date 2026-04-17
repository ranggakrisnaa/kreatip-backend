package repository

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LedgerRepository interface {
	CreateBatch(ctx context.Context, entries []*entity.LedgerEntry) error
	SumByWalletAndAccount(ctx context.Context, walletID string, account entity.LedgerAccount) (int64, error)
}

type ledgerRepository struct {
	Repository[entity.LedgerEntry]
	Log *logrus.Logger
}

func NewLedgerRepository(db *gorm.DB, log *logrus.Logger) LedgerRepository {
	return &ledgerRepository{
		Repository: Repository[entity.LedgerEntry]{DB: db},
		Log:        log,
	}
}

func (r *ledgerRepository) CreateBatch(ctx context.Context, entries []*entity.LedgerEntry) error {
	return r.DB.WithContext(ctx).Create(entries).Error
}

func (r *ledgerRepository) SumByWalletAndAccount(ctx context.Context, walletID string, account entity.LedgerAccount) (int64, error) {
	var result struct {
		Credits int64
		Debits  int64
	}
	err := r.DB.WithContext(ctx).Model(&entity.LedgerEntry{}).
		Select("SUM(CASE WHEN direction = 'credit' THEN amount ELSE 0 END) as credits, SUM(CASE WHEN direction = 'debit' THEN amount ELSE 0 END) as debits").
		Where("wallet_id = ? AND account = ?", walletID, account).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Credits - result.Debits, nil
}
