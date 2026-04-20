package model

// Task type constants — dipakai oleh producer (usecase) dan consumer (worker).
const (
	TypeEmailVerification = "email:verification"
	TypeEmailReceipt      = "email:receipt"
	TypeEmailPasswordReset = "email:password_reset"
)

type VerificationEmailPayload struct {
	UserID          string `json:"user_id"`
	Email           string `json:"email"`
	VerificationURL string `json:"verification_url"`
}

type ReceiptEmailPayload struct {
	Email      string `json:"email"`
	DonationID string `json:"donation_id"`
}

type PasswordResetEmailPayload struct {
	Email    string `json:"email"`
	ResetURL string `json:"reset_url"`
}
