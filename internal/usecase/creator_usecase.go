package usecase

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/model"
)

type CreatorUseCase interface {
	GetProfile(ctx context.Context, userID string) (*model.PublicProfileResponse, error)
	GetPublicProfile(ctx context.Context, username string) (*model.PublicProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req *model.UpdateProfileRequest) (*model.PublicProfileResponse, error)
	UploadAvatar(ctx context.Context, userID string, filename string, data []byte) (string, error)
	GetAlertToken(ctx context.Context, userID string) (string, error)
	RotateAlertToken(ctx context.Context, userID string) (string, error)
	SaveAlertSettings(ctx context.Context, userID string, req *model.AlertSettingsRequest) error
}
