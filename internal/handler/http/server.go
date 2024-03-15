package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gw "github.com/maxuanquang/idm/internal/generated/grpc/idm"
)

type Server interface {
	Start(ctx context.Context) error
}

func NewServer() Server {
	return &server{}
}

type server struct{}

func (s *server) Start(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterIdmServiceHandlerFromEndpoint(
		ctx,
		mux,
		":8080",
		opts,
	)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	fmt.Println("http server is running on :8081")
	return http.ListenAndServe(":8081", mux)
}
