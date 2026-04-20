package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/model"
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

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account. A wallet is automatically created alongside the user.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      model.RegisterRequest  true  "Register payload"
// @Success      201   {object}  model.WebResponse[model.UserResponse]
// @Failure      400   {object}  model.ErrorResponse  "Invalid request body"
// @Failure      409   {object}  model.ErrorResponse  "Email already taken"
// @Failure      422   {object}  model.ErrorResponse  "Validation error"
// @Router       /auth/register [post]
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	req := new(model.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := ctrl.UseCase.Register(c.UserContext(), req)
	if err != nil {
		return err // ErrorHandler yang map ke HTTP response
	}

	return c.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{Data: resp})
}

// POST /api/v1/auth/login
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/refresh
func (ctrl *AuthController) Refresh(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/verify-email
func (ctrl *AuthController) VerifyEmail(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/forgot-password
func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/reset-password
func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
