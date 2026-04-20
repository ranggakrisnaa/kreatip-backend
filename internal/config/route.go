package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
	handler "github.com/kreatip/kreatip-backend/internal/delivery/http"
	"github.com/kreatip/kreatip-backend/internal/delivery/http/route"
	"github.com/kreatip/kreatip-backend/internal/repository"
	"github.com/kreatip/kreatip-backend/internal/usecase/impl"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// BootstrapConfig holds all initialized infrastructure dependencies.
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
	// repositories
	userRepo := repository.NewUserRepository(cfg.DB, cfg.Log)
	walletRepo := repository.NewWalletRepository(cfg.DB, cfg.Log)

	// usecases
	authUC := impl.NewAuthUseCase(cfg.DB, cfg.Log, cfg.Validate, userRepo, walletRepo)

	// controllers
	authCtrl := handler.NewAuthController(authUC, cfg.Log)

	// healthcheck
	cfg.App.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// swagger UI — disable di production
	if cfg.Config.App.Env != "production" {
		cfg.App.Get("/docs/*", fiberSwagger.HandlerDefault)
	}

	// routes
	route.Setup(cfg.App, &route.Controllers{
		Auth: authCtrl,
	})
}
