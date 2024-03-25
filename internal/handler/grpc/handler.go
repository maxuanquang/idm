package grpc

import (
	"context"

	idm "github.com/maxuanquang/idm/internal/generated/grpc/idm"
	"github.com/maxuanquang/idm/internal/logic"
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

// CreateAccount implements idm.IdmServiceServer.
func (h *Handler) CreateAccount(ctx context.Context, in *idm.CreateAccountRequest) (*idm.CreateAccountResponse, error) {
	err := in.ValidateAll()
	if err != nil {
		return nil, clientResponseError(err)
	}

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
	err := in.ValidateAll()
	if err != nil {
		return nil, clientResponseError(err)
	}

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

	return &idm.CreateSessionResponse{
		Token: session.Token,
		Account: &idm.Account{
			Id:          session.AccountID,
			AccountName: session.AccountName,
		},
	}, nil
}

// CreateDownloadTask implements idm.IdmServiceServer.
func (h *Handler) CreateDownloadTask(ctx context.Context, in *idm.CreateDownloadTaskRequest) (*idm.CreateDownloadTaskResponse, error) {
	err := in.ValidateAll()
	if err != nil {
		return nil, clientResponseError(err)
	}

	out, err := h.downloadTaskLogic.CreateDownloadTask(ctx, logic.CreateDownloadTaskInput{
		Token: in.Token,
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
	err := in.ValidateAll()
	if err != nil {
		return nil, clientResponseError(err)
	}

	out, err := h.downloadTaskLogic.GetDownloadTaskList(ctx, logic.GetDownloadTaskListInput{
		Token:  in.Token,
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
	err := in.ValidateAll()
	if err != nil {
		return nil, clientResponseError(err)
	}

	out, err := h.downloadTaskLogic.UpdateDownloadTask(ctx, logic.UpdateDownloadTaskInput{
		Token:          in.Token,
		DownloadTaskID: in.DownloadTaskId,
		URL:            in.Url,
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
	err := in.ValidateAll()
	if err != nil {
		return nil, clientResponseError(err)
	}

	err = h.downloadTaskLogic.DeleteDownloadTask(ctx, logic.DeleteDownloadTaskInput{
		Token:          in.Token,
		DownloadTaskID: in.DownloadTaskId,
	})
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &idm.DeleteDownloadTaskResponse{}, nil
}
