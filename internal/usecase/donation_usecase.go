package usecase

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/model"
)

type DonationUseCase interface {
	Create(ctx context.Context, req *model.CreateDonationRequest) (*model.DonationResponse, error)
	GetByID(ctx context.Context, id string) (*model.DonationResponse, error)
	ListByCreator(ctx context.Context, creatorID string, filter model.DonationFilter) (*model.WebResponsePaged[model.DonationResponse], error)
}
