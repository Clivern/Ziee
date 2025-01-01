// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package storage

import (
	"context"
	"fmt"
)

// Store persists document text at "{uuid}.txt".
type Store interface {
	Config() Config
	PutDocument(ctx context.Context, id string, text string) (string, error)
	GetDocument(ctx context.Context, id string) (string, error)
	DeleteDocument(ctx context.Context, id string) error
	DocumentURI(id string) (string, error)
}

// New returns a document store from app.storage config.
func New() (Store, error) {
	return NewFromConfig(GetConfig())
}

// NewFromConfig returns a document store for the given configuration.
func NewFromConfig(config Config) (Store, error) {
	switch config.Provider {
	case "local":
		return NewLocalStore(config)
	case "s3":
		return NewS3Store(config)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, config.Provider)
	}
}
