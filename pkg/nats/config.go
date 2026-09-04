// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package nats

import (
	"github.com/spf13/viper"
)

// Config holds NATS connection settings.
type Config struct {
	URL   string
	Name  string
	Queue string
}

// GetConfig loads NATS settings from app.nats config.
func GetConfig() Config {
	return Config{
		URL:   viper.GetString("app.nats.url"),
		Name:  viper.GetString("app.nats.name"),
		Queue: viper.GetString("app.nats.queue"),
	}
}
