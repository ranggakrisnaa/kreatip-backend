package usecase

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/model"
)

type AuthUseCase interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error)
	VerifyEmail(ctx context.Context, token string) error
	Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*model.AuthResponse, error)
	ForgotPassword(ctx context.Context, req *model.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error
	Logout(ctx context.Context, userID string, jti string) error
}
