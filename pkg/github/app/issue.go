// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
)

// User is a GitHub user on an issue, pull request, or assignee list.
type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

// Issue is a GitHub issue.
type Issue struct {
	ID        int64   `json:"id"`
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	State     string  `json:"state"`
	Body      string  `json:"body"`
	HTMLURL   string  `json:"html_url"`
	User      User    `json:"user"`
	Labels    []Label `json:"labels"`
	Assignees []User  `json:"assignees"`
}

// Comment is a GitHub issue or pull request comment.
type Comment struct {
	ID      int64  `json:"id"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    User   `json:"user"`
}

// GetIssue fetches an issue by number.
func (a *App) GetIssue(ctx context.Context, installationID int64, owner, repo string, number int) (*Issue, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var issue Issue
	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d", a.apiURL, owner, repo, number)
	err = call(ctx, http.MethodGet, path, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, nil, &issue)
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
		err = call(ctx, http.MethodGet, path, token.Token, map[string]string{
			"Accept":               AppAccept,
			"User-Agent":           AppUserAgent,
			"X-GitHub-Api-Version": AppAPIVersion,
		}, nil, &issues)
		if err != nil {
			return nil, err
		}

		all = append(all, issues...)
		if len(issues) < AppPerPage {
			return all, nil
		}
	}
}

// CreateIssue opens an issue in a repository.
func (a *App) CreateIssue(ctx context.Context, installationID int64, owner, repo, title, body string) (*Issue, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var issue Issue
	path := fmt.Sprintf("%s/repos/%s/%s/issues", a.apiURL, owner, repo)
	err = call(ctx, http.MethodPost, path, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, map[string]string{
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
	err = call(ctx, http.MethodPost, path, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, map[string]string{"body": body}, &comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// CloseIssue closes an issue or pull request.
func (a *App) CloseIssue(ctx context.Context, installationID int64, owner, repo string, number int) error {
	return a.setIssueState(ctx, installationID, owner, repo, number, "closed")
}

// ReopenIssue reopens an issue or pull request.
func (a *App) ReopenIssue(ctx context.Context, installationID int64, owner, repo string, number int) error {
	return a.setIssueState(ctx, installationID, owner, repo, number, "open")
}

// AddAssignees assigns users to an issue or pull request.
func (a *App) AddAssignees(ctx context.Context, installationID int64, owner, repo string, number int, users []string) error {
	return a.mutateAssignees(ctx, http.MethodPost, installationID, owner, repo, number, users)
}

// RemoveAssignees removes users from an issue or pull request.
func (a *App) RemoveAssignees(ctx context.Context, installationID int64, owner, repo string, number int, users []string) error {
	return a.mutateAssignees(ctx, http.MethodDelete, installationID, owner, repo, number, users)
}

func (a *App) setIssueState(ctx context.Context, installationID int64, owner, repo string, number int, state string) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d", a.apiURL, owner, repo, number)

	return call(ctx, http.MethodPatch, path, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, map[string]string{"state": state}, nil)
}

func (a *App) mutateAssignees(ctx context.Context, method string, installationID int64, owner, repo string, number int, users []string) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d/assignees", a.apiURL, owner, repo, number)

	return call(ctx, method, path, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, map[string][]string{"assignees": users}, nil)
}
