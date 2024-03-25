package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/dataaccess/cache"
	"github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type TokenLogic interface {
	CreateTokenString(ctx context.Context, accountID uint64) (string, time.Time, error)
	GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error)
	WithDatabase(database database.Database) TokenLogic
}

func NewTokenLogic(
	accountDataAccessor database.AccountDataAccessor,
	tokenPublicKeyDataAccessor database.TokenPublicKeyDataAccessor,
	logger *zap.Logger,
	authConfig configs.Auth,
	tokenPublicKeyCache cache.TokenPublicKey,
) (TokenLogic, error) {

	rsaKeyPair, err := generateRSAKeyPair(int(authConfig.Token.RS512KeyPairBitSize))
	if err != nil {
		logger.Error("failed to genereate RSA key pair", zap.Error(err))
		return nil, err
	}

	publicKeyBytes, err := pemEncodePublicKey(&rsaKeyPair.PublicKey)
	if err != nil {
		logger.Error("failed to encode public key in pem format", zap.Error(err))
		return nil, err
	}

	tokenPublicKeyID, err := tokenPublicKeyDataAccessor.CreatePublicKey(
		context.Background(),
		database.TokenPublicKey{
			TokenPublicKeyValue: publicKeyBytes,
		})
	if err != nil {
		logger.Error("failed to create token public key", zap.Error(err))
		return nil, err
	}
	err = tokenPublicKeyCache.Set(context.Background(), fmt.Sprint(tokenPublicKeyID), publicKeyBytes)
	if err != nil {
		logger.With(zap.Error(err)).Warn("can not set token public key in cache")
	}

	return &tokenLogic{
		accountDataAccessor:        accountDataAccessor,
		tokenPublicKeyDataAccessor: tokenPublicKeyDataAccessor,
		logger:                     logger,
		authConfig:                 authConfig,
		tokenPublicKeyID:           tokenPublicKeyID,
		tokenPrivateKeyValue:       rsaKeyPair,
		tokenPublicKeyCache:        tokenPublicKeyCache,
	}, nil
}

type tokenLogic struct {
	accountDataAccessor        database.AccountDataAccessor
	tokenPublicKeyDataAccessor database.TokenPublicKeyDataAccessor
	logger                     *zap.Logger
	authConfig                 configs.Auth
	tokenPublicKeyID           uint64
	tokenPrivateKeyValue       *rsa.PrivateKey
	tokenPublicKeyCache        cache.TokenPublicKey
}

// GetAccountIDAndExpireTime implements Token.
func (t *tokenLogic) GetAccountIDAndExpireTime(ctx context.Context, tokenString string) (uint64, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	keyFunc := func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Error("unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			logger.Error("cannot get token's claims")
			return nil, errors.New("cannot get token's claims")
		}

		tokenPublicKeyID, ok := claims["kid"].(float64)
		if !ok {
			logger.Error("cannot get token's kid claim")
			return nil, errors.New("cannot get token's kid claim")
		}

		tokenPublicKeyValue, err := t.getJWTPublicKeyValue(ctx, uint64(tokenPublicKeyID))
		if err != nil {
			logger.Error("cannot get public key's value")
			return nil, err
		}

		return tokenPublicKeyValue, nil
	}
	parsedToken, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		logger.Error("cannot parse token", zap.Error(err))
		return 0, time.Time{}, err
	}

	if !parsedToken.Valid {
		logger.Error("invalid token")
		return 0, time.Time{}, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("cannot get token's claims")
		return 0, time.Time{}, errors.New("cannot get token's claims")
	}

	accountID, ok := claims["sub"].(float64)
	if !ok {
		logger.Error("cannot get token's sub claim")
		return 0, time.Time{}, errors.New("cannot get token's sub claim")
	}

	expiresAtUnix, ok := claims["exp"].(float64)
	if !ok {
		logger.Error("cannot get token's exp claim")
		return 0, time.Time{}, errors.New("cannot get token's exp claim")
	}

	return uint64(accountID), time.Unix(int64(expiresAtUnix), 0), nil

}

// CreateTokenString implements Token.
func (t *tokenLogic) CreateTokenString(ctx context.Context, accountID uint64) (string, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	expiresAt := time.Now().Add(t.authConfig.Token.GetTokenDuration())
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": accountID,
		"exp": expiresAt.Unix(),
		"kid": t.tokenPublicKeyID,
	})

	tokenString, err := token.SignedString(t.tokenPrivateKeyValue)
	if err != nil {
		logger.Error("failed signing token", zap.Error(err))
		return "", time.Time{}, nil
	}

	return tokenString, expiresAt, nil
}

// WithDatabase implements Token.
func (t *tokenLogic) WithDatabase(database database.Database) TokenLogic {
	panic("unimplemented")
}

func (t *tokenLogic) getJWTPublicKeyValue(ctx context.Context, tokenPublicKeyID uint64) (*rsa.PublicKey, error) {
	logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("tokenPublicKeyID", tokenPublicKeyID))

	var tokenPublicKeyValue database.TokenPublicKey

	cacheHit := true
	bytes, err := t.tokenPublicKeyCache.Get(ctx, fmt.Sprintf("%d", tokenPublicKeyID))
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to get tokenPublicKeyValue from cache, will fall back to database")
		cacheHit = false
	} else {
		tokenPublicKeyValue = database.TokenPublicKey{
			TokenPublicKeyID:    tokenPublicKeyID,
			TokenPublicKeyValue: bytes,
		}
	}

	if !cacheHit {
		tokenPublicKeyValue, err = t.tokenPublicKeyDataAccessor.GetPublicKey(ctx, tokenPublicKeyID)
		if err != nil {
			logger.Error("cannot get token's public key from database", zap.Error(err))
			return nil, err
		}
	}

	return jwt.ParseRSAPublicKeyFromPEM(tokenPublicKeyValue.TokenPublicKeyValue)
}

func pemEncodePublicKey(pubKey *rsa.PublicKey) ([]byte, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	return pem.EncodeToMemory(block), nil
}

func generateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	privateKeyPair, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return privateKeyPair, nil
}
