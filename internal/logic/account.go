package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maxuanquang/idm/internal/dataaccess/cache"
	"github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

type Account interface {
	CreateAccount(ctx context.Context, in CreateAccountInput) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, in CreateSessionInput) (CreateSessionOutput, error)
}

func NewAccount(
	database database.Database,
	accountDataAccessor database.AccountDataAccessor,
	passwordDataAccessor database.AccountPasswordDataAccessor,
	hashLogic Hash,
	tokenLogic Token,
	takenAccountNameCache cache.TakenAccountName,
	logger *zap.Logger,
) Account {
	return &account{
		database:              database,
		accountDataAccessor:   accountDataAccessor,
		passwordDataAccessor:  passwordDataAccessor,
		hashLogic:             hashLogic,
		tokenLogic:            tokenLogic,
		takenAccountNameCache: takenAccountNameCache,
		logger:                logger,
	}
}

type account struct {
	database              database.Database
	accountDataAccessor   database.AccountDataAccessor
	passwordDataAccessor  database.AccountPasswordDataAccessor
	hashLogic             Hash
	tokenLogic            Token
	takenAccountNameCache cache.TakenAccountName
	logger                *zap.Logger
}

// CreateAccount implements Account.
func (a *account) CreateAccount(ctx context.Context, in CreateAccountInput) (CreateAccountOutput, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("account_name", in.AccountName))

	var createAccountOutput CreateAccountOutput

	exists, err := a.takenAccountNameCache.Has(ctx, in.AccountName)
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to get account name from taken set in cache, will fall back to database")
	}
	if exists {
		return CreateAccountOutput{}, errors.New("check account existence: account name existed")
	}

	err = a.database.Transaction(func(tx *gorm.DB) error {
		// check account name taken
		_, err := a.accountDataAccessor.GetAccountByName(ctx, in.AccountName)
		if err == nil {
			return errors.New("check account existence: account name existed")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check account existence: %w", err)
		}

		// create createdAccount
		createdAccount, err := a.accountDataAccessor.CreateAccount(
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

		_, err = a.passwordDataAccessor.CreatePassword(ctx, createdAccount.AccountID, hashedPassword)
		if err != nil {
			return fmt.Errorf("error creating password: %w", err)
		}

		createAccountOutput.ID = createdAccount.AccountID
		createAccountOutput.AccountName = createdAccount.AccountName
		return nil
	})
	if err != nil {
		return CreateAccountOutput{}, err
	}

	err = a.takenAccountNameCache.Add(ctx, createAccountOutput.AccountName)
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to set account name into taken set in cache")
	}

	return createAccountOutput, nil

}

// CreateSession implements Account.
func (a *account) CreateSession(ctx context.Context, in CreateSessionInput) (CreateSessionOutput, error) {
	foundAccount, err := a.accountDataAccessor.GetAccountByName(ctx, in.AccountName)
	if err != nil {
		return CreateSessionOutput{}, fmt.Errorf("account name does not exist: %w", err)
	}

	foundPassword, err := a.passwordDataAccessor.GetPassword(ctx, foundAccount.AccountID)
	if err != nil {
		return CreateSessionOutput{}, fmt.Errorf("password does not exist: %w", err)
	}

	matched, err := a.hashLogic.IsHashEqual(ctx, in.Password, foundPassword.Hashed)
	if err != nil || !matched {
		return CreateSessionOutput{}, fmt.Errorf("wrong account name or password: %w", err)
	}

	stringToken, expiresAt, err := a.tokenLogic.CreateTokenString(ctx, foundAccount.AccountID)
	if err != nil {
		return CreateSessionOutput{}, fmt.Errorf("can not create token: %w", err)
	}

	return CreateSessionOutput{
		Token:       stringToken,
		ExpiresAt:   expiresAt,
		AccountID:   foundAccount.AccountID,
		AccountName: foundAccount.AccountName,
	}, nil
}
