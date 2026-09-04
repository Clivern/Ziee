// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is a Redis-backed cache client.
type RedisClient struct {
	config Config
	client *redis.Client
}

// NewRedisClient creates a Redis-backed cache client.
func NewRedisClient(config Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		MinIdleConns: config.Redis.MinIdleConns,
	})

	return &RedisClient{
		config: config,
		client: client,
	}, nil
}

// Config returns the cache configuration.
func (c *RedisClient) Config() Config {
	return c.config
}

// Get returns the cached value for key.
func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("redis get: %w", err)
	}

	return value, nil
}

// Set stores value at key with an optional TTL.
func (c *RedisClient) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

// Delete removes key from the cache.
func (c *RedisClient) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis delete: %w", err)
	}

	return nil
}

// Exists reports whether key is present in the cache.
func (c *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists: %w", err)
	}

	return count > 0, nil
}

// Ping checks Redis connectivity.
func (c *RedisClient) Ping(ctx context.Context) error {
	err := c.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	return nil
}

// Close closes the Redis client and its connection pool.
func (c *RedisClient) Close() error {
	err := c.client.Close()
	if err != nil {
		return fmt.Errorf("redis close: %w", err)
	}

	return nil
}
