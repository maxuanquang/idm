package servemuxoption

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const (
	GRPCMetadataPrefix = "Grpc-Metadata-"
)

func WithAuthCookieToAuthMetadata(authCookieName string, authTokenMetadataName string) runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
		cookie, err := r.Cookie(authCookieName)
		if err != nil {
			return make(metadata.MD)
		}
		return metadata.New(
			map[string]string{
				authTokenMetadataName: cookie.Value,
			},
		)
	})
}

func WithAuthMetadataToAuthCookie(authCookieName string, authTokenMetadataName string, expiresInDuration time.Duration) runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(
		func(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
			md, ok := runtime.ServerMetadataFromContext(ctx)
			if !ok {
				return nil
			}

			authTokenMetadataValues := md.HeaderMD.Get(authTokenMetadataName)
			if len(authTokenMetadataValues) == 0 {
				return nil
			}
			
			http.SetCookie(w, &http.Cookie{
				Name:     authCookieName,
				Value:    authTokenMetadataValues[0],
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(expiresInDuration),
			})
			return nil
		},
	)
}
