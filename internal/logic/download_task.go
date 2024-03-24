package logic

import (
	"context"

	"github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/generated/grpc/idm"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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
	DownloadTaskList       []idm.DownloadTask
	TotalDownloadTaskCount uint64
}

type UpdateDownloadTaskInput struct {
	Token          string
	DownloadTaskID uint64
	URL            string
}

type UpdateDownloadTaskOutput struct {
	DownloadTask idm.DownloadTask
}

type DeleteDownloadTaskInput struct {
	Token          string
	DownloadTaskID uint64
}

type DownloadTaskLogic interface {
	CreateDownloadTask(ctx context.Context, in CreateDownloadTaskInput) (CreateDownloadTaskOutput, error)
	GetDownloadTaskList(ctx context.Context, in GetDownloadTaskListInput) (GetDownloadTaskListOutput, error)
	UpdateDownloadTask(ctx context.Context, in UpdateDownloadTaskInput) (UpdateDownloadTaskOutput, error)
	DeleteDownloadTask(ctx context.Context, in DeleteDownloadTaskInput) error
}

func NewDownloadTaskLogic() (DownloadTaskLogic, error) {
	return &downloadTaskLogic{}, nil
}

type downloadTaskLogic struct {
	tokenLogic               TokenLogic
	downloadTaskDataAccessor database.DownloadTaskDataAccessor
	// messageProducer
	// messageConsumer
	database database.Database
	logger   *zap.Logger
}

// CreateDownloadTask implements DownloadTaskLogic.
func (d *downloadTaskLogic) CreateDownloadTask(ctx context.Context, in CreateDownloadTaskInput) (CreateDownloadTaskOutput, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("create_download_task_input", in))
	
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account id and expire time from token")
		return CreateDownloadTaskOutput{}, status.Error(codes.Unauthenticated, "authentication token is invalid")
	}
	
	var createdDownloadTask database.DownloadTask
	txErr := d.database.Transaction(func(tx *gorm.DB) error {
		var err error
		
		createdDownloadTask, err = d.downloadTaskDataAccessor.WithDatabase(tx).CreateDownloadTask(ctx, database.DownloadTask{
			OfAccountID:    accountID,
			DownloadType:   int16(in.Type),
			DownloadURL:    in.URL,
			DownloadStatus: int16(idm.DownloadStatus_Pending),
			Metadata:       "{}",
		})
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to create download task")
			return err
		}
		
		// TODO: add operations from message queue
		// 1. After done creating task in db, send a message to queue
		// 2. Downloader will receive message from queue, do downloading

		return nil
	})
	if txErr != nil {
		return CreateDownloadTaskOutput{}, status.Error(codes.Internal, "failed to create download task")
	}

	return CreateDownloadTaskOutput{
		DownloadTask: idm.DownloadTask{
			Id:             createdDownloadTask.DownloadTaskID,
			OfAccount:      nil,
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
	var outTaskList []idm.DownloadTask
	for _, task := range downloadTasks {
		outTaskList = append(outTaskList, idm.DownloadTask{
			Id:             task.DownloadTaskID,
			OfAccount:      nil,
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

	_, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account id and expire time from token")
		return UpdateDownloadTaskOutput{}, status.Error(codes.Unauthenticated, "authentication token is invalid")
	}

	// Implement the logic to update the download task based on the input parameters
	var updatedTask database.DownloadTask
	txErr := d.database.Transaction(func(tx *gorm.DB) error {
		downloadTask, err := d.downloadTaskDataAccessor.WithDatabase(tx).GetDownloadTask(ctx, in.DownloadTaskID)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to get download task")
			return err
		}

		downloadTask.DownloadURL = in.URL
		err = d.downloadTaskDataAccessor.WithDatabase(tx).UpdateDownloadTask(ctx, downloadTask)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to update download task")
			return err
		}

		updatedTask = downloadTask
		return nil
	})
	if txErr != nil {
		return UpdateDownloadTaskOutput{}, txErr
	}

	// Return the updated download task in the output.
	return UpdateDownloadTaskOutput{
		DownloadTask: idm.DownloadTask{
			Id:             updatedTask.DownloadTaskID,
			OfAccount:      nil,
			DownloadType:   idm.DownloadType(updatedTask.DownloadType),
			Url:            updatedTask.DownloadURL,
			DownloadStatus: idm.DownloadStatus(updatedTask.DownloadStatus),
			Metadata:       updatedTask.Metadata,
		},
	}, nil
}

// DeleteDownloadTask implements DownloadTaskLogic.
func (d *downloadTaskLogic) DeleteDownloadTask(ctx context.Context, in DeleteDownloadTaskInput) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("delete_download_task_input", in))

	_, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account id and expire time from token")
		return status.Error(codes.Unauthenticated, "authentication token is invalid")
	}

	err = d.downloadTaskDataAccessor.DeleteDownloadTask(ctx, in.DownloadTaskID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to delete download task")
		return status.Error(codes.Internal, "failed to delete download task")
	}

	return nil
}