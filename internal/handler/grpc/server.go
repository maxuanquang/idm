package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/generated/grpc/idm"
	"google.golang.org/grpc"
)

type Server interface {
	Start(ctx context.Context) error
}

func NewServer(grpcConfig configs.GRPC, handler idm.IdmServiceServer) Server {
	return &server{
		grpcConfig: grpcConfig,
		handler:    handler,
	}
}

type server struct {
	grpcConfig configs.GRPC
	handler    idm.IdmServiceServer
}

// Start implements Server.
func (s *server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.grpcConfig.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	var opts = []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	idm.RegisterIdmServiceServer(server, s.handler)

	fmt.Printf("gRPC server is running on %s\n", s.grpcConfig.Address)
	return server.Serve(listener)
}
