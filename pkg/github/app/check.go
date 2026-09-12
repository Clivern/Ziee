// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
)

// CheckRun is a GitHub check run.
type CheckRun struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	HeadSHA    string `json:"head_sha"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	HTMLURL    string `json:"html_url"`
}

// CreateCheckRun starts an in-progress check run on a commit.
func (a *App) CreateCheckRun(ctx context.Context, installationID int64, owner, repo, name, sha string) (*CheckRun, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var run CheckRun
	url := fmt.Sprintf("%s/repos/%s/%s/check-runs", a.apiURL, owner, repo)
	err = call(ctx, http.MethodPost, url, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, map[string]string{
		"name":     name,
		"head_sha": sha,
		"status":   "in_progress",
	}, &run)
	if err != nil {
		return nil, err
	}

	return &run, nil
}

// CompleteCheckRun marks a check run completed with a conclusion and output.
func (a *App) CompleteCheckRun(ctx context.Context, installationID int64, owner, repo string, id int64, conclusion, title, summary string) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/repos/%s/%s/check-runs/%d", a.apiURL, owner, repo, id)

	return call(ctx, http.MethodPatch, url, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, map[string]any{
		"status":     "completed",
		"conclusion": conclusion,
		"output": map[string]string{
			"title":   title,
			"summary": summary,
		},
	}, nil)
}
