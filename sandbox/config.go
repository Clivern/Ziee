// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"github.com/spf13/viper"
)

// Config holds sandbox configuration.
type Config struct {
	APIKey      string
	DockerImage string
	TempDir     string
	Model       string
	MinPort     int
}

// GetConfig loads sandbox settings from app.sandbox config.
func GetConfig() Config {
	return Config{
		APIKey:      viper.GetString("app.sandbox.api_key"),
		DockerImage: viper.GetString("app.sandbox.docker_image"),
		TempDir:     viper.GetString("app.sandbox.temp_dir"),
		Model:       viper.GetString("app.sandbox.model"),
		MinPort:     viper.GetInt("app.sandbox.min_port"),
	}
}
