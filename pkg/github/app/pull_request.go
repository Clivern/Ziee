// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
)

// PullRequest is a GitHub pull request.
type PullRequest struct {
	ID      int64  `json:"id"`
	Number  int    `json:"number"`
	Title   string `json:"title"`
	State   string `json:"state"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	Draft   bool   `json:"draft"`
	Merged  bool   `json:"merged"`
	User    User   `json:"user"`
	Head    struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

// NewPullRequest is a branch, file, and pull request to open against the default branch.
type NewPullRequest struct {
	Branch  string
	Path    string
	Content string
	Title   string
	Body    string
}

// GetPullRequest fetches a pull request by number.
func (a *App) GetPullRequest(ctx context.Context, installationID int64, owner, repo string, number int) (*PullRequest, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var pull PullRequest
	path := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", a.apiURL, owner, repo, number)
	err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &pull)
	if err != nil {
		return nil, err
	}

	return &pull, nil
}

// ListPullRequests lists pull requests in a repository.
func (a *App) ListPullRequests(ctx context.Context, installationID int64, owner, repo string) ([]PullRequest, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var all []PullRequest
	for page := 1; ; page++ {
		var pulls []PullRequest
		path := fmt.Sprintf("%s/repos/%s/%s/pulls?per_page=%d&page=%d", a.apiURL, owner, repo, AppPerPage, page)
		err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &pulls)
		if err != nil {
			return nil, err
		}

		all = append(all, pulls...)
		if len(pulls) < AppPerPage {
			return all, nil
		}
	}
}

// CreatePullRequest creates a branch, adds a file, and opens a pull request.
func (a *App) CreatePullRequest(ctx context.Context, installationID int64, owner, repo string, in NewPullRequest) (*PullRequest, error) {
	repository, err := a.GetRepository(ctx, installationID, owner, repo)
	if err != nil {
		return nil, err
	}

	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var ref struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	path := fmt.Sprintf("%s/repos/%s/%s/git/ref/heads/%s", a.apiURL, owner, repo, repository.DefaultBranch)
	err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &ref)
	if err != nil {
		return nil, err
	}

	path = fmt.Sprintf("%s/repos/%s/%s/git/refs", a.apiURL, owner, repo)
	err = Call(ctx, http.MethodPost, path, token.Token, GetHeaders(), map[string]string{
		"ref": "refs/heads/" + in.Branch,
		"sha": ref.Object.SHA,
	}, nil)
	if err != nil {
		return nil, err
	}

	err = a.CreateFile(ctx, installationID, owner, repo, in.Branch, in.Path, in.Content, in.Title)
	if err != nil {
		return nil, err
	}

	var pull PullRequest
	path = fmt.Sprintf("%s/repos/%s/%s/pulls", a.apiURL, owner, repo)
	err = Call(ctx, http.MethodPost, path, token.Token, GetHeaders(), map[string]string{
		"title": in.Title,
		"body":  in.Body,
		"head":  in.Branch,
		"base":  repository.DefaultBranch,
	}, &pull)
	if err != nil {
		return nil, err
	}

	return &pull, nil
}
