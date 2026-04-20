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
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{Data: resp})
}

// Login godoc
// @Summary      Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      model.LoginRequest  true  "Login payload"
// @Success      200   {object}  model.WebResponse[model.AuthResponse]
// @Failure      400   {object}  model.ErrorResponse
// @Failure      401   {object}  model.ErrorResponse  "Invalid credentials"
// @Failure      403   {object}  model.ErrorResponse  "Email not verified"
// @Failure      422   {object}  model.ErrorResponse  "Validation error"
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	req := new(model.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}
	req.UserAgent = c.Get("User-Agent")
	req.IP = c.IP()

	resp, err := ctrl.UseCase.Login(c.UserContext(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[*model.AuthResponse]{Data: resp})
}

// Refresh godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      model.RefreshTokenRequest  true  "Refresh token"
// @Success      200   {object}  model.WebResponse[model.AuthResponse]
// @Failure      401   {object}  model.ErrorResponse  "Invalid or expired refresh token"
// @Router       /auth/refresh [post]
func (ctrl *AuthController) Refresh(c *fiber.Ctx) error {
	req := new(model.RefreshTokenRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := ctrl.UseCase.Refresh(c.UserContext(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[*model.AuthResponse]{Data: resp})
}

// VerifyEmail godoc
// @Summary      Verify email address
// @Tags         auth
// @Param        token  query  string  true  "Verification token"
// @Success      200    {object}  model.WebResponse[string]
// @Failure      400    {object}  model.ErrorResponse  "Invalid or expired token"
// @Router       /auth/verify-email [get]
func (ctrl *AuthController) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return fiber.ErrBadRequest
	}

	if err := ctrl.UseCase.VerifyEmail(c.UserContext(), token); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[string]{Data: "email verified"})
}

// ResendVerification godoc
// @Summary      Resend email verification
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  model.ResendVerificationRequest  true  "Email"
// @Success      200   {object}  model.WebResponse[string]
// @Failure      422   {object}  model.ErrorResponse
// @Router       /auth/resend-verification [post]
func (ctrl *AuthController) ResendVerification(c *fiber.Ctx) error {
	req := new(model.ResendVerificationRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := ctrl.UseCase.ResendVerification(c.UserContext(), req.Email); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[string]{Data: "verification email sent"})
}

// POST /api/v1/auth/forgot-password
func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/auth/reset-password
func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
