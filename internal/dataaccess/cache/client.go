package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

var (
	ErrCacheMissed = errors.New("cache miss")
)

type Client interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	AddToSet(ctx context.Context, key string, value ...any) error
	IsValueInSet(ctx context.Context, key string, value any) (bool, error)
}

func NewClient(
	cacheConfig configs.Cache,
	logger *zap.Logger,
) (Client, error) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cacheConfig.Addr,
		Username: cacheConfig.Username,
		Password: cacheConfig.Password,
		DB:       cacheConfig.DB,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("can not connect to redis client", zap.Error(err))
		return nil, err
	}

	return &client{
		redisClient: redisClient,
		logger:      logger,
	}, nil
}

type client struct {
	redisClient *redis.Client
	logger      *zap.Logger
}

// AddToSet implements Client.
func (c *client) AddToSet(ctx context.Context, key string, value ...interface{}) error {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("key", key)).
		With(zap.Any("value", value))

	err := c.redisClient.SAdd(ctx, key, value).Err()
	if err != nil {
		logger.Error("failed to add value to set", zap.Error(err))
		return err
	}
	return nil
}

// IsValueInSet implements Client.
func (c *client) IsValueInSet(ctx context.Context, key string, value interface{}) (bool, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("key", key)).
		With(zap.Any("value", value))

	exists, err := c.redisClient.SIsMember(ctx, key, value).Result()
	if err != nil {
		logger.Error("failed to find value in set", zap.Error(err))
		return false, err
	}
	return exists, nil
}

// Get implements Client.
func (c *client) Get(ctx context.Context, key string) (any, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("key", key))

	value, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMissed
		}
		logger.Error("failed to get key from cache", zap.Error(err))
		return nil, err
	}

	return value, nil
}

// Set implements Client.
func (c *client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("key", key)).
		With(zap.Any("value", value)).
		With(zap.Duration("ttl", ttl))

	if err := c.redisClient.Set(ctx, key, value, ttl).Err(); err != nil {
		logger.Error("failed to set value into cache", zap.Error(err))
		return err
	}

	return nil
}
