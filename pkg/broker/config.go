// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package broker

import (
	"github.com/spf13/viper"
)

// Config holds NATS connection settings.
type Config struct {
	URL   string
	Name  string
	Queue string
}

// GetConfig loads broker settings from app.broker config.
func GetConfig() Config {
	return Config{
		URL:   viper.GetString("app.broker.url"),
		Name:  viper.GetString("app.broker.name"),
		Queue: viper.GetString("app.broker.queue"),
	}
}
