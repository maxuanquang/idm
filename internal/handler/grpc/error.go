package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func clientResponseError(err error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Internal:
			return status.Error(codes.Internal, "Something went wrong")
		default:
			return err
		}
	}
	return err
}
