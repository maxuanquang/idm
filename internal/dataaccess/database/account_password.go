package database

import (
	"context"

	"go.uber.org/zap"
)

type AccountPassword struct {
	OfAccountID uint64 `gorm:"column:of_account_id;primaryKey"`
	Hashed      string `gorm:"column:hashed"`
}

type AccountPasswordDataAccessor interface {
	CreatePassword(ctx context.Context, ofAccountID uint64, hashedPassword string) (AccountPassword, error)
	GetPassword(ctx context.Context, ofAccountID uint64) (AccountPassword, error)
	WithDatabaseTransaction(database Database) AccountPasswordDataAccessor
}

type accountPasswordDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewAccountPasswordDataAccessor(
	database Database,
	logger *zap.Logger,
) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
		logger:   logger,
	}
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

// WithDatabaseTransaction implements AccountPasswordDataAccessor.
func (a *accountPasswordDataAccessor) WithDatabaseTransaction(database Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
		logger:   a.logger,
	}
}
