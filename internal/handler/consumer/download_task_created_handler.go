package consumer

import (
	"context"

	"github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/logic"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type DownloadTaskCreatedHandler interface {
	Handle(ctx context.Context, event database.DownloadTask) error
}

func NewDownloadTaskCreatedHandler(
	downloadTaskLogic logic.DownloadTaskLogic,
	logger *zap.Logger,
) (DownloadTaskCreatedHandler, error) {
	return &downloadTaskCreatedHandler{
		downloadTaskLogic: downloadTaskLogic,
		logger:            logger,
	}, nil
}

type downloadTaskCreatedHandler struct {
	downloadTaskLogic logic.DownloadTaskLogic
	logger            *zap.Logger
}

// Handle implements DownloadTaskCreatedHandler.
func (d *downloadTaskCreatedHandler) Handle(ctx context.Context, event database.DownloadTask) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("event", event))

	logger.Info("download task received at handlerFunc")
	err := d.downloadTaskLogic.ExecuteDownloadTask(
		ctx,
		logic.ExecuteDownloadTaskInput{
			DownloadTaskID: event.DownloadTaskID,
		},
	)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to download event")
		return err
	}

	return nil
}
