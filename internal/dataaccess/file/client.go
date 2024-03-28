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
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

func NewS3Client(downloadConfig configs.Download, logger *zap.Logger) (Client, error) {

	minioClient, err := minio.New(
		downloadConfig.Address,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				downloadConfig.Username,
				downloadConfig.Password,
				"",
			),
		},
	)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create minio client")
		return nil, err
	}

	// Make a new bucket
	bucketName := downloadConfig.Bucket
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			logger.Info("bucket existed")
		} else {
			logger.With(zap.Error(err)).Error("failed to create bucket")
			return nil, err
		}
	}

	return &s3Client{
		minioClient: minioClient,
		logger:      logger,
		bucketName:  bucketName,
	}, nil
}

type s3Client struct {
	minioClient *minio.Client
	logger      *zap.Logger
	bucketName  string
}

// Read implements Client.
func (s *s3Client) Read(ctx context.Context, fileName string) (io.ReadCloser, error) {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("file_name", fileName))

	object, err := s.minioClient.GetObject(ctx, s.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get object")
		return nil, err
	}

	return object, nil
}

// Write implements Client.
func (s *s3Client) Write(ctx context.Context, fileName string) (io.WriteCloser, error) {
	return newS3ReadWriteCloser(
		ctx,
		s.minioClient,
		s.logger,
		s.bucketName,
		fileName,
	)
}

func newS3ReadWriteCloser(
	ctx context.Context,
	minioClient *minio.Client,
	logger *zap.Logger,
	bucketName string,
	fileName string,
) (*s3ReadWriteCloser, error) {
	logger = utils.LoggerWithContext(ctx, logger)
	readWriteCloser := &s3ReadWriteCloser{
		buffer:   make([]byte, 0),
		isClosed: false,
	}

	go func() {
		if _, err := minioClient.PutObject(
			ctx,
			bucketName,
			fileName,
			readWriteCloser,
			-1,
			minio.PutObjectOptions{},
		); err != nil {
			logger.With(zap.Error(err)).Error("failed to put object")
		}
	}()

	return readWriteCloser, nil
}

type s3ReadWriteCloser struct {
	buffer   []byte
	isClosed bool
}

func (s *s3ReadWriteCloser) Read(p []byte) (int, error) {
	if len(s.buffer) > 0 {
		readLength := copy(p, s.buffer)
		s.buffer = s.buffer[readLength:]
		return readLength, nil
	}

	if s.isClosed {
		return 0, io.EOF
	}

	return 0, nil
}

func (s *s3ReadWriteCloser) Write(p []byte) (int, error) {
	s.buffer = append(s.buffer, p...)
	return len(p), nil
}

func (s *s3ReadWriteCloser) Close() error {
	s.isClosed = true
	return nil
}
