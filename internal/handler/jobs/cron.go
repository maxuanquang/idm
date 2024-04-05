package jobs

import (
	"context"
	"syscall"

	"github.com/go-co-op/gocron/v2"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type Job interface {
	Run(context.Context) error
	GetSchedule() string
}

type Cron interface {
	Start(context.Context) error
}

func NewCron(
	executeAllPendingDownloadTaskJob ExecuteAllPendingDownloadTaskJob,
	updateFailedDownloadTaskStatusToPending UpdateFailedDownloadTaskStatusToPendingJob,
	logger *zap.Logger,
) (Cron, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.Error("failed to create scheduler", zap.Error(err))
		return nil, err
	}

	err = scheduleCronJobs(scheduler, logger, executeAllPendingDownloadTaskJob, updateFailedDownloadTaskStatusToPending)
	if err != nil {
		logger.Error("failed to schedule jobs", zap.Error(err))
		return nil, err
	}

	return &cron{
		logger:    logger,
		scheduler: scheduler,
	}, nil
}

type cron struct {
	logger    *zap.Logger
	scheduler gocron.Scheduler
}

// Start implements Cron.
func (c *cron) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, c.logger)

	c.scheduler.Start()
	defer func() {
		if err := c.scheduler.Shutdown(); err != nil {
			logger.Error("failed to shutdown scheduler", zap.Error(err))
		}
	}()

	utils.WaitForSignals(syscall.SIGINT, syscall.SIGTERM)
	return nil
}

func scheduleCronJobs(scheduler gocron.Scheduler, logger *zap.Logger, jobs ...Job) error {
	for _, job := range jobs {
		if _, err := scheduler.NewJob(
			gocron.CronJob(job.GetSchedule(), true),
			gocron.NewTask(func() {
				err := job.Run(context.Background())
				if err != nil {
					logger.Error("failed to run job", zap.Error(err))
				}
			}),
		); err != nil {
			logger.Error("failed to schedule job", zap.Error(err))
			return err
		}
	}

	return nil
}
