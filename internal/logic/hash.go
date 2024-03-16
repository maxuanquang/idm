package logic

import (
	"context"

	"github.com/maxuanquang/idm/internal/configs"
	"golang.org/x/crypto/bcrypt"
)

type Hash interface {
	HashPassword(ctx context.Context, plainPassword string) (string, error)
	IsHashEqual(ctx context.Context, plainPassword string, hashedPassword string) (bool, error)
}

type hash struct {
	authConfig configs.Auth
}

func NewHash() Hash {
	return &hash{}
}

// HashPassword implements Hash.
func (h *hash) HashPassword(ctx context.Context, plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), h.authConfig.Hash.Cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// IsHashEqual implements Hash.
func (h *hash) IsHashEqual(ctx context.Context, hashedPassword string, plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return false, err
	}
	return true, nil
}
