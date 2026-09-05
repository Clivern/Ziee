// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"github.com/spf13/viper"
)

// Config holds GitHub App settings.
type Config struct {
	ClientID       string
	PrivateKeyPath string
}

// GetConfig loads GitHub App settings from app.oauth.github config.
func GetConfig() Config {
	return Config{
		ClientID:       viper.GetString("app.oauth.github.client_id"),
		PrivateKeyPath: viper.GetString("app.oauth.github.private_key_path"),
	}
}
