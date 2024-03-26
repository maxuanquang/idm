package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/maxuanquang/idm/internal/configs"
	gw "github.com/maxuanquang/idm/internal/generated/grpc/idm"
	grpcHandler "github.com/maxuanquang/idm/internal/handler/grpc"
	"github.com/maxuanquang/idm/internal/handler/http/servemuxoption"
)

const (
	AuthCookieName = "IDM_AUTH"
)

type Server interface {
	Start(ctx context.Context) error
}

func NewServer(
	httpConfig configs.HTTP,
	grpcConfig configs.GRPC,
	authConfig configs.Auth,
	logger *zap.Logger,
) Server {
	return &server{
		httpConfig: httpConfig,
		grpcConfig: grpcConfig,
		authConfig: authConfig,
		logger:     logger,
	}
}

type server struct {
	httpConfig configs.HTTP
	grpcConfig configs.GRPC
	authConfig configs.Auth
	logger     *zap.Logger
}

func (s *server) Start(ctx context.Context) error {
	mux := runtime.NewServeMux(
		// servemuxoption.WithAuthCookieToAuthMetadata(AuthCookieName, grpcHandler.AuthTokenMetadataName),
		servemuxoption.WithAuthMetadataToAuthCookie(AuthCookieName, grpcHandler.AuthTokenMetadataName, s.authConfig.Token.GetTokenDuration()),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterIdmServiceHandlerFromEndpoint(
		ctx,
		mux,
		s.grpcConfig.Address,
		opts,
	)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	fmt.Printf("http server is running on %s\n", s.httpConfig.Address)
	return http.ListenAndServe(s.httpConfig.Address, mux)
}
