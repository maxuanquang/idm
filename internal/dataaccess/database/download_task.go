package database

import (
	"context"
	"errors"

	"github.com/maxuanquang/idm/internal/generated/grpc/idm"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrDownloadTaskNotFound = errors.New("download task not found")
)

type DownloadTask struct {
	DownloadTaskID uint64 `gorm:"column:download_task_id;primaryKey"`
	OfAccountID    uint64 `gorm:"column:of_account_id"`
	DownloadType   uint16 `gorm:"column:download_type"`
	DownloadURL    string `gorm:"column:download_url"`
	DownloadStatus uint16 `gorm:"column:download_status"`
	Metadata       string `gorm:"column:metadata"`
}

type DownloadTaskDataAccessor interface {
	CreateDownloadTask(ctx context.Context, downloadTask DownloadTask) (DownloadTask, error)
	GetDownloadTask(ctx context.Context, downloadTaskID uint64) (DownloadTask, error)
	GetDownloadTaskForUpdate(ctx context.Context, downloadTaskID uint64) (DownloadTask, error)
	GetPendingDownloadTaskIDList(ctx context.Context) ([]uint64, error)
	GetDownloadTaskListOfAccount(ctx context.Context, accountID, offset, limit uint64) ([]DownloadTask, error)
	GetDownloadTaskCountOfAccount(ctx context.Context, accountID uint64) (uint64, error)
	UpdateDownloadTask(ctx context.Context, downloadTaskID uint64, downloadStatus uint16, metadata string) error
	UpdateFailedDownloadTaskStatusToPending(ctx context.Context) error
	DeleteDownloadTask(ctx context.Context, downloadTaskID uint64) error
	WithDatabaseTransaction(database Database) DownloadTaskDataAccessor
}

func NewDownloadTaskDataAccessor(database Database, logger *zap.Logger) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
		logger:   logger,
	}
}

type downloadTaskDataAccessor struct {
	database Database
	logger   *zap.Logger
}

// CreateDownloadTask implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) CreateDownloadTask(ctx context.Context, downloadTask DownloadTask) (DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("downloadTask", nil))

	var createdDownloadTask = downloadTask
	createdDownloadTask.DownloadTaskID = 0

	result := d.database.Create(&createdDownloadTask)
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("error creating download task")
		return DownloadTask{}, result.Error
	}

	return createdDownloadTask, nil
}

// GetDownloadTask implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) GetDownloadTask(ctx context.Context, downloadTaskID uint64) (DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("downloadTaskID", downloadTaskID))

	var downloadTask DownloadTask
	result := d.database.Where("download_task_id = ?", downloadTaskID).First(&downloadTask)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return DownloadTask{}, ErrDownloadTaskNotFound
		}

		logger.With(zap.Error(result.Error)).Error("error getting download task")
		return DownloadTask{}, result.Error
	}

	return downloadTask, nil
}

func (d *downloadTaskDataAccessor) GetPendingDownloadTaskIDList(ctx context.Context) ([]uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger)

	var downloadTaskIDs []uint64
	result := d.database.Model(&DownloadTask{}).Where("download_status = ?", uint16(idm.DownloadStatus_Pending)).Pluck("download_task_id", &downloadTaskIDs)
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("error getting pending download task id list")
		return nil, result.Error
	}

	return downloadTaskIDs, nil
}

// GetDownloadTaskForUpdate implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) GetDownloadTaskForUpdate(ctx context.Context, downloadTaskID uint64) (DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("downloadTaskID", downloadTaskID))

	var downloadTask DownloadTask
	result := d.database.Clauses(clause.Locking{Strength: "UPDATE"}).Where("download_task_id = ?", downloadTaskID).First(&downloadTask)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return DownloadTask{}, ErrDownloadTaskNotFound
		}

		logger.With(zap.Error(result.Error)).Error("error getting download task")
		return DownloadTask{}, result.Error
	}

	return downloadTask, nil
}

// GetDownloadTaskListOfAccount implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) GetDownloadTaskListOfAccount(ctx context.Context, accountID, offset, limit uint64) ([]DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("accountID", accountID)).With(zap.Uint64("offset", offset)).With(zap.Uint64("limit", limit))

	var downloadTasks []DownloadTask
	result := d.database.Where("of_account_id = ?", accountID).Offset(int(offset)).Limit(int(limit)).Find(&downloadTasks)
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("error getting download task list")
		return nil, result.Error
	}

	return downloadTasks, nil
}

// GetDownloadTaskCountOfAccount implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) GetDownloadTaskCountOfAccount(ctx context.Context, accountID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("accountID", accountID))

	var downloadTasks []DownloadTask
	result := d.database.Where("of_account_id = ?", accountID).Find(&downloadTasks)
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("error getting download task count")
		return 0, result.Error
	}

	return uint64(result.RowsAffected), nil
}

// UpdateDownloadTask implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) UpdateDownloadTask(ctx context.Context, downloadTaskID uint64, downloadStatus uint16, metadata string) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("downloadTaskID", downloadTaskID)).With(zap.Uint16("downloadStatus", downloadStatus)).With(zap.String("metadata", metadata))

	downloadTask, err := d.GetDownloadTaskForUpdate(ctx, downloadTaskID)
	if err != nil {
		return err
	}

	if downloadStatus != 0 {
		downloadTask.DownloadStatus = downloadStatus
	}
	if metadata != "" {
		downloadTask.Metadata = metadata
	}

	result := d.database.Save(&downloadTask)
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("failed to update download task")
		return result.Error
	}

	return nil
}

func (d *downloadTaskDataAccessor) UpdateFailedDownloadTaskStatusToPending(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, d.logger)

	result := d.database.Model(&DownloadTask{}).Where("download_status = ?", uint16(idm.DownloadStatus_Failed)).Update("download_status", uint16(idm.DownloadStatus_Pending))
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("failed to update failed download task status to pending")
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Info("no failed task found")
	}
	return nil
}

// DeleteDownloadTask implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) DeleteDownloadTask(ctx context.Context, downloadTaskID uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("downloadTaskID", downloadTaskID))

	result := d.database.Delete(&DownloadTask{}, downloadTaskID)
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("error deleting download task")
		return result.Error
	}

	return nil
}

// WithDatabaseTransaction implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) WithDatabaseTransaction(database Database) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
		logger:   d.logger,
	}
}
