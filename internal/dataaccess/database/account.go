package database

import (
	"context"

	"gorm.io/gorm"
)

type Account struct {
	AccountID   uint64 `gorm:"column:account_id"`
	AccountName string `gorm:"column:account_name"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account Account) (Account, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByName(ctx context.Context, name string) (Account, error)
}

func NewAccountDataAccessor(database *gorm.DB) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}

type accountDataAccessor struct {
	database *gorm.DB
}

// CreateAccount implements AccountDataAccessor.
func (a *accountDataAccessor) CreateAccount(ctx context.Context, account Account) (Account, error) {
	createdAccount := Account{
		AccountName: account.AccountName,
	}
	result := a.database.Create(&createdAccount)
	if result.Error != nil {
		return Account{}, result.Error
	}

	return createdAccount, nil
}

// GetAccountByID implements AccountDataAccessor.
func (a *accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	var foundAccount Account
	result := a.database.First(&foundAccount, id)
	if result.Error != nil {
		return Account{}, result.Error
	}

	return foundAccount, nil
}

// GetAccountByName implements AccountDataAccessor.
func (a *accountDataAccessor) GetAccountByName(ctx context.Context, name string) (Account, error) {
	var foundAccount Account
	result := a.database.Where("account_name = ?", name).First(&foundAccount)
	if result.Error != nil {
		return Account{}, result.Error
	}

	return foundAccount, nil
}
