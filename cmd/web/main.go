// @title           Kreatip API
// @version         1.0
// @description     Donation platform for creators & streamers
// @termsOfService  https://kreatip.id/terms

// @contact.name   Kreatip Team
// @contact.email  dev@kreatip.id

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJ..."

package main

import (
	"fmt"

	_ "github.com/kreatip/kreatip-backend/docs" // swag generated docs
	"github.com/kreatip/kreatip-backend/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.json")
	viper.SetEnvPrefix("KREATIP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("error reading config: %v", err)
	}

	cfg := new(config.Config)
	if err := viper.Unmarshal(cfg); err != nil {
		logrus.Fatalf("error unmarshaling config: %v", err)
	}

	log := config.NewLogger(cfg)
	log.Infof("starting %s in %s mode", cfg.App.Name, cfg.App.Env)

	db := config.NewDatabase(cfg, log)
	rdb := config.NewRedis(cfg, log)
	validate := config.NewValidator()
	app := config.NewFiber(cfg)

	config.Bootstrap(&config.BootstrapConfig{
		App:      app,
		DB:       db,
		Redis:    rdb,
		Log:      log,
		Validate: validate,
		Config:   cfg,
	})

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Infof("server listening on %s", addr)
	log.Infof("swagger UI: http://localhost%s/docs/", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
