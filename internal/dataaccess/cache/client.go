package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
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
	switch cacheConfig.Type {
	case configs.CacheTypeInMemory:
		return NewInMemoryClient(cacheConfig, logger)
	case configs.CacheTypeRedis:
		return NewRedisClient(cacheConfig, logger)
	default:
		err := fmt.Errorf(`invalid cache type, expect one of ["redis", "in-memory"], got %s`, string(cacheConfig.Type))
		logger.With(zap.Error(err)).Error("invalid cache type")
		return nil, err
	}
}

func NewRedisClient(
	cacheConfig configs.Cache,
	logger *zap.Logger,
) (Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     cacheConfig.Addr,
		Username: cacheConfig.Username,
		Password: cacheConfig.Password,
		DB:       cacheConfig.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("can not connect to redis client", zap.Error(err))
		return nil, err
	}

	return &redisClient{
		client: client,
		logger: logger,
	}, nil
}

type redisClient struct {
	client *redis.Client
	logger *zap.Logger
}

// AddToSet implements Client.
func (c *redisClient) AddToSet(ctx context.Context, key string, value ...interface{}) error {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("key", key)).
		With(zap.Any("value", value))

	err := c.client.SAdd(ctx, key, value).Err()
	if err != nil {
		logger.Error("failed to add value to set", zap.Error(err))
		return err
	}
	return nil
}

// IsValueInSet implements Client.
func (c *redisClient) IsValueInSet(ctx context.Context, key string, value interface{}) (bool, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("key", key)).
		With(zap.Any("value", value))

	exists, err := c.client.SIsMember(ctx, key, value).Result()
	if err != nil {
		logger.Error("failed to find value in set", zap.Error(err))
		return false, err
	}
	return exists, nil
}

// Get implements Client.
func (c *redisClient) Get(ctx context.Context, key string) (any, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("key", key))

	value, err := c.client.Get(ctx, key).Result()
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
func (c *redisClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("key", key)).
		With(zap.Any("value", value)).
		With(zap.Duration("ttl", ttl))

	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		logger.Error("failed to set value into cache", zap.Error(err))
		return err
	}

	return nil
}

func NewInMemoryClient(
	cacheConfig configs.Cache,
	logger *zap.Logger,
) (Client, error) {
	return &inMemoryClient{
		cache:      make(map[string]any),
		cacheMutex: &sync.Mutex{},
		logger:     logger,
	}, nil
}

type inMemoryClient struct {
	cache      map[string]any
	cacheMutex *sync.Mutex
	logger     *zap.Logger
}

// AddToSet implements Client.
func (i *inMemoryClient) AddToSet(ctx context.Context, key string, value ...any) error {
	i.cacheMutex.Lock()
	defer i.cacheMutex.Unlock()

	if _, ok := i.cache[key]; !ok {
		i.cache[key] = make(map[any]struct{})
	}

	set := i.cache[key].(map[any]struct{})
	for _, v := range value {
		set[v] = struct{}{}
	}

	return nil
}

// Get implements Client.
func (i *inMemoryClient) Get(ctx context.Context, key string) (any, error) {
	if val, ok := i.cache[key]; ok {
		return val, nil
	}

	return nil, ErrCacheMissed
}

// IsValueInSet implements Client.
func (i *inMemoryClient) IsValueInSet(ctx context.Context, key string, value any) (bool, error) {
	if set, ok := i.cache[key].(map[any]struct{}); ok {
		if _, exists := set[value]; exists {
			return true, nil
		}
	}

	return false, nil
}

// Set implements Client.
func (i *inMemoryClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	i.cacheMutex.Lock()
	defer i.cacheMutex.Unlock()

	i.cache[key] = value

	return nil
}
