package logic

import (
	"context"
	"errors"
	"fmt"

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
	AccountName uint64
	Password    string
}

type CreateSessionOutput struct {
	// TODO: Implement this
}

type Account interface {
	CreateAccount(ctx context.Context, in CreateAccountInput) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, in CreateSessionInput) (CreateSessionOutput, error)
}

func NewAccount(
	gormDatabase *gorm.DB,
	accountDataAccessor database.AccountDataAccessor,
	passwordDataAccessor database.AccountPasswordDataAccessor,
	hashLogic Hash,
) Account {
	return &account{
		gormDatabase:         gormDatabase,
		accountDataAccessor:  accountDataAccessor,
		passwordDataAccessor: passwordDataAccessor,
		hashLogic:            hashLogic,
	}
}

type account struct {
	gormDatabase         *gorm.DB
	accountDataAccessor  database.AccountDataAccessor
	passwordDataAccessor database.AccountPasswordDataAccessor
	hashLogic            Hash
}

// CreateAccount implements Account.
func (a *account) CreateAccount(ctx context.Context, in CreateAccountInput) (CreateAccountOutput, error) {
	var createAccountOutput CreateAccountOutput

	err := a.gormDatabase.Transaction(func(tx *gorm.DB) error {
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
	panic("unimplemented")
}
