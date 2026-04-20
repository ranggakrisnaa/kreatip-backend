package consumer

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	emailgw "github.com/kreatip/kreatip-backend/internal/gateway/email"
	"github.com/kreatip/kreatip-backend/internal/model"
	tmpl "github.com/kreatip/kreatip-backend/internal/template"
	"github.com/sirupsen/logrus"
)

type EmailConsumer struct {
	Log     *logrus.Logger
	Gateway emailgw.Gateway
}

func NewEmailConsumer(log *logrus.Logger, gw emailgw.Gateway) *EmailConsumer {
	return &EmailConsumer{Log: log, Gateway: gw}
}

func (c *EmailConsumer) HandleVerification(ctx context.Context, t *asynq.Task) error {
	var p model.VerificationEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	html, err := tmpl.RenderVerification(p.Email, p.VerificationURL)
	if err != nil {
		c.Log.WithError(err).Error("email: render verification template")
		return err
	}

	return c.Gateway.Send(ctx, &emailgw.SendRequest{
		To:      []string{p.Email},
		Subject: "Verify your Kreatip account",
		HTML:    html,
	})
}
