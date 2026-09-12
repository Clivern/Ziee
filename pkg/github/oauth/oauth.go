// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package oauth

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

const OauthAuthorizeURL = "https://github.com/login/oauth/authorize"

var (
	OauthTokenURL  = "https://github.com/login/oauth/access_token"
	OauthUserURL   = "https://api.github.com/user"
	OauthEmailsURL = "https://api.github.com/user/emails"
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
	cfg OAuthConfig
}

// Token is an OAuth access token response.
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type OauthTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	State        string `json:"state,omitempty"`
}

// UserInfo is the authenticated GitHub user.
type UserInfo struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Email is a GitHub account email address.
type Email struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// PrimaryEmail returns the primary verified email, else the first verified, else fallback.
func PrimaryEmail(emails []Email, fallback string) string {
	var firstVerified string
	for _, e := range emails {
		if !e.Verified {
			continue
		}
		if e.Primary {
			return e.Email
		}
		if firstVerified == "" {
			firstVerified = e.Email
		}
	}
	if firstVerified != "" {
		return firstVerified
	}

	return fallback
}

// NewOAuth returns an OAuth helper for the given app configuration.
func NewOAuth(cfg OAuthConfig) *OAuth {
	return &OAuth{cfg: cfg}
}

// AuthorizeURL builds the GitHub authorization URL for the provided state.
func (o *OAuth) AuthorizeURL(state string) string {
	q := url.Values{
		"client_id":    {o.cfg.ClientID},
		"redirect_uri": {o.cfg.RedirectURL},
		"scope":        {strings.Join(o.cfg.Scopes, ",")},
		"state":        {state},
	}
	if !o.cfg.AllowSignup {
		q.Set("allow_signup", "false")
	}

	return OauthAuthorizeURL + "?" + q.Encode()
}

// Exchange validates state and exchanges an authorization code for an access token.
func (o *OAuth) Exchange(ctx context.Context, code, state, expectedState string) (*Token, error) {
	if state != expectedState {
		return nil, ErrInvalidOAuthState
	}

	var token Token
	err := Call(ctx, http.MethodPost, OauthTokenURL, "", map[string]string{
		"Accept": "application/json",
	}, OauthTokenRequest{
		ClientID:     o.cfg.ClientID,
		ClientSecret: o.cfg.ClientSecret,
		Code:         code,
		State:        expectedState,
	}, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// User fetches the authenticated GitHub user.
func (o *OAuth) User(ctx context.Context, accessToken string) (*UserInfo, error) {
	var user UserInfo
	err := Call(ctx, http.MethodGet, OauthUserURL, accessToken, map[string]string{
		"Accept": "application/json",
	}, nil, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Emails fetches the authenticated user's email addresses.
func (o *OAuth) Emails(ctx context.Context, accessToken string) ([]Email, error) {
	var emails []Email
	err := Call(ctx, http.MethodGet, OauthEmailsURL, accessToken, map[string]string{
		"Accept": "application/json",
	}, nil, &emails)
	if err != nil {
		return nil, err
	}

	return emails, nil
}
