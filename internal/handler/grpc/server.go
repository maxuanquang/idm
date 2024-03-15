package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/maxuanquang/idm/internal/generated/grpc/idm"
	"google.golang.org/grpc"
)

type Server interface {
	Start(ctx context.Context) error
}

func NewServer(handler idm.IdmServiceServer) Server {
	return &server{
		handler: handler,
	}
}

type server struct {
	handler idm.IdmServiceServer
}

// Start implements Server.
func (s *server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		return err
	}
	defer listener.Close()

	var opts = []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	idm.RegisterIdmServiceServer(server, s.handler)
	return server.Serve(listener)
}
