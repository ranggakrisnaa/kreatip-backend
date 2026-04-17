package consumer

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const TypeProcessDonation = "donation:process"

type DonationPayload struct {
	DonationID string `json:"donation_id"`
}

type DonationConsumer struct {
	Log *logrus.Logger
	// TODO: inject usecase
}

func NewDonationConsumer(log *logrus.Logger) *DonationConsumer {
	return &DonationConsumer{Log: log}
}

func (c *DonationConsumer) Handle(ctx context.Context, t *asynq.Task) error {
	var payload DonationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	c.Log.Infof("processing donation job: %s", payload.DonationID)
	// TODO: implement
	return nil
}
