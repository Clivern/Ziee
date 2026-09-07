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

// GetPullRequest fetches a pull request by number.
func (a *App) GetPullRequest(ctx context.Context, installationID int64, owner, repo string, number int) (*PullRequest, error) {
	token, err := a.CreateInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var pull PullRequest
	path := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", a.apiURL, owner, repo, number)
	err = call(ctx, http.MethodGet, path, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, nil, &pull)
	if err != nil {
		return nil, err
	}

	return &pull, nil
}

// ListPullRequests lists pull requests in a repository.
func (a *App) ListPullRequests(ctx context.Context, installationID int64, owner, repo string) ([]PullRequest, error) {
	token, err := a.CreateInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var all []PullRequest
	for page := 1; ; page++ {
		var pulls []PullRequest
		path := fmt.Sprintf("%s/repos/%s/%s/pulls?per_page=%d&page=%d", a.apiURL, owner, repo, AppPerPage, page)
		err = call(ctx, http.MethodGet, path, token.Token, map[string]string{
			"Accept":               AppAccept,
			"User-Agent":           AppUserAgent,
			"X-GitHub-Api-Version": AppAPIVersion,
		}, nil, &pulls)
		if err != nil {
			return nil, err
		}

		all = append(all, pulls...)
		if len(pulls) < AppPerPage {
			return all, nil
		}
	}
}
