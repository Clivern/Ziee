// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package oauth

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

var httpClient = req.C()

// Call sends an authenticated JSON request to GitHub OAuth APIs.
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
		return fmt.Errorf("http %s %s: status %d: %s", method, url, resp.StatusCode, resp.Bytes())
	}

	return nil
}
