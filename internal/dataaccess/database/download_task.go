package database

import (
	"gorm.io/gorm"
)

type DownloadTask struct {
	DownloadTaskID uint64 `gorm:"column:download_task_id;primaryKey"`
	OfAccountID    uint64 `gorm:"column:of_account_id"`
	DownloadType   int16  `gorm:"column:download_type"`
	DownloadURL    string `gorm:"column:download_url"`
	DownloadStatus int16  `gorm:"column:download_status"`
	Metadata       string `gorm:"column:metadata"`
}

type DownloadTaskDataAccessor interface{}

func NewDownloadTaskDataAccessor(database *gorm.DB) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{database: database}
}

type downloadTaskDataAccessor struct {
	database *gorm.DB
}
