// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	oauthAuthorizeURL = "https://github.com/login/oauth/authorize"
	oauthTokenURL     = "https://github.com/login/oauth/access_token"
)

// OAuthConfig holds GitHub OAuth app credentials.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AllowSignup  bool
}

// OAuth exchanges authorization codes for access tokens.
type OAuth struct {
	cfg        OAuthConfig
	httpClient *http.Client
}

// Token is an OAuth access token response.
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type oauthTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	State        string `json:"state,omitempty"`
}

// NewOAuth returns an OAuth helper for the given app configuration.
func NewOAuth(cfg OAuthConfig, httpClient *http.Client) *OAuth {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &OAuth{cfg: cfg, httpClient: httpClient}
}

// AuthorizeURL builds the GitHub authorization URL for the provided state.
func (o *OAuth) AuthorizeURL(state string) (string, error) {
	u, err := url.Parse(oauthAuthorizeURL)
	if err != nil {
		return "", fmt.Errorf("github oauth authorize url: %w", err)
	}

	q := u.Query()
	q.Set("client_id", o.cfg.ClientID)
	q.Set("redirect_uri", o.cfg.RedirectURL)
	q.Set("scope", strings.Join(o.cfg.Scopes, ","))
	q.Set("state", state)
	if !o.cfg.AllowSignup {
		q.Set("allow_signup", "false")
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// Exchange validates state and exchanges an authorization code for an access token.
func (o *OAuth) Exchange(ctx context.Context, code, state, expectedState string) (*Token, error) {
	if state != expectedState {
		return nil, fmt.Errorf("%w: mismatch", ErrInvalidOAuthState)
	}

	payload, err := json.Marshal(oauthTokenRequest{
		ClientID:     o.cfg.ClientID,
		ClientSecret: o.cfg.ClientSecret,
		Code:         code,
		State:        expectedState,
	})
	if err != nil {
		return nil, fmt.Errorf("github oauth encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthTokenURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("github oauth build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github oauth request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("github oauth read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newAPIError(resp.StatusCode, body)
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("github oauth decode body: %w", err)
	}

	return &token, nil
}
