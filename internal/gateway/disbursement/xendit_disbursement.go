package disbursement

import (
	"context"
	"errors"

	"github.com/kreatip/kreatip-backend/internal/config"
	"github.com/sirupsen/logrus"
)

type DisburseRequest struct {
	ExternalID        string // idempotency key (withdrawal.id)
	BankCode          string
	AccountNumber     string
	AccountHolderName string
	Amount            int64
	Description       string
}

type DisburseResponse struct {
	ExternalID string
	Status     string // PENDING | COMPLETED | FAILED
	RawResponse string
}

type Gateway interface {
	Disburse(ctx context.Context, req *DisburseRequest) (*DisburseResponse, error)
}

type xenditDisbursement struct {
	cfg *config.XenditConfig
	log *logrus.Logger
}

func NewXenditDisbursement(cfg *config.XenditConfig, log *logrus.Logger) Gateway {
	return &xenditDisbursement{cfg: cfg, log: log}
}

func (g *xenditDisbursement) Disburse(ctx context.Context, req *DisburseRequest) (*DisburseResponse, error) {
	// TODO: call Xendit Disbursement API
	// https://developers.xendit.co/api-reference/#create-disbursement
	g.log.WithField("external_id", req.ExternalID).Info("creating xendit disbursement")
	return nil, errors.New("not implemented")
}
