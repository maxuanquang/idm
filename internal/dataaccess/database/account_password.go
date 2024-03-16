package database

import (
	"context"
)

type AccountPassword struct {
	OfAccountID uint64 `gorm:"column:of_account_id;primaryKey"`
	Hashed      string `gorm:"column:hashed"`
}

type AccountPasswordDataAccessor interface {
	CreatePassword(ctx context.Context, ofAccountID uint64, hashedPassword string) (AccountPassword, error)
	ComparePassword(ctx context.Context, ofAccountID uint64, password string) error
}

type accountPasswordDataAccessor struct {
	database Database
}

func NewAccountPasswordDataAccessor(database Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{database: database}
}

// ComparePassword implements AccountPasswordDataAccessor.
func (a *accountPasswordDataAccessor) ComparePassword(ctx context.Context, ofAccountID uint64, password string) error {
	panic("unimplemented")
}

// CreatePassword implements AccountPasswordDataAccessor.
func (a *accountPasswordDataAccessor) CreatePassword(ctx context.Context, ofAccountID uint64, hashedPassword string) (AccountPassword, error) {
	var createdPassword = AccountPassword{
		OfAccountID: ofAccountID,
		Hashed:      hashedPassword,
	}
	result := a.database.Create(&createdPassword)
	if result.Error != nil {
		return AccountPassword{}, result.Error
	}

	return createdPassword, nil
}
