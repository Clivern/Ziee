// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package util

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/drone/envsubst"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// GetBaseDir returns the project base dir
func GetBaseDir(dirName string) string {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to read working directory")
	}
	cacheDir := fmt.Sprintf("%s/%s", baseDir, dirName)

	for {
		if fi, err := os.Stat(cacheDir); err == nil {
			if fi.Mode().IsDir() {
				return baseDir
			}
		}
		baseDir = filepath.Dir(baseDir)
		cacheDir = fmt.Sprintf("%s/%s", baseDir, dirName)
	}
}

// LoadConfigs load configs for testing purposes using viper
func LoadConfigs(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	data1, err := envsubst.EvalEnv(string(data))
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer([]byte(data1)))

	return nil
}
