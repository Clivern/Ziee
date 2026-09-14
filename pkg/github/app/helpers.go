// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

// httpClient is the default HTTP client for the GitHub API.
var httpClient = req.C()

// StatusError is a non-2xx GitHub API response.
type StatusError struct {
	Method string
	URL    string
	Status int
	Body   string
}

// Error returns the error message.
func (e *StatusError) Error() string {
	return fmt.Sprintf("http %s %s: status %d: %s", e.Method, e.URL, e.Status, e.Body)
}

// GetHeaders returns the default GitHub API request headers.
func GetHeaders() map[string]string {
	return map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}
}

// Call sends an authenticated JSON request to the GitHub API.
func Call(ctx context.Context, method, url, token string, headers map[string]string, payload, dest any) error {
	r := httpClient.R().
		SetContext(ctx).
		SetHeaders(headers)

	if lo.IsNotEmpty(token) {
		r.SetBearerAuthToken(token)
	}

	if !lo.IsNil(payload) {
		r.SetBody(payload)
	}

	if !lo.IsNil(dest) {
		r.SetSuccessResult(dest)
	}

	resp, err := r.Send(method, url)
	if err != nil {
		return fmt.Errorf("http %s %s: %w", method, url, err)
	}

	if !resp.IsSuccessState() {
		return &StatusError{
			Method: method,
			URL:    url,
			Status: resp.StatusCode,
			Body:   resp.String(),
		}
	}

	return nil
}
