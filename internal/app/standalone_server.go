package app

import (
	"context"
	"syscall"

	"github.com/maxuanquang/idm/internal/handler/consumer"
	"github.com/maxuanquang/idm/internal/handler/grpc"
	"github.com/maxuanquang/idm/internal/handler/http"
	"github.com/maxuanquang/idm/internal/handler/jobs"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type StandaloneServer struct {
	grpcServer grpc.Server
	httpServer http.Server
	mqConsumer consumer.RootConsumer
	cron       jobs.Cron
	logger     *zap.Logger
}

func NewStandaloneServer(
	grpcServer grpc.Server,
	httpServer http.Server,
	mqConsumer consumer.RootConsumer,
	cron jobs.Cron,
	logger *zap.Logger,
) (StandaloneServer, error) {
	return StandaloneServer{
		grpcServer: grpcServer,
		httpServer: httpServer,
		mqConsumer: mqConsumer,
		cron:       cron,
		logger:     logger,
	}, nil
}

func (s *StandaloneServer) Start() {

	go func() {
		err := s.grpcServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("can not start GRPC Server")
	}()

	go func() {
		err := s.httpServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("can not start HTTP Server")
	}()

	go func() {
		err := s.mqConsumer.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("can not start message queue consumer")
	}()

	go func() {
		err := s.cron.Start(context.Background())
		s.logger.With(zap.Error(err)).Error("can not start cron jobs")
	}()

	utils.WaitForSignals(syscall.SIGINT, syscall.SIGTERM)
}
