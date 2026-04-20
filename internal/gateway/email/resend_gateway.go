package email

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

// Config holds email gateway settings — defined here to avoid importing internal/config.
type Config struct {
	Provider string
	APIKey   string
	From     string
	SMTPHost string
	SMTPPort int
}

type SendRequest struct {
	To      []string
	Subject string
	HTML    string
}

type Gateway interface {
	Send(ctx context.Context, req *SendRequest) error
}

type resendGateway struct {
	cfg *Config
	log *logrus.Logger
}

func NewResendGateway(cfg *Config, log *logrus.Logger) Gateway {
	return &resendGateway{cfg: cfg, log: log}
}

// NewGateway returns the right implementation based on Provider field.
// "resend" → Resend API, anything else → SMTP (MailHog in dev).
func NewGateway(cfg *Config, log *logrus.Logger) Gateway {
	if cfg.Provider == "resend" {
		return NewResendGateway(cfg, log)
	}
	return NewSMTPGateway(cfg, log)
}

func (g *resendGateway) Send(ctx context.Context, req *SendRequest) error {
	// TODO: implement Resend HTTP API call
	g.log.WithField("to", req.To).Info("resend: sending email")
	return errors.New("resend gateway not implemented")
}
