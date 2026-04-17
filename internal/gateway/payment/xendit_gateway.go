package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/kreatip/kreatip-backend/internal/config"
	"github.com/sirupsen/logrus"
)

type xenditGateway struct {
	cfg *config.XenditConfig
	log *logrus.Logger
}

func NewXenditGateway(cfg *config.XenditConfig, log *logrus.Logger) Gateway {
	return &xenditGateway{cfg: cfg, log: log}
}

func (g *xenditGateway) CreateCharge(ctx context.Context, req *CreateChargeRequest) (*CreateChargeResponse, error) {
	// TODO: call Xendit invoice/eWallet/QRIS API based on req.Method
	// https://developers.xendit.co/api-reference
	g.log.WithField("external_id", req.ExternalID).Info("creating xendit charge")
	return nil, errors.New("not implemented")
}

func (g *xenditGateway) VerifyWebhookSignature(payload []byte, signature string) error {
	mac := hmac.New(sha256.New, []byte(g.cfg.WebhookToken))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return errors.New("invalid webhook signature")
	}
	return nil
}
