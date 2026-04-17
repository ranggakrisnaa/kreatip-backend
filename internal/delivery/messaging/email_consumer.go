package consumer

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	TypeEmailReceipt          = "email:receipt"
	TypeEmailCreatorNotif     = "email:creator_notif"
	TypeEmailWithdrawalStatus = "email:withdrawal_status"
	TypeEmailVerification     = "email:verification"
	TypeEmailPasswordReset    = "email:password_reset"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	RefID   string `json:"ref_id"`   // donation or withdrawal ID
	RefType string `json:"ref_type"` // donation | withdrawal
}

type EmailConsumer struct {
	Log *logrus.Logger
	// TODO: inject email gateway
}

func NewEmailConsumer(log *logrus.Logger) *EmailConsumer {
	return &EmailConsumer{Log: log}
}

func (c *EmailConsumer) HandleReceipt(ctx context.Context, t *asynq.Task) error {
	var p EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	c.Log.Infof("sending receipt email to %s", p.To)
	// TODO: render template + call email gateway
	return nil
}

func (c *EmailConsumer) HandleVerification(ctx context.Context, t *asynq.Task) error {
	var p EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	c.Log.Infof("sending verification email to %s", p.To)
	// TODO: render template + call email gateway
	return nil
}
