// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package cache

import (
	"errors"
)

var (
	ErrUnsupportedProvider = errors.New("unsupported cache provider")
	ErrNotFound            = errors.New("cache key not found")
)
