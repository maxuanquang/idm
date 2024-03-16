package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maxuanquang/idm/internal/dataaccess/database"
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
) Account {
	return &account{
		database:             database,
		accountDataAccessor:  accountDataAccessor,
		passwordDataAccessor: passwordDataAccessor,
		hashLogic:            hashLogic,
		tokenLogic:           tokenLogic,
	}
}

type account struct {
	database             database.Database
	accountDataAccessor  database.AccountDataAccessor
	passwordDataAccessor database.AccountPasswordDataAccessor
	hashLogic            Hash
	tokenLogic           Token
}

// CreateAccount implements Account.
func (a *account) CreateAccount(ctx context.Context, in CreateAccountInput) (CreateAccountOutput, error) {
	var createAccountOutput CreateAccountOutput

	err := a.database.Transaction(func(tx *gorm.DB) error {
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

	stringToken, expiresAt, err := a.tokenLogic.CreateToken(ctx, foundAccount.AccountID)
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
