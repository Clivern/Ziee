// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const checksAccept = "application/vnd.github+json"

type callOptions struct {
	method       string
	path         string
	query        url.Values
	body         any
	accept       string
	expectedCode int
}

func (c *Client) repoPath(path string, args ...any) string {
	repo := fmt.Sprintf("/repos/%s/%s", c.owner, c.repository)
	if len(args) > 0 {
		path = fmt.Sprintf(path, args...)
	}
	return c.baseURL + repo + path
}

func (c *Client) call(ctx context.Context, opts callOptions, out any) error {
	var body io.Reader
	if opts.body != nil {
		encoded, err := json.Marshal(opts.body)
		if err != nil {
			return fmt.Errorf("github api encode body: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	target := opts.path
	if len(opts.query) > 0 {
		u, err := url.Parse(target)
		if err != nil {
			return fmt.Errorf("github api parse url: %w", err)
		}
		u.RawQuery = opts.query.Encode()
		target = u.String()
	}

	req, err := http.NewRequestWithContext(ctx, opts.method, target, body)
	if err != nil {
		return fmt.Errorf("github api build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	accept := opts.accept
	if accept == "" {
		accept = "application/vnd.github+json"
	}
	req.Header.Set("Accept", accept)
	if opts.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("github api request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("github api read body: %w", err)
	}

	if resp.StatusCode != opts.expectedCode {
		return newAPIError(resp.StatusCode, respBody)
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("github api decode body: %w", err)
	}

	return nil
}

func (c *Client) callNoContent(ctx context.Context, method, path string, expectedCode int) error {
	return c.call(ctx, callOptions{
		method:       method,
		path:         path,
		expectedCode: expectedCode,
	}, nil)
}
