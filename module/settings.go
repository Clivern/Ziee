// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"fmt"
	"strings"

	"github.com/actx0/ziee/db"

	"github.com/samber/lo"
)

// Settings manages app settings.
type Settings struct {
	ConfigRepository db.ConfigRepository
}

// NewSettings creates a settings module with the given config repo.
func NewSettings(config db.ConfigRepository) *Settings {
	return &Settings{ConfigRepository: config}
}

// Update persists the provided settings to the config repository.
func (s *Settings) Update(platformEmail string, maintenanceMode bool) error {
	updates := map[string]string{
		"platform_email":   strings.TrimSpace(platformEmail),
		"maintenance_mode": lo.Ternary(maintenanceMode, "1", "0"),
	}

	for key, value := range updates {
		err := s.ConfigRepository.Update(key, value)
		if err != nil {
			return fmt.Errorf("failed to update config %s: %w", key, err)
		}
	}

	return nil
}

// GetSettings retrieves the current settings from the config repository.
func (s *Settings) GetSettings() (map[string]string, error) {
	keys := []string{
		"platform_email",
		"maintenance_mode",
	}

	settings := make(map[string]string, len(keys))

	for _, key := range keys {
		config, err := s.ConfigRepository.Get(key)
		if err != nil {
			return nil, err
		}
		if config != nil {
			settings[key] = config.Value
		} else {
			settings[key] = ""
		}
	}

	return settings, nil
}
