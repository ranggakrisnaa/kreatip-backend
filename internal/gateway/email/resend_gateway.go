package email

import (
	"context"
	"errors"

	"github.com/kreatip/kreatip-backend/internal/config"
	"github.com/sirupsen/logrus"
)

type SendRequest struct {
	To      []string
	Subject string
	HTML    string
}

type Gateway interface {
	Send(ctx context.Context, req *SendRequest) error
}

type resendGateway struct {
	cfg *config.EmailConfig
	log *logrus.Logger
}

func NewResendGateway(cfg *config.EmailConfig, log *logrus.Logger) Gateway {
	return &resendGateway{cfg: cfg, log: log}
}

func (g *resendGateway) Send(ctx context.Context, req *SendRequest) error {
	// TODO: call Resend API or fallback to SMTP (MailHog in dev)
	g.log.WithField("to", req.To).Info("sending email")
	return errors.New("not implemented")
}
