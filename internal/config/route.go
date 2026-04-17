package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// BootstrapConfig holds all initialized infrastructure dependencies.
// Passed to Bootstrap to wire repositories, usecases, and controllers together.
type BootstrapConfig struct {
	App      *fiber.App
	DB       *gorm.DB
	Redis    *redis.Client
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *Config
}

// Bootstrap wires all layers: repository → usecase → controller → route.
func Bootstrap(cfg *BootstrapConfig) {
	// TODO: inject repositories
	// userRepo := repository.NewUserRepository(cfg.DB, cfg.Log)
	// creatorRepo := repository.NewCreatorProfileRepository(cfg.DB, cfg.Log)
	// ...

	// TODO: inject usecases
	// authUseCase := usecase.NewAuthUseCase(cfg.Config, cfg.Log, cfg.Validate, userRepo)
	// ...

	// TODO: inject controllers
	// authController := handler.NewAuthController(authUseCase, cfg.Log)
	// ...

	// TODO: register routes
	// route.Setup(cfg.App, authController, ...)

	cfg.App.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
