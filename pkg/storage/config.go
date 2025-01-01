// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package storage

import (
	"github.com/spf13/viper"
)

// Config holds document storage configuration.
type Config struct {
	Provider string
	Local    LocalConfig
	S3       S3Config
}

// LocalConfig holds local filesystem storage settings.
type LocalConfig struct {
	Path string
}

// S3Config holds S3-compatible object storage settings.
type S3Config struct {
	Bucket          string
	Region          string
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
	UsePathStyle    bool
	Prefix          string
}

// GetConfig loads document storage settings from app.storage config.
func GetConfig() Config {
	return Config{
		Provider: viper.GetString("app.storage.provider"),
		Local: LocalConfig{
			Path: viper.GetString("app.storage.local.path"),
		},
		S3: S3Config{
			Bucket:          viper.GetString("app.storage.s3.bucket"),
			Region:          viper.GetString("app.storage.s3.region"),
			Endpoint:        viper.GetString("app.storage.s3.endpoint"),
			AccessKeyId:     viper.GetString("app.storage.s3.access_key_id"),
			SecretAccessKey: viper.GetString("app.storage.s3.secret_access_key"),
			UsePathStyle:    viper.GetBool("app.storage.s3.use_path_style"),
			Prefix:          viper.GetString("app.storage.s3.prefix"),
		},
	}
}
