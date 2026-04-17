package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	UseCase usecase.AuthUseCase
	Log     *logrus.Logger
}

func NewAuthController(uc usecase.AuthUseCase, log *logrus.Logger) *AuthController {
	return &AuthController{UseCase: uc, Log: log}
}

// POST /api/v1/auth/register
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	// TODO: implement
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/login
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	// TODO: implement
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/refresh
func (ctrl *AuthController) Refresh(c *fiber.Ctx) error {
	// TODO: implement
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/verify-email
func (ctrl *AuthController) VerifyEmail(c *fiber.Ctx) error {
	// TODO: implement
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/forgot-password
func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	// TODO: implement
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/reset-password
func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	// TODO: implement
	return fiber.ErrNotImplemented
}
