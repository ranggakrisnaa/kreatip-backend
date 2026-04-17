package usecase

import "context"

// WebhookPayload is the raw parsed webhook body from the payment gateway.
type WebhookPayload struct {
	Provider       string
	EventType      string
	ExternalID     string
	IdempotencyKey string
	RawPayload     string
}

type PaymentUseCase interface {
	// HandleDonationWebhook processes an incoming payment webhook.
	// Idempotent: safe to call multiple times with the same idempotency key.
	HandleDonationWebhook(ctx context.Context, payload WebhookPayload) error
}
