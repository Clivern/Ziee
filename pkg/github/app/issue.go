// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// GetIssue fetches an issue by number.
func (a *App) GetIssue(ctx context.Context, installationID int64, owner, repo string, number int) (*Issue, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var issue Issue
	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d", a.apiURL, owner, repo, number)
	err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &issue)
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

// ListIssues lists issues in a repository.
func (a *App) ListIssues(ctx context.Context, installationID int64, owner, repo string) ([]Issue, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var all []Issue
	for page := 1; ; page++ {
		var issues []Issue
		path := fmt.Sprintf("%s/repos/%s/%s/issues?per_page=%d&page=%d", a.apiURL, owner, repo, AppPerPage, page)
		err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &issues)
		if err != nil {
			return nil, err
		}

		all = append(all, issues...)
		if len(issues) < AppPerPage {
			return all, nil
		}
	}
}

// IsFirstIssue reports whether login has at most one issue in the repository.
func (a *App) IsFirstIssue(ctx context.Context, installationID int64, owner, repo, login string) (bool, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return false, err
	}

	var result struct {
		TotalCount int `json:"total_count"`
	}
	path := fmt.Sprintf("%s/search/issues?per_page=1&q=%s", a.apiURL, url.QueryEscape(
		fmt.Sprintf("repo:%s/%s author:%s type:issue", owner, repo, login),
	))
	err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &result)
	if err != nil {
		return false, err
	}

	return result.TotalCount <= 1, nil
}

// IsFirstPullRequest reports whether login has at most one pull request in the repository.
func (a *App) IsFirstPullRequest(ctx context.Context, installationID int64, owner, repo, login string) (bool, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return false, err
	}

	var result struct {
		TotalCount int `json:"total_count"`
	}
	path := fmt.Sprintf("%s/search/issues?per_page=1&q=%s", a.apiURL, url.QueryEscape(
		fmt.Sprintf("repo:%s/%s author:%s type:pr", owner, repo, login),
	))

	err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &result)
	if err != nil {
		return false, err
	}

	return result.TotalCount <= 1, nil
}

// CreateIssue opens an issue in a repository.
func (a *App) CreateIssue(ctx context.Context, installationID int64, owner, repo, title, body string) (*Issue, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var issue Issue
	path := fmt.Sprintf("%s/repos/%s/%s/issues", a.apiURL, owner, repo)
	err = Call(ctx, http.MethodPost, path, token.Token, GetHeaders(), map[string]string{
		"title": title,
		"body":  body,
	}, &issue)
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

// CreateComment adds a comment on an issue or pull request.
func (a *App) CreateComment(ctx context.Context, installationID int64, owner, repo string, number int, body string) (*Comment, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var comment Comment
	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", a.apiURL, owner, repo, number)
	err = Call(ctx, http.MethodPost, path, token.Token, GetHeaders(), map[string]string{"body": body}, &comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// CloseIssue closes an issue or pull request.
func (a *App) CloseIssue(ctx context.Context, installationID int64, owner, repo string, number int) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d", a.apiURL, owner, repo, number)

	return Call(ctx, http.MethodPatch, path, token.Token, GetHeaders(), map[string]string{"state": "closed"}, nil)
}

// ReopenIssue reopens an issue or pull request.
func (a *App) ReopenIssue(ctx context.Context, installationID int64, owner, repo string, number int) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d", a.apiURL, owner, repo, number)

	return Call(ctx, http.MethodPatch, path, token.Token, GetHeaders(), map[string]string{"state": "open"}, nil)
}

// AddAssignees assigns users to an issue or pull request.
func (a *App) AddAssignees(ctx context.Context, installationID int64, owner, repo string, number int, users []string) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d/assignees", a.apiURL, owner, repo, number)

	return Call(ctx, http.MethodPost, path, token.Token, GetHeaders(), map[string][]string{"assignees": users}, nil)
}

// RemoveAssignees removes users from an issue or pull request.
func (a *App) RemoveAssignees(ctx context.Context, installationID int64, owner, repo string, number int, users []string) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d/assignees", a.apiURL, owner, repo, number)

	return Call(ctx, http.MethodDelete, path, token.Token, GetHeaders(), map[string][]string{"assignees": users}, nil)
}
