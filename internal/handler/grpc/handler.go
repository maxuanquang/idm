package grpc

import (
	"context"
	// "strings"

	idm "github.com/maxuanquang/idm/internal/generated/grpc/idm"
	"github.com/maxuanquang/idm/internal/logic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	AuthTokenMetadataName         = "IDM_AUTH"
	GRPCGatewayCookieMetadataName = "grpcgateway-cookie"
)

func NewHandler(
	accountLogic logic.AccountLogic,
	downloadTaskLogic logic.DownloadTaskLogic,
) idm.IdmServiceServer {
	return &Handler{
		accountLogic:      accountLogic,
		downloadTaskLogic: downloadTaskLogic,
	}
}

type Handler struct {
	idm.UnimplementedIdmServiceServer
	accountLogic      logic.AccountLogic
	downloadTaskLogic logic.DownloadTaskLogic
}

func (h *Handler) getAuthTokenFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	authTokenValues := md.Get(AuthTokenMetadataName)
	if len(authTokenValues) == 0 {
		return ""
	}

	return authTokenValues[0]
}

// CreateAccount implements idm.IdmServiceServer.
func (h *Handler) CreateAccount(ctx context.Context, in *idm.CreateAccountRequest) (*idm.CreateAccountResponse, error) {
	account, err := h.accountLogic.CreateAccount(ctx, logic.CreateAccountInput{
		AccountName: in.AccountName,
		Password:    in.Password,
	})
	if err != nil {
		return nil, clientResponseError(err)
	}
	return &idm.CreateAccountResponse{
		AccountId: account.ID,
	}, nil
}

// CreateSession implements idm.IdmServiceServer.
func (h *Handler) CreateSession(ctx context.Context, in *idm.CreateSessionRequest) (*idm.CreateSessionResponse, error) {
	session, err := h.accountLogic.CreateSession(
		ctx,
		logic.CreateSessionInput{
			AccountName: in.AccountName,
			Password:    in.Password,
		},
	)
	if err != nil {
		return nil, clientResponseError(err)
	}

	err = grpc.SendHeader(ctx, metadata.Pairs(AuthTokenMetadataName, session.Token))
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &idm.CreateSessionResponse{
		Account: &idm.Account{
			Id:          session.AccountID,
			AccountName: session.AccountName,
		},
	}, nil
}

// DeleteSession implements idm.IdmServiceServer.
func (h *Handler) DeleteSession(ctx context.Context, in *idm.DeleteSessionRequest) (*idm.DeleteSessionResponse, error) {
	err := h.accountLogic.DeleteSession(
		ctx,
		logic.DeleteSessionInput{
			Token: h.getAuthTokenFromMetadata(ctx),
		},
	)
	if err != nil {
		return nil, clientResponseError(err)
	}

	err = grpc.SendHeader(ctx, metadata.Pairs(AuthTokenMetadataName, ""))
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &idm.DeleteSessionResponse{}, nil
}

// CreateDownloadTask implements idm.IdmServiceServer.
func (h *Handler) CreateDownloadTask(ctx context.Context, in *idm.CreateDownloadTaskRequest) (*idm.CreateDownloadTaskResponse, error) {
	out, err := h.downloadTaskLogic.CreateDownloadTask(ctx, logic.CreateDownloadTaskInput{
		Token: h.getAuthTokenFromMetadata(ctx),
		Type:  in.DownloadType,
		URL:   in.Url,
	})
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &idm.CreateDownloadTaskResponse{
		DownloadTask: &out.DownloadTask,
	}, nil
}

// GetDownloadTaskFile implements idm.IdmServiceServer.
func (h *Handler) GetDownloadTaskFile(*idm.GetDownloadTaskFileRequest, idm.IdmService_GetDownloadTaskFileServer) error {
	panic("unimplemented")
}

// GetDownloadTaskList implements idm.IdmServiceServer.
func (h *Handler) GetDownloadTaskList(ctx context.Context, in *idm.GetDownloadTaskListRequest) (*idm.GetDownloadTaskListResponse, error) {
	out, err := h.downloadTaskLogic.GetDownloadTaskList(ctx, logic.GetDownloadTaskListInput{
		Token:  h.getAuthTokenFromMetadata(ctx),
		Offset: in.Offset,
		Limit:  in.Limit,
	})
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &idm.GetDownloadTaskListResponse{
		DownloadTaskList:       out.DownloadTaskList,
		TotalDownloadTaskCount: out.TotalDownloadTaskCount,
	}, nil
}

// UpdateDownloadTask implements idm.IdmServiceServer.
func (h *Handler) UpdateDownloadTask(ctx context.Context, in *idm.UpdateDownloadTaskRequest) (*idm.UpdateDownloadTaskResponse, error) {
	var (
		downloadStatus uint16 = 0
		metadata       string = ""
	)

	if in.DownloadStatus != nil {
		downloadStatus = uint16(*in.DownloadStatus)
	}
	if in.Metadata != nil {
		metadata = *in.Metadata
	}

	out, err := h.downloadTaskLogic.UpdateDownloadTask(ctx, logic.UpdateDownloadTaskInput{
		Token:          h.getAuthTokenFromMetadata(ctx),
		DownloadTaskID: in.DownloadTaskId,
		DownloadStatus: downloadStatus,
		Metadata:       metadata,
	})
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &idm.UpdateDownloadTaskResponse{
		DownloadTask: &out.DownloadTask,
	}, nil
}

// DeleteDownloadTask implements idm.IdmServiceServer.
func (h *Handler) DeleteDownloadTask(ctx context.Context, in *idm.DeleteDownloadTaskRequest) (*idm.DeleteDownloadTaskResponse, error) {
	err := h.downloadTaskLogic.DeleteDownloadTask(ctx, logic.DeleteDownloadTaskInput{
		Token:          h.getAuthTokenFromMetadata(ctx),
		DownloadTaskID: in.DownloadTaskId,
	})
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &idm.DeleteDownloadTaskResponse{}, nil
}
