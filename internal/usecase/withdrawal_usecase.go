package usecase

import (
	"context"

	"github.com/kreatip/kreatip-backend/internal/model"
)

type WithdrawalUseCase interface {
	Request(ctx context.Context, userID string, req *model.CreateWithdrawalRequest) (*model.WithdrawalResponse, error)
	Cancel(ctx context.Context, userID string, withdrawalID string) error
	// Admin actions
	Approve(ctx context.Context, adminID string, withdrawalID string) error
	Reject(ctx context.Context, adminID string, withdrawalID string, reason string) error
	// Called by worker after approval
	ProcessDisbursement(ctx context.Context, withdrawalID string) error
	// Called by disbursement webhook
	HandleDisbursementWebhook(ctx context.Context, payload WebhookPayload) error
}
