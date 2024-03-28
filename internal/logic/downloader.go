package logic

import (
	"context"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	HTTPMetadataKeyContentType    = "Content-Type"
	HTTPResponseHeaderContentType = "Content-Type"
)

type Downloader interface {
	Download(ctx context.Context, writer io.Writer) (map[string]any, error)
}

func NewHTTPDownloader(
	url string,
	logger *zap.Logger,
) (Downloader, error) {
	return &httpDownloader{
		url:    url,
		logger: logger,
	}, nil
}

func NewFTPDownloader() (Downloader, error) {
	panic("unimplemented")
}

type httpDownloader struct {
	url    string
	logger *zap.Logger
}

// Download implements Downloader.
func (h *httpDownloader) Download(ctx context.Context, writer io.Writer) (map[string]any, error) {
	logger := utils.LoggerWithContext(ctx, h.logger)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, http.NoBody)
	if err != nil {
		logger.With(zap.Error(err)).Error("can not create request")
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.With(zap.Error(err)).Error("can not do request")
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to copy http response body to file writer")
		return nil, err
	}

	metadata := map[string]any{
		HTTPMetadataKeyContentType: resp.Header.Get(HTTPResponseHeaderContentType),
	}

	return metadata, nil
}
