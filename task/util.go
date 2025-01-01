// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package task

import (
	"encoding/json"
	"errors"
)

var ErrInvalidPayload = errors.New("invalid async task payload")

// EncodeResult JSON-encodes a task result map.
func EncodeResult(result any) (string, error) {
	raw, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

// DecodePayload JSON-decodes a task payload into a map.
func DecodePayload(payload *string) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(*payload), &result)
	if err != nil {
		return nil, ErrInvalidPayload
	}

	return result, nil
}
