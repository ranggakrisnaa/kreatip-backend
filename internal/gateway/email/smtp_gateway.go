package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/sirupsen/logrus"
)

type smtpGateway struct {
	cfg *Config
	log *logrus.Logger
}

func NewSMTPGateway(cfg *Config, log *logrus.Logger) Gateway {
	return &smtpGateway{cfg: cfg, log: log}
}

func (g *smtpGateway) Send(_ context.Context, req *SendRequest) error {
	addr := fmt.Sprintf("%s:%d", g.cfg.SMTPHost, g.cfg.SMTPPort)

	header := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		g.cfg.From,
		strings.Join(req.To, ", "),
		req.Subject,
	)

	msg := []byte(header + req.HTML)

	// MailHog (dev) tidak perlu auth — nil auth
	err := smtp.SendMail(addr, nil, g.cfg.From, req.To, msg)
	if err != nil {
		g.log.WithError(err).WithField("to", req.To).Error("smtp: failed to send email")
		return err
	}

	g.log.WithField("to", req.To).Info("smtp: email sent")
	return nil
}
