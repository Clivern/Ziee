// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cache

import (
	"context"
	"fmt"
	"time"
)

// Cache is a low-latency key-value cache.
type Cache interface {
	Config() Config
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Ping(ctx context.Context) error
	Close() error
}

// New returns a cache from app.cache config.
func New() (Cache, error) {
	config := GetConfig()

	switch config.Provider {
	case "redis":
		return NewRedisClient(config)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, config.Provider)
	}
}
