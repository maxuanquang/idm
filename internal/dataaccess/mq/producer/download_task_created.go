package producer

import (
	"context"
	"encoding/json"

	// "github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

const (
	MessageQueueDownloadTaskCreated = "download_task_created"
)

type DownloadTaskCreatedProducer interface {
	Produce(ctx context.Context, downloadTaskID uint64) error
}

func NewDownloadTaskCreatedProducer(client Client, logger *zap.Logger) (DownloadTaskCreatedProducer, error) {
	return &downloadTaskCreatedProducer{
		client: client,
		logger: logger,
	}, nil
}

type downloadTaskCreatedProducer struct {
	client Client
	logger *zap.Logger
}

// Produce implements DownloadTaskCreatedProducer.
func (d *downloadTaskCreatedProducer) Produce(ctx context.Context, downloadTaskID uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("download_task_id", downloadTaskID))

	payload, err := json.Marshal(downloadTaskID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to marshal event download task created")
		return err
	}

	err = d.client.Produce(ctx, MessageQueueDownloadTaskCreated, payload)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to produce message download task created")
		return err
	}

	return nil
}
