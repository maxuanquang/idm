package http

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/maxuanquang/idm/web"
	"go.uber.org/zap"
)

type SPAHandler http.Handler

func NewSPAHandler(logger *zap.Logger) SPAHandler {
	return &spaHandler{
		fileSystem: http.FS(web.StaticContent),
		fileServer: http.FileServer(http.FS(web.StaticContent)),
		logger:     logger,
	}
}

type spaHandler struct {
	fileSystem http.FileSystem
	fileServer http.Handler
	logger     *zap.Logger
}

// ServeHTTP implements SPAHandler.
func (s *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	fileName := filepath.Base(r.URL.Path)
	if strings.Contains(fileName, ".") {
		r.URL.Path = "dist/" + r.URL.Path
		r.RequestURI = "/dist" + r.RequestURI
		s.fileServer.ServeHTTP(w, r)
	}

	indexFile, err := s.fileSystem.Open("dist/index.html")
	if err != nil {
		s.logger.With(zap.Error(err)).Error("Failed to open index.html")
		http.Error(w, "could not open embedded file", http.StatusInternalServerError)
	}
	defer indexFile.Close()

	if _, err = io.Copy(w, indexFile); err != nil {
		s.logger.With(zap.Error(err)).Error("Failed to copy index.html")
		http.Error(w, "could not copy embedded file", http.StatusInternalServerError)
	}
}
