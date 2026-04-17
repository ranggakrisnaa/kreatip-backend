package usecase

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/model"
)

type WalletUseCase interface {
	GetWallet(ctx context.Context, userID string) (*model.WalletResponse, error)
	ListLedger(ctx context.Context, userID string, cursor string, limit int) (*model.WebResponsePaged[model.LedgerEntryResponse], error)
}
