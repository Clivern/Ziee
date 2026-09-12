// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// DefaultZieeYML is the starter policy file added on repository bootstrap.
const DefaultZieeYML = "version: 1.0.0\n"

// FileExists reports whether a path exists in the repository's default branch.
func (a *App) FileExists(ctx context.Context, installationID int64, owner, repo, path string) (bool, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", a.apiURL, owner, repo, path)
	err = call(ctx, http.MethodGet, url, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, nil, nil)
	if err == nil {
		return true, nil
	}

	var status *StatusError
	if errors.As(err, &status) && status.Status == http.StatusNotFound {
		return false, nil
	}

	return false, err
}

// GetFile returns the decoded contents of a path on the default branch.
func (a *App) GetFile(ctx context.Context, installationID int64, owner, repo, path string) ([]byte, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var file struct {
		Content string `json:"content"`
	}
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", a.apiURL, owner, repo, path)
	err = call(ctx, http.MethodGet, url, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, nil, &file)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(strings.ReplaceAll(file.Content, "\n", ""))
}

// CreateFile creates a file on a branch.
func (a *App) CreateFile(ctx context.Context, installationID int64, owner, repo, branch, path, content, message string) error {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", a.apiURL, owner, repo, path)

	return call(ctx, http.MethodPut, url, token.Token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, map[string]string{
		"message": message,
		"content": base64.StdEncoding.EncodeToString([]byte(content)),
		"branch":  branch,
	}, nil)
}
