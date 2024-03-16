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
	GetPassword(ctx context.Context, ofAccountID uint64) (AccountPassword, error)
}

type accountPasswordDataAccessor struct {
	database Database
}

func NewAccountPasswordDataAccessor(database Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{database: database}
}

// GetPassword implements AccountPasswordDataAccessor.
func (a *accountPasswordDataAccessor) GetPassword(ctx context.Context, ofAccountID uint64) (AccountPassword, error) {
	var foundPassword AccountPassword
	result := a.database.Where("of_account_id = ?", ofAccountID).First(&foundPassword)
	if result.Error != nil {
		return AccountPassword{}, result.Error
	}

	return foundPassword, nil
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
