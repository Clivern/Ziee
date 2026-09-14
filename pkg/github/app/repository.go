// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
)

// Repository is a repo accessible to an installation.
type Repository struct {
	ID            int64  `json:"id"`
	NodeID        string `json:"node_id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type RepositoriesResponse struct {
	TotalCount   int          `json:"total_count"`
	Repositories []Repository `json:"repositories"`
}

// ListRepositories lists all repos the installation can access.
func (a *App) ListRepositories(ctx context.Context, installationID int64) ([]Repository, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var all []Repository
	for page := 1; ; page++ {
		var body RepositoriesResponse
		path := fmt.Sprintf("%s/installation/repositories?per_page=%d&page=%d", a.apiURL, AppPerPage, page)
		err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &body)
		if err != nil {
			return nil, err
		}

		all = append(all, body.Repositories...)
		if len(body.Repositories) < AppPerPage {
			return all, nil
		}
	}
}

// GetRepository fetches a repository by owner and name.
func (a *App) GetRepository(ctx context.Context, installationID int64, owner, repo string) (*Repository, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var repository Repository
	path := fmt.Sprintf("%s/repos/%s/%s", a.apiURL, owner, repo)
	err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &repository)
	if err != nil {
		return nil, err
	}

	return &repository, nil
}
