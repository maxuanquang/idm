package grpc

import (
	"context"

	idm "github.com/maxuanquang/idm/internal/generated/grpc/idm"
	"github.com/maxuanquang/idm/internal/logic"
)

func NewHandler(accountLogic logic.Account) idm.IdmServiceServer {
	return &Handler{
		accountLogic: accountLogic,
	}
}

type Handler struct {
	idm.UnimplementedIdmServiceServer
	accountLogic logic.Account
}

// CreateAccount implements idm.IdmServiceServer.
func (h *Handler) CreateAccount(ctx context.Context, in *idm.CreateAccountRequest) (*idm.CreateAccountResponse, error) {
	err := in.ValidateAll()
	if err != nil {
		return nil, responseError(err)
	}

	account, err := h.accountLogic.CreateAccount(ctx, logic.CreateAccountInput{
		AccountName: in.AccountName,
		Password:    in.Password,
	})
	if err != nil {
		return nil, responseError(err)
	}
	return &idm.CreateAccountResponse{
		AccountId: account.ID,
	}, nil
}

// CreateDownloadTask implements idm.IdmServiceServer.
func (h *Handler) CreateDownloadTask(ctx context.Context, in *idm.CreateDownloadTaskRequest) (*idm.CreateDownloadTaskResponse, error) {
	panic("unimplemented")
}

// CreateSession implements idm.IdmServiceServer.
func (h *Handler) CreateSession(ctx context.Context, in *idm.CreateSessionRequest) (*idm.CreateSessionResponse, error) {
	err := in.ValidateAll()
	if err != nil {
		return nil, responseError(err)
	}

	session, err := h.accountLogic.CreateSession(
		ctx,
		logic.CreateSessionInput{
			AccountName: in.AccountName,
			Password:    in.Password,
		},
	)
	if err != nil {
		return nil, responseError(err)
	}

	return &idm.CreateSessionResponse{
		Token: session.Token,
		Account: &idm.Account{
			Id:          session.AccountID,
			AccountName: session.AccountName,
		},
	}, nil
}

// DeleteDownloadTask implements idm.IdmServiceServer.
func (h *Handler) DeleteDownloadTask(ctx context.Context, in *idm.DeleteDownloadTaskRequest) (*idm.DeleteDownloadTaskResponse, error) {
	panic("unimplemented")
}

// GetDownloadTaskFile implements idm.IdmServiceServer.
func (h *Handler) GetDownloadTaskFile(*idm.GetDownloadTaskFileRequest, idm.IdmService_GetDownloadTaskFileServer) error {
	panic("unimplemented")
}

// GetDownloadTaskList implements idm.IdmServiceServer.
func (h *Handler) GetDownloadTaskList(ctx context.Context, in *idm.GetDownloadTaskListRequest) (*idm.GetDownloadTaskListResponse, error) {
	panic("unimplemented")
}

// UpdateDownloadTask implements idm.IdmServiceServer.
func (h *Handler) UpdateDownloadTask(ctx context.Context, in *idm.UpdateDownloadTaskRequest) (*idm.UpdateDownloadTaskResponse, error) {
	panic("unimplemented")
}
