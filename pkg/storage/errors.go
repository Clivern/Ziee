// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package storage

import (
	"errors"
)

var (
	ErrUnsupportedProvider = errors.New("unsupported storage provider")
	ErrDocumentNotFound    = errors.New("document not found")
)
