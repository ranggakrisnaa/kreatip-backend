package config

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(cfg *Config) *logrus.Logger {
	log := logrus.New()

	if cfg.App.Env == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
		log.SetLevel(logrus.DebugLevel)
	}

	return log
}
