// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package cache

import (
	"github.com/spf13/viper"
)

// Config holds cache configuration.
type Config struct {
	Provider string
	Redis    RedisConfig
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// GetConfig loads cache settings from app.cache config.
func GetConfig() Config {
	return Config{
		Provider: viper.GetString("app.cache.provider"),
		Redis: RedisConfig{
			Addr:         viper.GetString("app.cache.redis.addr"),
			Password:     viper.GetString("app.cache.redis.password"),
			DB:           viper.GetInt("app.cache.redis.db"),
			PoolSize:     viper.GetInt("app.cache.redis.pool_size"),
			MinIdleConns: viper.GetInt("app.cache.redis.min_idle_conns"),
		},
	}
}
