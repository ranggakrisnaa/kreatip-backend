package usecase

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/entity"
)

type AlertUseCase interface {
	// AuthorizeToken validates an alert token and returns the userID.
	AuthorizeToken(ctx context.Context, token string) (string, error)
	// BroadcastDonation pushes a donation event to the creator's WS room.
	BroadcastDonation(ctx context.Context, donation *entity.Donation) error
}
