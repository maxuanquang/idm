package jobs

import (
	"context"

	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/logic"
)

type UpdateFailedDownloadTaskStatusToPendingJob interface {
	Run(ctx context.Context) error
	GetSchedule() string
}

func NewUpdateFailedDownloadTaskStatusToPendingJob(
	downloadTaskLogic logic.DownloadTaskLogic,
	cronConfig configs.Cron,
) UpdateFailedDownloadTaskStatusToPendingJob {
	return &updateFailedDownloadTaskStatusToPendingJob{
		downloadTaskLogic: downloadTaskLogic,
		cronConfig:        cronConfig,
	}
}

type updateFailedDownloadTaskStatusToPendingJob struct {
	downloadTaskLogic logic.DownloadTaskLogic
	cronConfig        configs.Cron
}

// GetSchedule implements UpdateFailedDownloadTaskStatusToPendingJob.
func (e *updateFailedDownloadTaskStatusToPendingJob) GetSchedule() string {
	return e.cronConfig.UpdateFailedDownloadTaskStatusToPending.Schedule
}

// Run implements updateFailedDownloadTaskStatusToPendingJob.
func (e *updateFailedDownloadTaskStatusToPendingJob) Run(ctx context.Context) error {
	return e.downloadTaskLogic.UpdateFailedDownloadTaskStatusToPending(ctx)
}
