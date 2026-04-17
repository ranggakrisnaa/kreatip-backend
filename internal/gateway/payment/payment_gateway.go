package payment

import "context"

// CreateChargeRequest is the provider-agnostic payment request.
type CreateChargeRequest struct {
	ExternalID  string // must be unique, maps to donation.id
	Amount      int64
	Currency    string
	Method      string // qris | va_bca | gopay | ovo | dana | shopeepay | ...
	Description string
	CallbackURL string
	SuccessURL  string
	FailureURL  string
	// Payer info (optional)
	PayerName  string
	PayerEmail string
}

// CreateChargeResponse is the provider-agnostic payment response.
type CreateChargeResponse struct {
	ExternalID  string
	CheckoutURL string
	ExpiresAt   string
	RawResponse string // JSON string for audit
}

// Gateway is the interface that all payment providers must implement.
type Gateway interface {
	CreateCharge(ctx context.Context, req *CreateChargeRequest) (*CreateChargeResponse, error)
	// VerifyWebhookSignature verifies the raw payload + signature header.
	VerifyWebhookSignature(payload []byte, signature string) error
}
