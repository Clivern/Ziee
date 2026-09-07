// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"crypto/rsa"
	"fmt"
	"os"
)

const (
	AppAPIURL     = "https://api.github.com"
	AppUserAgent  = "ziee"
	AppAccept     = "application/vnd.github+json"
	AppAPIVersion = "2022-11-28"
	AppPerPage    = 100
)

var app *App

// App authenticates as a GitHub App using the private key.
type App struct {
	config Config
	key    *rsa.PrivateKey
	apiURL string
}

// Init loads the GitHub App private key from config.
func Init() error {
	client, err := New()
	if err != nil {
		return err
	}

	app = client

	return nil
}

// Get returns the GitHub App client loaded at startup.
func Get() *App {
	return app
}

// New returns a GitHub App client with the private key loaded from config.
func New() (*App, error) {
	return NewFromConfig(GetConfig())
}

// NewFromConfig returns a GitHub App client with the private key loaded from config.
func NewFromConfig(config Config) (*App, error) {
	pem, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("github app read key: %w", err)
	}

	key, err := ParsePrivateKey(pem)
	if err != nil {
		return nil, err
	}

	return &App{
		config: config,
		key:    key,
		apiURL: AppAPIURL,
	}, nil
}
