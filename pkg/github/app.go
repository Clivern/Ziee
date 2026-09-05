// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	AppAPIURL       = "https://api.github.com"
	AppUserAgent    = "ziee"
	AppAccept       = "application/vnd.github+json"
	AppAPIVersion   = "2022-11-28"
	AppReposPerPage = 100
)

var app *App

// App authenticates as a GitHub App using the private key.
type App struct {
	config Config
	key    *rsa.PrivateKey
	apiURL string
}

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

// Repository is a repo accessible to an installation.
type Repository struct {
	ID       int64  `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type repositoriesResponse struct {
	TotalCount   int          `json:"total_count"`
	Repositories []Repository `json:"repositories"`
}

// Init loads the GitHub App private key from config.
func Init() error {
	client, err := New()
	if err != nil {
		return err
	}

	app = client

	return nil
}

// Get returns the GitHub App client loaded at startup.
func Get() *App {
	return app
}

// New returns a GitHub App client with the private key loaded from config.
func New() (*App, error) {
	return NewFromConfig(GetConfig())
}

// NewFromConfig returns a GitHub App client with the private key loaded from config.
func NewFromConfig(config Config) (*App, error) {
	pem, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("github app read key: %w", err)
	}

	key, err := ParsePrivateKey(pem)
	if err != nil {
		return nil, err
	}

	return &App{
		config: config,
		key:    key,
		apiURL: AppAPIURL,
	}, nil
}

// Installation fetches a GitHub App installation by id.
func (a *App) Installation(ctx context.Context, installationID int64) (*Installation, error) {
	token, err := a.JWT()
	if err != nil {
		return nil, err
	}

	var installation Installation
	err = a.do(ctx, http.MethodGet, "/app/installations/"+strconv.FormatInt(installationID, 10), token, &installation)
	if err != nil {
		return nil, err
	}

	return &installation, nil
}

// InstallationToken creates an installation access token.
func (a *App) InstallationToken(ctx context.Context, installationID int64) (*InstallationToken, error) {
	jwt, err := a.JWT()
	if err != nil {
		return nil, err
	}

	var token InstallationToken
	err = a.do(ctx, http.MethodPost, "/app/installations/"+strconv.FormatInt(installationID, 10)+"/access_tokens", jwt, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Repositories lists repos the installation can access.
func (a *App) Repositories(ctx context.Context, installationID int64) ([]Repository, error) {
	token, err := a.InstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var body repositoriesResponse
	err = a.do(ctx, http.MethodGet, "/installation/repositories?per_page="+strconv.Itoa(AppReposPerPage), token.Token, &body)
	if err != nil {
		return nil, err
	}

	return body.Repositories, nil
}

func (a *App) do(ctx context.Context, method, path, token string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, method, a.apiURL+path, nil)
	if err != nil {
		return fmt.Errorf("github app request: %w", err)
	}

	req.Header.Set("Accept", AppAccept)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", AppUserAgent)
	req.Header.Set("X-GitHub-Api-Version", AppAPIVersion)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("github app %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("github app read body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("github app %s %s: status %d: %s", method, path, resp.StatusCode, body)
	}

	if dest == nil {
		return nil
	}

	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("github app decode body: %w", err)
	}

	return nil
}
