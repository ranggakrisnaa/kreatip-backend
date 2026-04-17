package usecase

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
	"gorm.io/gorm"
)

type LedgerUseCase interface {
	// PostDonationEntries records ledger entries for a paid donation within tx.
	// Must be called inside an existing DB transaction.
	PostDonationEntries(ctx context.Context, donation *entity.Donation, tx *gorm.DB) error

	// HoldFunds moves amount from creator_balance → pending_withdrawal within tx.
	HoldFunds(ctx context.Context, walletID string, amount int64, tx *gorm.DB) error

	// ReleaseHeldFunds reverses a hold (failed withdrawal).
	ReleaseHeldFunds(ctx context.Context, walletID string, amount int64, tx *gorm.DB) error

	// SettleWithdrawal moves amount from pending_withdrawal → payment_gateway (disbursed).
	SettleWithdrawal(ctx context.Context, walletID string, amount int64, tx *gorm.DB) error
}
