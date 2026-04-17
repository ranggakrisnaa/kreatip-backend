package storage

import (
	"context"
	"errors"

	"github.com/kreatip/kreatip-backend/internal/config"
	"github.com/sirupsen/logrus"
)

type UploadRequest struct {
	Key         string // path in bucket, e.g. "avatars/user-id.jpg"
	Data        []byte
	ContentType string
}

type Gateway interface {
	Upload(ctx context.Context, req *UploadRequest) (publicURL string, err error)
	Delete(ctx context.Context, key string) error
}

type s3Gateway struct {
	cfg *config.StorageConfig
	log *logrus.Logger
}

func NewS3Gateway(cfg *config.StorageConfig, log *logrus.Logger) Gateway {
	return &s3Gateway{cfg: cfg, log: log}
}

func (g *s3Gateway) Upload(ctx context.Context, req *UploadRequest) (string, error) {
	// TODO: use aws-sdk-go-v2 or minio-go to upload
	g.log.WithField("key", req.Key).Info("uploading file")
	return "", errors.New("not implemented")
}

func (g *s3Gateway) Delete(ctx context.Context, key string) error {
	// TODO: implement delete
	return errors.New("not implemented")
}
