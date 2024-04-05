package configs

type ExecuteAllPendingDownloadTask struct {
	Schedule         string `yaml:"schedule"`
	ConcurrencyLimit int    `yaml:"concurrency_limit"`
}

type UpdateFailedDownloadTaskStatusToPending struct {
	Schedule string `yaml:"schedule"`
}

type Cron struct {
	ExecuteAllPendingDownloadTask           ExecuteAllPendingDownloadTask           `yaml:"execute_all_pending_download_task"`
	UpdateFailedDownloadTaskStatusToPending UpdateFailedDownloadTaskStatusToPending `yaml:"update_failed_download_task_status_to_pending"`
}
