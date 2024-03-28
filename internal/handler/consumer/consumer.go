package consumer

import (
	"context"
	"encoding/json"

	"github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/dataaccess/mq/consumer"
	"github.com/maxuanquang/idm/internal/dataaccess/mq/producer"
	"go.uber.org/zap"
)

type RootConsumer interface {
	Start(ctx context.Context) error
}

func NewRootConsumer(
	downloadTaskCreatedHandler DownloadTaskCreatedHandler,
	mqConsumer consumer.Consumer,
	logger *zap.Logger,
) RootConsumer {
	return &rootConsumer{
		downloadTaskCreatedHandler: downloadTaskCreatedHandler,
		mqConsumer:                 mqConsumer,
		logger:                     logger,
	}
}

type rootConsumer struct {
	downloadTaskCreatedHandler DownloadTaskCreatedHandler
	mqConsumer                 consumer.Consumer
	logger                     *zap.Logger
}

// Start implements RootConsumer.
func (r *rootConsumer) Start(ctx context.Context) error {
	r.mqConsumer.RegisterHandler(
		producer.MessageQueueDownloadTaskCreated,
		func(ctx context.Context, payload []byte) error {

			var downloadTaskID uint64

			err := json.Unmarshal(payload, &downloadTaskID)
			if err != nil {
				return err
			}

			return r.downloadTaskCreatedHandler.Handle(ctx, database.DownloadTask{
				DownloadTaskID: downloadTaskID,
			})
		},
	)

	return r.mqConsumer.Start(ctx)
}
