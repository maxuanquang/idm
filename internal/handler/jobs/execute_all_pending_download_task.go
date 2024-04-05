package jobs

import (
	"context"

	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/logic"
)

type ExecuteAllPendingDownloadTaskJob interface {
	Run(ctx context.Context) error
	GetSchedule() string
}

func NewExecuteAllPendingDownloadTaskJob(
	downloadTaskLogic logic.DownloadTaskLogic,
	cronConfig configs.Cron,
) ExecuteAllPendingDownloadTaskJob {
	return &executeAllPendingDownloadTaskJob{
		cronConfig:        cronConfig,
		downloadTaskLogic: downloadTaskLogic,
	}
}

type executeAllPendingDownloadTaskJob struct {
	downloadTaskLogic logic.DownloadTaskLogic
	cronConfig        configs.Cron
}

// GetSchedule implements ExecuteAllPendingDownloadTaskJob.
func (e *executeAllPendingDownloadTaskJob) GetSchedule() string {
	return e.cronConfig.ExecuteAllPendingDownloadTask.Schedule
}

// Run implements ExecuteAllPendingDownloadTaskJob.
func (e *executeAllPendingDownloadTaskJob) Run(ctx context.Context) error {
	return e.downloadTaskLogic.ExecuteAllPendingDownloadTask(ctx)
}
