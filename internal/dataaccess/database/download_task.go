package database

import (
	"context"
	"errors"

	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrDownloadTaskNotFound = errors.New("download task not found")
)

type DownloadTask struct {
	DownloadTaskID uint64 `gorm:"column:download_task_id;primaryKey"`
	OfAccountID    uint64 `gorm:"column:of_account_id"`
	DownloadType   int16  `gorm:"column:download_type"`
	DownloadURL    string `gorm:"column:download_url"`
	DownloadStatus int16  `gorm:"column:download_status"`
	Metadata       string `gorm:"column:metadata"`
}

type DownloadTaskDataAccessor interface {
	CreateDownloadTask(ctx context.Context, downloadTask DownloadTask) (DownloadTask, error)
	GetDownloadTask(ctx context.Context, downloadTaskID uint64) (DownloadTask, error)
	GetDownloadTaskListOfAccount(ctx context.Context, accountID, offset, limit uint64) ([]DownloadTask, error)
	GetDownloadTaskCountOfAccount(ctx context.Context, accountID uint64) (uint64, error)
	UpdateDownloadTask(ctx context.Context, downloadTask DownloadTask) error
	DeleteDownloadTask(ctx context.Context, downloadTaskID uint64) error
	WithDatabase(database Database) DownloadTaskDataAccessor
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
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("downloadTask", downloadTask))

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

	result := d.database.Model(&DownloadTask{}).Where("of_account_id = ?", accountID)
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("error getting download task count")
		return 0, result.Error
	}

	return uint64(result.RowsAffected), nil
}

// UpdateDownloadTask implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) UpdateDownloadTask(ctx context.Context, downloadTask DownloadTask) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("downloadTask", downloadTask))

	result := d.database.Save(&downloadTask)
	if result.Error != nil {
		logger.With(zap.Error(result.Error)).Error("failed to update download task")
		return result.Error
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

// WithDatabase implements DownloadTaskDataAccessor.
func (d *downloadTaskDataAccessor) WithDatabase(database Database) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
	}
}
