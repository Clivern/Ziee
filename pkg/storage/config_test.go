// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package storage

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitGetConfig(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("app.storage.provider", "s3")
	viper.Set("app.storage.local.path", "/data")
	viper.Set("app.storage.s3.bucket", "docs")
	viper.Set("app.storage.s3.region", "us-east-1")
	viper.Set("app.storage.s3.endpoint", "http://localhost:9000")
	viper.Set("app.storage.s3.access_key_id", "ak")
	viper.Set("app.storage.s3.secret_access_key", "sk")
	viper.Set("app.storage.s3.use_path_style", true)
	viper.Set("app.storage.s3.prefix", "ziee/")

	cfg := GetConfig()
	assert.Equal(t, "s3", cfg.Provider)
	assert.Equal(t, "/data", cfg.Local.Path)
	assert.Equal(t, "docs", cfg.S3.Bucket)
	assert.Equal(t, "us-east-1", cfg.S3.Region)
	assert.Equal(t, "http://localhost:9000", cfg.S3.Endpoint)
	assert.Equal(t, "ak", cfg.S3.AccessKeyId)
	assert.Equal(t, "sk", cfg.S3.SecretAccessKey)
	assert.True(t, cfg.S3.UsePathStyle)
	assert.Equal(t, "ziee/", cfg.S3.Prefix)
}

func TestUnitDocumentKey(t *testing.T) {
	assert.Equal(t, "abc.txt", DocumentKey("abc"))
}

func TestUnitLocalStore(t *testing.T) {
	dir := t.TempDir()
	store, err := NewLocalStore(Config{Local: LocalConfig{Path: dir}})
	assert.NoError(t, err)
	assert.Equal(t, dir, store.Config().Local.Path)

	uri, err := store.PutDocument(context.Background(), "doc-1", "hello")
	assert.NoError(t, err)
	assert.NotEmpty(t, uri)

	text, err := store.GetDocument(context.Background(), "doc-1")
	assert.NoError(t, err)
	assert.Equal(t, "hello", text)

	docURI, err := store.DocumentURI("doc-1")
	assert.NoError(t, err)
	assert.NotEmpty(t, docURI)

	assert.NoError(t, store.DeleteDocument(context.Background(), "doc-1"))
}
