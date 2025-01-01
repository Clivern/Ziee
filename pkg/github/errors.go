// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidWebhookSignature = errors.New("invalid github webhook signature")
	ErrInvalidOAuthState       = errors.New("invalid oauth state")
)

// APIError is returned when the GitHub REST API responds with a non-success status.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("github api: status %d: %s", e.StatusCode, e.Body)
}

func newAPIError(status int, body []byte) error {
	return &APIError{
		StatusCode: status,
		Body:       string(body),
	}
}
