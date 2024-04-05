package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gammazero/workerpool"
	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/dataaccess/file"
	"github.com/maxuanquang/idm/internal/dataaccess/mq/producer"
	"github.com/maxuanquang/idm/internal/generated/grpc/idm"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	DownloadTaskMetadataKeyFileName = "file-name"
)

var (
	ErrPermissionDenied         = status.Error(codes.PermissionDenied, "permission denied")
	ErrDownloadTaskNotCompleted = status.Error(codes.FailedPrecondition, "download task not completed")
)

type CreateDownloadTaskInput struct {
	Token string
	Type  idm.DownloadType
	URL   string
}

type CreateDownloadTaskOutput struct {
	DownloadTask idm.DownloadTask
}

type GetDownloadTaskListInput struct {
	Token  string
	Offset uint64
	Limit  uint64
}

type GetDownloadTaskListOutput struct {
	DownloadTaskList       []*idm.DownloadTask
	TotalDownloadTaskCount uint64
}

type UpdateDownloadTaskInput struct {
	Token          string
	DownloadTaskID uint64
	DownloadStatus uint16
	Metadata       string
}

type UpdateDownloadTaskOutput struct {
	DownloadTask idm.DownloadTask
}

type DeleteDownloadTaskInput struct {
	Token          string
	DownloadTaskID uint64
}

type ExecuteDownloadTaskInput struct {
	DownloadTaskID uint64
}

type GetDownloadTaskFileInput struct {
	Token          string
	DownloadTaskID uint64
}

type GetDownloadTaskFileOutput struct {
	Reader io.ReadCloser
}

type DownloadTaskLogic interface {
	CreateDownloadTask(ctx context.Context, in CreateDownloadTaskInput) (CreateDownloadTaskOutput, error)
	GetDownloadTaskList(ctx context.Context, in GetDownloadTaskListInput) (GetDownloadTaskListOutput, error)
	UpdateDownloadTask(ctx context.Context, in UpdateDownloadTaskInput) (UpdateDownloadTaskOutput, error)
	UpdateFailedDownloadTaskStatusToPending(ctx context.Context) error
	DeleteDownloadTask(ctx context.Context, in DeleteDownloadTaskInput) error

	ExecuteDownloadTask(ctx context.Context, in ExecuteDownloadTaskInput) error
	ExecuteAllPendingDownloadTask(ctx context.Context) error

	GetDownloadTaskFile(ctx context.Context, in GetDownloadTaskFileInput) (GetDownloadTaskFileOutput, error)
}

func NewDownloadTaskLogic(
	tokenLogic TokenLogic,
	accountDataAccessor database.AccountDataAccessor,
	downloadTaskDataAccessor database.DownloadTaskDataAccessor,
	downloadTaskCreatedProducer producer.DownloadTaskCreatedProducer,
	fileClient file.Client,
	database database.Database,
	logger *zap.Logger,
	cronConfig configs.Cron,
) (DownloadTaskLogic, error) {
	return &downloadTaskLogic{
		tokenLogic:                  tokenLogic,
		accountDataAccessor:         accountDataAccessor,
		downloadTaskDataAccessor:    downloadTaskDataAccessor,
		downloadTaskCreatedProducer: downloadTaskCreatedProducer,
		fileClient:                  fileClient,
		database:                    database,
		logger:                      logger,
		cronConfig:                  cronConfig,
	}, nil
}

type downloadTaskLogic struct {
	tokenLogic                  TokenLogic
	accountDataAccessor         database.AccountDataAccessor
	downloadTaskDataAccessor    database.DownloadTaskDataAccessor
	downloadTaskCreatedProducer producer.DownloadTaskCreatedProducer
	fileClient                  file.Client
	database                    database.Database
	logger                      *zap.Logger
	cronConfig                  configs.Cron
}

// CreateDownloadTask implements DownloadTaskLogic.
func (d *downloadTaskLogic) CreateDownloadTask(ctx context.Context, in CreateDownloadTaskInput) (CreateDownloadTaskOutput, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("create_download_task_input", in))

	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account id and expire time from token")
		return CreateDownloadTaskOutput{}, status.Error(codes.Unauthenticated, "authentication token is invalid")
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account from database")
		return CreateDownloadTaskOutput{}, status.Error(codes.NotFound, "account not found")
	}

	var createdDownloadTask database.DownloadTask
	txErr := d.database.Transaction(func(tx *gorm.DB) error {
		var err error

		createdDownloadTask, err = d.downloadTaskDataAccessor.WithDatabaseTransaction(tx).CreateDownloadTask(ctx, database.DownloadTask{
			OfAccountID:    accountID,
			DownloadType:   uint16(in.Type),
			DownloadURL:    in.URL,
			DownloadStatus: uint16(idm.DownloadStatus_Pending),
			Metadata:       "{}",
		})
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to create download task")
			return err
		}

		err = d.downloadTaskCreatedProducer.Produce(ctx, createdDownloadTask.DownloadTaskID)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to produce message download task created")
			return err
		}

		return nil
	})
	if txErr != nil {
		return CreateDownloadTaskOutput{}, status.Error(codes.Internal, "failed to create download task")
	}

	return CreateDownloadTaskOutput{
		DownloadTask: idm.DownloadTask{
			Id: createdDownloadTask.DownloadTaskID,
			OfAccount: &idm.Account{
				Id:          account.AccountID,
				AccountName: account.AccountName,
			},
			DownloadType:   idm.DownloadType(createdDownloadTask.DownloadType),
			Url:            createdDownloadTask.DownloadURL,
			DownloadStatus: idm.DownloadStatus(createdDownloadTask.DownloadStatus),
			Metadata:       createdDownloadTask.Metadata,
		},
	}, nil
}

// GetDownloadTaskList implements DownloadTaskLogic.
func (d *downloadTaskLogic) GetDownloadTaskList(ctx context.Context, in GetDownloadTaskListInput) (GetDownloadTaskListOutput, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("get_download_task_list_input", in))

	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account id and expire time from token")
		return GetDownloadTaskListOutput{}, status.Error(codes.Unauthenticated, "authentication token is invalid")
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account from database")
		return GetDownloadTaskListOutput{}, status.Error(codes.NotFound, "account not found")
	}

	// Get the list of download tasks for the account from the data accessor.
	downloadTasks, err := d.downloadTaskDataAccessor.GetDownloadTaskListOfAccount(ctx, accountID, in.Offset, in.Limit)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task list from data accessor")
		return GetDownloadTaskListOutput{}, status.Error(codes.Internal, "failed to get download task list")
	}

	// Get the total count of download tasks for the account from the data accessor.
	totalDownloadTaskCount, err := d.downloadTaskDataAccessor.GetDownloadTaskCountOfAccount(ctx, accountID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get total download task count from data accessor")
		return GetDownloadTaskListOutput{}, status.Error(codes.Internal, "failed to get total download task count")
	}

	// Construct the output with the retrieved download tasks and total count.
	var outTaskList []*idm.DownloadTask
	for _, task := range downloadTasks {
		outTaskList = append(outTaskList, &idm.DownloadTask{
			Id: task.DownloadTaskID,
			OfAccount: &idm.Account{
				Id:          account.AccountID,
				AccountName: account.AccountName,
			},
			DownloadType:   idm.DownloadType(task.DownloadType),
			Url:            task.DownloadURL,
			DownloadStatus: idm.DownloadStatus(task.DownloadStatus),
			Metadata:       task.Metadata,
		})
	}
	output := GetDownloadTaskListOutput{
		DownloadTaskList:       outTaskList,
		TotalDownloadTaskCount: totalDownloadTaskCount,
	}

	return output, nil
}

// UpdateDownloadTask implements DownloadTaskLogic.
func (d *downloadTaskLogic) UpdateDownloadTask(ctx context.Context, in UpdateDownloadTaskInput) (UpdateDownloadTaskOutput, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("update_download_task_input", in))

	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account id and expire time from token")
		return UpdateDownloadTaskOutput{}, status.Error(codes.Unauthenticated, "authentication token is invalid")
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account from database")
		return UpdateDownloadTaskOutput{}, status.Error(codes.NotFound, "account not found")
	}

	downloadTask, err := d.downloadTaskDataAccessor.GetDownloadTask(ctx, in.DownloadTaskID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task from database")
		return UpdateDownloadTaskOutput{}, status.Error(codes.NotFound, "download task not found")
	}

	if accountID != downloadTask.OfAccountID {
		return UpdateDownloadTaskOutput{}, status.Error(codes.PermissionDenied, "user do not have permission to update download task")
	}

	// Implement the logic to update the download task based on the input parameters
	var updatedTask database.DownloadTask
	txErr := d.database.Transaction(func(tx *gorm.DB) error {
		err = d.downloadTaskDataAccessor.WithDatabaseTransaction(tx).UpdateDownloadTask(ctx, in.DownloadTaskID, in.DownloadStatus, in.Metadata)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to update download task")
			return err
		}

		updatedTask, err = d.downloadTaskDataAccessor.GetDownloadTaskForUpdate(ctx, in.DownloadTaskID)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		logger.With(zap.Error(txErr)).Error("transaction failed")
		return UpdateDownloadTaskOutput{}, status.Error(codes.Internal, txErr.Error())
	}

	// Return the updated download task in the output.
	return UpdateDownloadTaskOutput{
		DownloadTask: idm.DownloadTask{
			Id: updatedTask.DownloadTaskID,
			OfAccount: &idm.Account{
				Id:          account.AccountID,
				AccountName: account.AccountName,
			},
			DownloadType:   idm.DownloadType(updatedTask.DownloadType),
			Url:            updatedTask.DownloadURL,
			DownloadStatus: idm.DownloadStatus(updatedTask.DownloadStatus),
			Metadata:       updatedTask.Metadata,
		},
	}, nil
}

// UpdateFailedDownloadTaskStatusToPending implements DownloadTaskLogic.
func (d *downloadTaskLogic) UpdateFailedDownloadTaskStatusToPending(ctx context.Context) error {
	return d.downloadTaskDataAccessor.UpdateFailedDownloadTaskStatusToPending(ctx)
}

// DeleteDownloadTask implements DownloadTaskLogic.
func (d *downloadTaskLogic) DeleteDownloadTask(ctx context.Context, in DeleteDownloadTaskInput) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("delete_download_task_input", in))

	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account id and expire time from token")
		return status.Error(codes.Unauthenticated, "authentication token is invalid")
	}

	downloadTask, err := d.downloadTaskDataAccessor.GetDownloadTask(ctx, in.DownloadTaskID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task from database")
		return status.Error(codes.NotFound, "download task not found")
	}

	if accountID != downloadTask.OfAccountID {
		return status.Error(codes.PermissionDenied, "user do not have permission to delete download task")
	}

	err = d.downloadTaskDataAccessor.DeleteDownloadTask(ctx, in.DownloadTaskID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to delete download task")
		return status.Error(codes.Internal, "failed to delete download task")
	}

	return nil
}

// ExecuteDownloadTask implements DownloadTaskLogic.
func (d *downloadTaskLogic) ExecuteDownloadTask(ctx context.Context, in ExecuteDownloadTaskInput) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("execute_download_task_input", in))

	downloadTask, err := d.updateDownloadTaskStatusFromPendingToDownloading(ctx, in.DownloadTaskID)
	if err != nil {
		logger.With(zap.Error(err)).Error("can not update task status from pending to downloading")
		return err
	}

	// Create downloader
	var downloader Downloader

	switch downloadTask.DownloadType {
	case uint16(idm.DownloadType_HTTP):
		downloader, err = NewHTTPDownloader(downloadTask.DownloadURL, d.logger)
		if err != nil {
			logger.With(zap.Error(err)).Error("can not create http downloader")
			return err
		}
	default:
		logger.With(zap.Error(err)).Error("download type not supported")
		return err
	}

	fileName := fmt.Sprintf("%d", downloadTask.DownloadTaskID)
	fileWriteCloser, err := d.fileClient.Write(ctx, fileName)
	if err != nil {
		logger.With(zap.Error(err)).Error("can not create file writer")
		return err
	}
	defer fileWriteCloser.Close()

	metadata, err := downloader.Download(ctx, fileWriteCloser)
	if err != nil {
		logger.With(zap.Error(err)).Error("filed to download file")
		return err
	}

	// Update donwloadTask in database
	metadata[DownloadTaskMetadataKeyFileName] = fileName
	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		logger.With(zap.Error(err)).Error("can not marshal metadata")
		return err
	}
	err = d.downloadTaskDataAccessor.UpdateDownloadTask(ctx, downloadTask.DownloadTaskID, uint16(idm.DownloadStatus_Success), string(jsonMetadata))
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to update download task status to success")
		return err
	}

	logger.Info("download task executed successfully")

	return nil
}

// ExecuteAllPendingDownloadTask implements DownloadTaskLogic.
func (d *downloadTaskLogic) ExecuteAllPendingDownloadTask(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, d.logger)

	pendingDownloadTaskIDList, err := d.downloadTaskDataAccessor.GetPendingDownloadTaskIDList(ctx)
	if err != nil {
		return err
	}
	if len(pendingDownloadTaskIDList) == 0 {
		logger.Info("no pending download task found")
		return nil
	}

	logger.
		With(zap.Int("len(pending_download_task_id_list)", len(pendingDownloadTaskIDList))).
		Info("pending download task found")

	workerPool := workerpool.New(d.cronConfig.ExecuteAllPendingDownloadTask.ConcurrencyLimit)
	for _, id := range pendingDownloadTaskIDList {
		workerPool.Submit(func() {
			if executeDownloadTaskErr := d.ExecuteDownloadTask(ctx, ExecuteDownloadTaskInput{
				DownloadTaskID: id,
			}); executeDownloadTaskErr != nil {
				logger.
					With(zap.Uint64("download_task_id", id)).
					With(zap.Error(executeDownloadTaskErr)).
					Error("failed to execute download_task")
			}
		})
	}

	workerPool.StopWait()
	return nil
}

func (d *downloadTaskLogic) updateDownloadTaskStatusFromPendingToDownloading(ctx context.Context, downloadTaskID uint64) (database.DownloadTask, error) {

	var (
		downloadTask database.DownloadTask
		txErr        error
		err          error
	)

	txErr = d.database.Transaction(func(tx *gorm.DB) error {
		downloadTask, err = d.downloadTaskDataAccessor.WithDatabaseTransaction(tx).GetDownloadTaskForUpdate(ctx, downloadTaskID)
		if err != nil {
			d.logger.With(zap.Error(err)).Error("failed to get download task")
			return err
		}

		if downloadTask.DownloadStatus != uint16(idm.DownloadStatus_Pending) {
			d.logger.Error("download task not in pending status to download")
			return fmt.Errorf("download task not in pending status to download")
		}

		err = d.downloadTaskDataAccessor.WithDatabaseTransaction(tx).UpdateDownloadTask(ctx, downloadTaskID, uint16(idm.DownloadStatus_Downloading), "")
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return database.DownloadTask{}, txErr
	}

	return downloadTask, nil
}

// GetDownloadTaskFile implements DownloadTaskLogic.
func (d *downloadTaskLogic) GetDownloadTaskFile(ctx context.Context, in GetDownloadTaskFileInput) (GetDownloadTaskFileOutput, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("get_download_task_file_input", in))

	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		return GetDownloadTaskFileOutput{}, err
	}

	downloadTask, err := d.downloadTaskDataAccessor.GetDownloadTask(ctx, in.DownloadTaskID)
	if err != nil {
		return GetDownloadTaskFileOutput{}, err
	}

	if accountID != downloadTask.OfAccountID {
		return GetDownloadTaskFileOutput{}, ErrPermissionDenied
	}

	if downloadTask.DownloadStatus != uint16(idm.DownloadStatus_Success) {
		return GetDownloadTaskFileOutput{}, ErrDownloadTaskNotCompleted
	}

	var downloadTaskMetadata map[string]any
	err = json.Unmarshal([]byte(downloadTask.Metadata), &downloadTaskMetadata)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to unmarshal metadata")
		return GetDownloadTaskFileOutput{}, err
	}

	fileName, ok := downloadTaskMetadata[DownloadTaskMetadataKeyFileName]
	if !ok {
		logger.Error("file name not found in metadata")
		return GetDownloadTaskFileOutput{}, ErrDownloadTaskNotCompleted
	}

	readCloser, err := d.fileClient.Read(ctx, fileName.(string))
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to read file")
		return GetDownloadTaskFileOutput{}, err
	}

	return GetDownloadTaskFileOutput{
		Reader: readCloser,
	}, nil
}
