package impl

import (
	"errors"

	"github.com/kreatip/kreatip-backend/internal/model"
)

// Sentinel — pakai errors.Is untuk cek
var ErrNotImplemented = errors.New("not implemented")

// Business errors — return langsung ke controller, ErrorHandler yang handle HTTP status
var (
	ErrEmailAlreadyTaken       = model.ErrConflict("email", "email already taken")
	ErrInvalidCredential       = model.ErrUnauthorized("invalid email or password")
	ErrEmailNotVerified        = model.ErrForbidden("email not verified")
	ErrInvalidVerifyToken      = model.ErrBadRequest("invalid or expired verification token")
	ErrInvalidRefreshToken     = model.ErrUnauthorized("invalid or expired refresh token")
	ErrInternalCheckEmail      = model.ErrInternal("failed to check email availability")
	ErrInternalProcessPassword = model.ErrInternal("failed to process password")
	ErrInternalCreateUser      = model.ErrInternal("failed to create user")
	ErrInternalCreateWallet    = model.ErrInternal("failed to create wallet")
	ErrInternalSaveToken       = model.ErrInternal("failed to save token")
)
