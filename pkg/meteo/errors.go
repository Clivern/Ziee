// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package meteo

import "errors"

var (
	ErrCityNotFound    = errors.New("city not found")
	ErrEmptyCityName   = errors.New("city name is required")
	ErrInvalidResponse = errors.New("invalid meteo api response")
)
