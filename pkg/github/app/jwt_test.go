// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitGetConfig(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("app.oauth.github.client_id", "app-123")
	viper.Set("app.oauth.github.private_key_path", "/tmp/key.pem")

	cfg := GetConfig()
	assert.Equal(t, "app-123", cfg.ClientID)
	assert.Equal(t, "/tmp/key.pem", cfg.PrivateKeyPath)
}

func TestUnitParsePrivateKeyAndJWT(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	pkcs1 := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	parsed, err := ParsePrivateKey(pkcs1)
	assert.NoError(t, err)
	assert.Equal(t, key.N, parsed.N)

	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(key)
	assert.NoError(t, err)
	pkcs8 := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8Bytes})
	parsed8, err := ParsePrivateKey(pkcs8)
	assert.NoError(t, err)
	assert.Equal(t, key.N, parsed8.N)

	client := &App{config: Config{ClientID: "Iv1.abc"}, key: key}
	token, err := client.JWT()
	assert.NoError(t, err)
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3)

	_, err = NewFromConfig(Config{PrivateKeyPath: "/missing/key.pem"})
	assert.Error(t, err)
}
