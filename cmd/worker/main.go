package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
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
	log.Infof("starting worker in %s mode", cfg.App.Env)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Addr, DB: cfg.Redis.DB},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	// TODO: register task handlers here
	// mux.HandleFunc(tasks.TypeEmailReceipt, handlers.HandleEmailReceipt)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("worker error: %v", err)
		}
	}()

	log.Info("worker started, waiting for tasks...")
	<-quit

	log.Info("shutting down worker...")
	srv.Shutdown()
}
