package cache

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	tokenPublicKeyIDPrefix string = "token_public_key:"
)

type TokenPublicKey interface {
	Set(ctx context.Context, tokenPublicKeyID string, tokenPublicKeyValue []byte) error
	Get(ctx context.Context, tokenPublicKeyID string) ([]byte, error)
}

func NewTokenPublicKey(client Client) (TokenPublicKey, error) {
	return &tokenPublicKey{
		client: client,
	}, nil
}

type tokenPublicKey struct {
	client Client
}

// Get implements TokenPublicKey.
func (t *tokenPublicKey) Get(ctx context.Context, tokenPublicKeyID string) ([]byte, error) {
	key := t.getCacheKey(tokenPublicKeyID)
	value, err := t.client.Get(ctx, key)
	if err != nil {
		return []byte{}, err
	}

	stringValue, ok := value.(string)
	if !ok {
		return []byte{}, errors.New("cached value is not a string")
	}

	tokenPublicKeyValue, err := t.decodeBase64(stringValue)
	if err != nil {
		return []byte{}, err
	}

	return tokenPublicKeyValue, nil
}

// Set implements TokenPublicKey.
func (t *tokenPublicKey) Set(ctx context.Context, tokenPublicKeyID string, tokenPublicKeyValue []byte) error {
	key := t.getCacheKey(tokenPublicKeyID)
	value := t.encodeBase64(tokenPublicKeyValue)
	return t.client.Set(ctx, key, value, 0)
}

func (t *tokenPublicKey) encodeBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func (t *tokenPublicKey) decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func (t *tokenPublicKey) getCacheKey(tokenPublicKeyID string) string {
	return fmt.Sprintf("%s:%s", tokenPublicKeyIDPrefix, tokenPublicKeyID)
}
