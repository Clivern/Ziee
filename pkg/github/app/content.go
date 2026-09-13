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

// FileExists returns the first path that exists on the repository's default branch.
func (a *App) FileExists(ctx context.Context, installationID int64, owner, repo string, paths []string) (string, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return "", err
	}

	for _, path := range paths {
		url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", a.apiURL, owner, repo, path)
		err = Call(ctx, http.MethodGet, url, token.Token, GetHeaders(), nil, nil)
		if err == nil {
			return path, nil
		}

		var status *StatusError
		if errors.As(err, &status) && status.Status == http.StatusNotFound {
			continue
		}

		return "", err
	}

	return "", nil
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
	err = Call(ctx, http.MethodGet, url, token.Token, GetHeaders(), nil, &file)
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

	return Call(ctx, http.MethodPut, url, token.Token, GetHeaders(), map[string]string{
		"message": message,
		"content": base64.StdEncoding.EncodeToString([]byte(content)),
		"branch":  branch,
	}, nil)
}
