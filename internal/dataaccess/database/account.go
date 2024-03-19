package database

import (
	"context"
	"errors"

	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
)

type Account struct {
	AccountID   uint64 `gorm:"column:account_id;primaryKey"`
	AccountName string `gorm:"column:account_name"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account Account) (Account, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByName(ctx context.Context, name string) (Account, error)
}

func NewAccountDataAccessor(database Database, logger *zap.Logger) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   logger,
	}
}

type accountDataAccessor struct {
	database Database
	logger   *zap.Logger
}

// CreateAccount implements AccountDataAccessor.
func (a *accountDataAccessor) CreateAccount(ctx context.Context, account Account) (Account, error) {
	createdAccount := Account{
		AccountName: account.AccountName,
	}
	result := a.database.Create(&createdAccount)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return Account{}, nil
		}

		logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("account_name", account.AccountName))
		logger.Error("error creating account", zap.Error(result.Error))
		return Account{}, result.Error
	}

	return createdAccount, nil
}

// GetAccountByID implements AccountDataAccessor.
func (a *accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	var foundAccount Account
	result := a.database.First(&foundAccount, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Account{}, nil
		}

		logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("account_id", id))
		logger.Error("error getting account", zap.Error(result.Error))
		return Account{}, result.Error
	}

	return foundAccount, nil
}

// GetAccountByName implements AccountDataAccessor.
func (a *accountDataAccessor) GetAccountByName(ctx context.Context, name string) (Account, error) {
	var foundAccount Account
	result := a.database.Where("account_name = ?", name).First(&foundAccount)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Account{}, nil
		}

		logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("account_name", name))
		logger.Error("error getting account", zap.Error(result.Error))
		return Account{}, result.Error
	}

	return foundAccount, nil
}
