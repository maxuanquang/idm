package app

import (
	"context"
	"syscall"

	"github.com/maxuanquang/idm/internal/handler/grpc"
	"github.com/maxuanquang/idm/internal/handler/http"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type Server struct {
	grpcServer grpc.Server
	httpServer http.Server
	logger     *zap.Logger
}

func NewServer(
	grpcServer grpc.Server,
	httpServer http.Server,
	logger *zap.Logger,
) (Server, error) {
	return Server{
		grpcServer: grpcServer,
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

func (s *Server) Start() {

	go func() {
		err := s.grpcServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("can not start GRPC Server")
	}()

	go func() {
		err := s.httpServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("can not start HTTP Server")
	}()

	utils.WaitForSignals(syscall.SIGINT, syscall.SIGTERM)
}
