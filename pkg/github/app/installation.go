// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
)

// Installation is a GitHub App installation.
type Installation struct {
	ID                  int64               `json:"id"`
	AppID               int64               `json:"app_id"`
	AppSlug             string              `json:"app_slug"`
	TargetID            int64               `json:"target_id"`
	TargetType          string              `json:"target_type"`
	RepositorySelection string              `json:"repository_selection"`
	HTMLURL             string              `json:"html_url"`
	Account             InstallationAccount `json:"account"`
	Permissions         map[string]string   `json:"permissions"`
}

// InstallationAccount is the user or org the app is installed on.
type InstallationAccount struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Type      string `json:"type"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

// InstallationToken is a short-lived token for installation APIs.
type InstallationToken struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// GetInstallation fetches a GitHub App installation by id.
func (a *App) GetInstallation(ctx context.Context, installationID int64) (*Installation, error) {
	token, err := a.JWT()
	if err != nil {
		return nil, err
	}

	var installation Installation
	path := fmt.Sprintf("%s/app/installations/%d", a.apiURL, installationID)
	err = call(ctx, http.MethodGet, path, token, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, nil, &installation)
	if err != nil {
		return nil, err
	}

	return &installation, nil
}

// CreateInstallationToken creates an installation access token.
func (a *App) CreateInstallationToken(ctx context.Context, installationID int64) (*InstallationToken, error) {
	jwt, err := a.JWT()
	if err != nil {
		return nil, err
	}

	var token InstallationToken
	path := fmt.Sprintf("%s/app/installations/%d/access_tokens", a.apiURL, installationID)
	err = call(ctx, http.MethodPost, path, jwt, map[string]string{
		"Accept":               AppAccept,
		"User-Agent":           AppUserAgent,
		"X-GitHub-Api-Version": AppAPIVersion,
	}, nil, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
