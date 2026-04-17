package model

import "time"

type CreateDonationRequest struct {
	Username   string  `json:"username"    validate:"required"`
	DonorName  *string `json:"donor_name"`                          // nil = anonymous
	DonorEmail *string `json:"donor_email" validate:"omitempty,email"`
	Amount     int64   `json:"amount"      validate:"required,min=1000"`
	Message    *string `json:"message"     validate:"omitempty,max=500"`
	Method     string  `json:"method"      validate:"required"`     // qris | va_bca | gopay | ...
}

type DonationResponse struct {
	ID          string         `json:"id"`
	CreatorID   string         `json:"creator_id"`
	DonorName   *string        `json:"donor_name"`
	Amount      int64          `json:"amount"`
	NetAmount   int64          `json:"net_amount"`
	Message     *string        `json:"message"`
	Status      string         `json:"status"`
	CheckoutURL *string        `json:"checkout_url,omitempty"`
	PaidAt      *time.Time     `json:"paid_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

type DonationFilter struct {
	Status string `query:"status"`
	From   string `query:"from"`
	To     string `query:"to"`
	Cursor string `query:"cursor"`
	Limit  int    `query:"limit"`
}
