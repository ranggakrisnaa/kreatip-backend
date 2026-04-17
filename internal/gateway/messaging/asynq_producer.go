package messaging

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// Producer wraps the Asynq client for enqueuing background jobs.
type Producer struct {
	client *asynq.Client
	log    *logrus.Logger
}

func NewProducer(redisAddr string, log *logrus.Logger) *Producer {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &Producer{client: client, log: log}
}

func (p *Producer) Enqueue(taskType string, payload any, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(taskType, data, opts...)
	info, err := p.client.Enqueue(task)
	if err != nil {
		p.log.WithError(err).Errorf("failed to enqueue task: %s", taskType)
		return err
	}
	p.log.Debugf("enqueued task %s id=%s queue=%s", taskType, info.ID, info.Queue)
	return nil
}
