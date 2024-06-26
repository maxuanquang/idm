package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/maxuanquang/idm/internal/dataaccess/cache"
	"github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateAccountInput struct {
	AccountName string
	Password    string
}

type CreateAccountOutput struct {
	ID          uint64
	AccountName string
}

type CreateSessionInput struct {
	AccountName string
	Password    string
}

type CreateSessionOutput struct {
	Token       string
	ExpiresAt   time.Time
	AccountID   uint64
	AccountName string
}

type DeleteSessionInput struct {
	Token string
}

type DeleteSessionOutput struct {
	ExpiredToken string
}

type AccountLogic interface {
	CreateAccount(ctx context.Context, in CreateAccountInput) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, in CreateSessionInput) (CreateSessionOutput, error)
	DeleteSession(ctx context.Context, in DeleteSessionInput) error
}

func NewAccountLogic(
	database database.Database,
	accountDataAccessor database.AccountDataAccessor,
	passwordDataAccessor database.AccountPasswordDataAccessor,
	hashLogic HashLogic,
	tokenLogic TokenLogic,
	takenAccountNameCache cache.TakenAccountName,
	logger *zap.Logger,
) AccountLogic {
	return &accountLogic{
		database:              database,
		accountDataAccessor:   accountDataAccessor,
		passwordDataAccessor:  passwordDataAccessor,
		hashLogic:             hashLogic,
		tokenLogic:            tokenLogic,
		takenAccountNameCache: takenAccountNameCache,
		logger:                logger,
	}
}

type accountLogic struct {
	database              database.Database
	accountDataAccessor   database.AccountDataAccessor
	passwordDataAccessor  database.AccountPasswordDataAccessor
	hashLogic             HashLogic
	tokenLogic            TokenLogic
	takenAccountNameCache cache.TakenAccountName
	logger                *zap.Logger
}

// CreateAccount implements Account.
func (a *accountLogic) CreateAccount(ctx context.Context, in CreateAccountInput) (CreateAccountOutput, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Any("create_account_input", in))

	taken, err := a.isAccountNameTaken(ctx, in.AccountName)
	if err != nil {
		logger.Error("failed to check if account name taken", zap.Error(err))
		return CreateAccountOutput{}, status.Error(codes.Internal, "failed to check if account name taken")
	}
	if taken {
		return CreateAccountOutput{}, status.Error(codes.AlreadyExists, "account name already exists")
	}

	var createAccountOutput CreateAccountOutput
	txErr := a.database.Transaction(func(tx *gorm.DB) error {
		// create createdAccount
		createdAccount, err := a.accountDataAccessor.WithDatabaseTransaction(tx).CreateAccount(
			ctx,
			database.Account{
				AccountName: in.AccountName,
			},
		)
		if err != nil {
			return fmt.Errorf("error creating account: %w", err)
		}

		// create password
		hashedPassword, err := a.hashLogic.HashPassword(ctx, in.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		_, err = a.passwordDataAccessor.WithDatabaseTransaction(tx).CreatePassword(ctx, createdAccount.AccountID, hashedPassword)
		if err != nil {
			return fmt.Errorf("error creating password: %w", err)
		}

		createAccountOutput.ID = createdAccount.AccountID
		createAccountOutput.AccountName = createdAccount.AccountName
		return nil
	})
	if txErr != nil {
		logger.With(zap.Error(txErr)).Error("create account transaction failed")
		return CreateAccountOutput{}, status.Error(codes.Internal, txErr.Error())
	}

	err = a.takenAccountNameCache.Add(ctx, createAccountOutput.AccountName)
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to set account name into taken set in cache")
	}

	return createAccountOutput, nil

}

// CreateSession implements Account.
func (a *accountLogic) CreateSession(ctx context.Context, in CreateSessionInput) (CreateSessionOutput, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Any("create_session_input", in))

	foundAccount, err := a.accountDataAccessor.GetAccountByName(ctx, in.AccountName)
	if err != nil {
		logger.Error("failed to get account by name", zap.Error(err))
		return CreateSessionOutput{}, status.Error(codes.Internal, "error getting account")
	}
	if foundAccount.AccountID == 0 {
		return CreateSessionOutput{}, status.Error(codes.NotFound, "wrong account name or password")
	}

	foundPassword, err := a.passwordDataAccessor.GetPassword(ctx, foundAccount.AccountID)
	if err != nil {
		logger.Error("failed to get account password", zap.Error(err))
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to get account password")
	}

	matched, err := a.hashLogic.IsHashEqual(ctx, in.Password, foundPassword.Hashed)
	if err != nil {
		logger.Error("failed comparing password", zap.Error(err))
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed comparing password")
	}
	if !matched {
		return CreateSessionOutput{}, status.Error(codes.NotFound, "wrong account name or password")
	}

	stringToken, expiresAt, err := a.tokenLogic.CreateTokenString(ctx, foundAccount.AccountID)
	if err != nil {
		logger.Error("failed to create token", zap.Error(err))
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to create token")
	}

	return CreateSessionOutput{
		Token:       stringToken,
		ExpiresAt:   expiresAt,
		AccountID:   foundAccount.AccountID,
		AccountName: foundAccount.AccountName,
	}, nil
}

// DeleteSession implements AccountLogic.
func (a *accountLogic) DeleteSession(ctx context.Context, in DeleteSessionInput) error {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Any("delete_session_input", in))

	_, _, err := a.tokenLogic.GetAccountIDAndExpireTime(ctx, in.Token)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to find account from token")
		return status.Error(codes.NotFound, "account not found")
	}

	return nil
}

func (a *accountLogic) isAccountNameTaken(ctx context.Context, accountName string) (bool, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("account_name", accountName))

	// Check cache
	taken, err := a.takenAccountNameCache.Has(ctx, accountName)
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to get account name from taken set in cache, will fall back to database")
	}
	if taken {
		return true, nil
	}

	// check account name taken
	foundAccount, err := a.accountDataAccessor.GetAccountByName(ctx, accountName)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account name from database")
		return false, err
	}
	if foundAccount.AccountID == 0 {
		return false, nil
	}

	// add missed taken name to cache
	err = a.takenAccountNameCache.Add(ctx, accountName)
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to add missed taken account name to cache")
	}

	return true, nil
}
