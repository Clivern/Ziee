// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package cache

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitGetConfig(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("app.cache.provider", "redis")
	viper.Set("app.cache.redis.addr", "localhost:6379")
	viper.Set("app.cache.redis.password", "secret")
	viper.Set("app.cache.redis.db", 2)
	viper.Set("app.cache.redis.pool_size", 10)
	viper.Set("app.cache.redis.min_idle_conns", 1)

	cfg := GetConfig()
	assert.Equal(t, "redis", cfg.Provider)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "secret", cfg.Redis.Password)
	assert.Equal(t, 2, cfg.Redis.DB)
	assert.Equal(t, 10, cfg.Redis.PoolSize)
	assert.Equal(t, 1, cfg.Redis.MinIdleConns)
}
