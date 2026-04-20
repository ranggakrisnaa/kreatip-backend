package route

import (
	handler "github.com/kreatip/kreatip-backend/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

type Controllers struct {
	Auth *handler.AuthController
}

func Setup(app *fiber.App, ctrl *Controllers) {
	v1 := app.Group("/api/v1")

	// auth (public)
	auth := v1.Group("/auth")
	auth.Post("/register", ctrl.Auth.Register)
	auth.Post("/login", ctrl.Auth.Login)
	auth.Post("/refresh", ctrl.Auth.Refresh)
	auth.Post("/verify-email", ctrl.Auth.VerifyEmail)
	auth.Post("/forgot-password", ctrl.Auth.ForgotPassword)
	auth.Post("/reset-password", ctrl.Auth.ResetPassword)

	// TODO: add other controller groups as implemented
}
