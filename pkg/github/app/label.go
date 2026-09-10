// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Label is a GitHub repository label.
type Label struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// ListLabels lists labels in a repository.
func (a *App) ListLabels(ctx context.Context, installationID int64, owner, repo string) ([]Label, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var all []Label
	for page := 1; ; page++ {
		var labels []Label
		path := fmt.Sprintf("%s/repos/%s/%s/labels?per_page=%d&page=%d", a.apiURL, owner, repo, AppPerPage, page)
		err = call(ctx, http.MethodGet, path, token.Token, map[string]string{
			"Accept":               AppAccept,
			"User-Agent":           AppUserAgent,
			"X-GitHub-Api-Version": AppAPIVersion,
		}, nil, &labels)
		if err != nil {
			return nil, err
		}

		all = append(all, labels...)
		if len(labels) < AppPerPage {
			return all, nil
		}
	}
}

// AddLabels adds labels to an issue or pull request.
func (a *App) AddLabels(ctx context.Context, installationID int64, owner, repo string, number int, labels []string) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/repos/%s/%s/issues/%d/labels", a.apiURL, owner, repo, number)

	return call(ctx, http.MethodPost, path, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, labels, nil)
}

// RemoveLabels removes labels from an issue or pull request.
func (a *App) RemoveLabels(ctx context.Context, installationID int64, owner, repo string, number int, labels []string) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	for _, label := range labels {
		path := fmt.Sprintf("%s/repos/%s/%s/issues/%d/labels/%s", a.apiURL, owner, repo, number, url.PathEscape(label))
		err = call(ctx, http.MethodDelete, path, token.Token, map[string]string{
			"Accept":               AppAccept,
			"User-Agent":           AppUserAgent,
			"X-GitHub-Api-Version": AppAPIVersion,
		}, nil, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
