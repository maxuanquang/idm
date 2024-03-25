package database

import (
	"context"

	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type TokenPublicKey struct {
	TokenPublicKeyID    uint64 `gorm:"column:token_public_key_id;primaryKey"`
	TokenPublicKeyValue []byte `gorm:"column:token_public_key_value"`
}

type TokenPublicKeyDataAccessor interface {
	CreatePublicKey(ctx context.Context, tokenPublicKey TokenPublicKey) (uint64, error)
	GetPublicKey(ctx context.Context, tokenPublicKeyID uint64) (TokenPublicKey, error)
}

func NewTokenPublicKeyDataAccessor(
	database Database,
	logger *zap.Logger,
) (TokenPublicKeyDataAccessor, error) {
	return &tokenPublicKeyDataAccessor{
		database: database,
		logger:   logger,
	}, nil
}

type tokenPublicKeyDataAccessor struct {
	database Database
	logger   *zap.Logger
}

// CreatePublicKey implements TokenPublicKeyDataAccessor.
func (t *tokenPublicKeyDataAccessor) CreatePublicKey(ctx context.Context, tokenPublicKey TokenPublicKey) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	var createdTokenPublicKey = TokenPublicKey{
		TokenPublicKeyValue: tokenPublicKey.TokenPublicKeyValue,
	}
	result := t.database.Create(&createdTokenPublicKey)

	if result.Error != nil {
		logger.Error("failed to create new token public key", zap.Error(result.Error))
		return 0, result.Error
	}

	return createdTokenPublicKey.TokenPublicKeyID, nil
}

// GetPublicKey implements TokenPublicKeyDataAccessor.
func (t *tokenPublicKeyDataAccessor) GetPublicKey(ctx context.Context, tokenPublicKeyID uint64) (TokenPublicKey, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	var foundTokenPublicKey TokenPublicKey
	result := t.database.First(&foundTokenPublicKey, tokenPublicKeyID)
	if result.Error != nil {
		logger.Error("failed to get public key", zap.Error(result.Error))
		return TokenPublicKey{}, result.Error
	}

	return foundTokenPublicKey, nil
}
