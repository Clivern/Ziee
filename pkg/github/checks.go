// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"context"
	"net/http"
	"net/url"

	"github.com/actx0/ziee/pkg/github/response"
	"github.com/actx0/ziee/pkg/github/sender"
)

// CreateCheckRun creates a check run.
func (c *Client) CreateCheckRun(ctx context.Context, checkRun sender.CheckRun) (response.CheckRun, error) {
	var out response.CheckRun
	err := c.call(ctx, callOptions{
		method:       http.MethodPost,
		path:         c.repoPath("/check-runs"),
		body:         checkRun,
		accept:       checksAccept,
		expectedCode: http.StatusCreated,
	}, &out)
	return out, err
}

// UpdateCheckRun updates a check run.
func (c *Client) UpdateCheckRun(ctx context.Context, id int, checkRun sender.CheckRun) (response.CheckRun, error) {
	var out response.CheckRun
	err := c.call(ctx, callOptions{
		method:       http.MethodPatch,
		path:         c.repoPath("/check-runs/%d", id),
		body:         checkRun,
		accept:       checksAccept,
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// ListRefCheckRuns lists check runs for a git ref.
func (c *Client) ListRefCheckRuns(ctx context.Context, ref, checkName, status, filter string) (response.CheckRuns, error) {
	var out response.CheckRuns
	query := url.Values{}
	if checkName != "" {
		query.Set("check_name", checkName)
	}
	if status != "" {
		query.Set("status", status)
	}
	if filter != "" {
		query.Set("filter", filter)
	}

	err := c.call(ctx, callOptions{
		method:       http.MethodGet,
		path:         c.repoPath("/commits/%s/check-runs", ref),
		query:        query,
		accept:       checksAccept,
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// ListSuiteCheckRuns lists check runs in a check suite.
func (c *Client) ListSuiteCheckRuns(ctx context.Context, checkSuiteID int, checkName, status, filter string) (response.CheckRuns, error) {
	var out response.CheckRuns
	query := url.Values{}
	if checkName != "" {
		query.Set("check_name", checkName)
	}
	if status != "" {
		query.Set("status", status)
	}
	if filter != "" {
		query.Set("filter", filter)
	}

	err := c.call(ctx, callOptions{
		method:       http.MethodGet,
		path:         c.repoPath("/check-suites/%d/check-runs", checkSuiteID),
		query:        query,
		accept:       checksAccept,
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// GetCheckRun returns a single check run.
func (c *Client) GetCheckRun(ctx context.Context, id int) (response.CheckRun, error) {
	var out response.CheckRun
	err := c.call(ctx, callOptions{
		method:       http.MethodGet,
		path:         c.repoPath("/check-runs/%d", id),
		accept:       checksAccept,
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// ListCheckRunAnnotations lists annotations for a check run.
func (c *Client) ListCheckRunAnnotations(ctx context.Context, id int) (response.Annotations, error) {
	var out response.Annotations
	err := c.call(ctx, callOptions{
		method:       http.MethodGet,
		path:         c.repoPath("/check-runs/%d/annotations", id),
		accept:       checksAccept,
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}
