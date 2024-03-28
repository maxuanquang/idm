package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type Client interface {
	Write(ctx context.Context, fileName string) (io.WriteCloser, error)
	Read(ctx context.Context, fileName string) (io.ReadCloser, error)
}

func NewClient(downloadConfig configs.Download, logger *zap.Logger) (Client, error) {
	switch downloadConfig.Mode {
	case configs.DownloadModeLocal:
		return NewLocalClient(downloadConfig, logger)
	case configs.DownloadModeS3:
		return NewS3Client(downloadConfig, logger)
	default:
		return nil, fmt.Errorf("unsupported download mode: %s", downloadConfig.Mode)
	}
}

func NewLocalClient(downloadConfig configs.Download, logger *zap.Logger) (Client, error) {
	if err := os.MkdirAll(downloadConfig.DownloadDirectory, os.ModeDir); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
	}

	return &localClient{
		downloadDirectory: downloadConfig.DownloadDirectory,
		logger:            logger,
	}, nil
}

type localClient struct {
	downloadDirectory string
	logger            *zap.Logger
}

// Read implements Client.
func (l *localClient) Read(ctx context.Context, fileName string) (io.ReadCloser, error) {
	logger := utils.LoggerWithContext(ctx, l.logger).With(zap.String("file_name", fileName))

	filePath := path.Join(l.downloadDirectory, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		logger.With(zap.Error(err)).Error("can not open file")
		return nil, err
	}

	return file, nil
}

// Write implements Client.
func (l *localClient) Write(ctx context.Context, fileName string) (io.WriteCloser, error) {
	logger := utils.LoggerWithContext(ctx, l.logger).With(zap.String("file_name", fileName))

	filePath := path.Join(l.downloadDirectory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		logger.With(zap.Error(err)).Error("can not open file")
		return nil, err
	}

	return file, nil
}

// TODO: Implement this
func NewS3Client(downloadConfig configs.Download, logger *zap.Logger) (Client, error) {
	return nil, nil
}
