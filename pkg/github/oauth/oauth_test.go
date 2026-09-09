// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitOAuth(t *testing.T) {
	t.Run("PrimaryEmail", func(t *testing.T) {
		assert.Equal(t, "fallback@example.com", PrimaryEmail(nil, "fallback@example.com"))

		assert.Equal(t, "b@example.com", PrimaryEmail([]Email{
			{Email: "a@example.com", Verified: true},
			{Email: "b@example.com", Primary: true, Verified: true},
		}, "fallback"))

		assert.Equal(t, "first@example.com", PrimaryEmail([]Email{
			{Email: "unverified@example.com", Primary: true},
			{Email: "first@example.com", Verified: true},
			{Email: "second@example.com", Verified: true},
		}, "fallback"))
	})

	t.Run("AuthorizeURL", func(t *testing.T) {
		oauth := NewOAuth(OAuthConfig{
			ClientID:    "cid",
			RedirectURL: "https://app.example/callback",
			Scopes:      []string{"read:user", "user:email"},
			AllowSignup: true,
		})

		raw := oauth.AuthorizeURL("state-1")
		u, err := url.Parse(raw)
		assert.NoError(t, err)
		assert.Equal(t, OauthAuthorizeURL, u.Scheme+"://"+u.Host+u.Path)
		q := u.Query()
		assert.Equal(t, "cid", q.Get("client_id"))
		assert.Equal(t, "https://app.example/callback", q.Get("redirect_uri"))
		assert.Equal(t, "read:user,user:email", q.Get("scope"))
		assert.Equal(t, "state-1", q.Get("state"))
		assert.Empty(t, q.Get("allow_signup"))

		blocked := NewOAuth(OAuthConfig{AllowSignup: false})
		assert.Contains(t, blocked.AuthorizeURL("s"), "allow_signup=false")
	})

	t.Run("Exchange invalid state", func(t *testing.T) {
		oauth := NewOAuth(OAuthConfig{})
		token, err := oauth.Exchange(context.Background(), "code", "a", "b")
		assert.Nil(t, token)
		assert.ErrorIs(t, err, ErrInvalidOAuthState)
	})

	t.Run("Exchange User Emails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/login/oauth/access_token":
				assert.Equal(t, http.MethodPost, r.Method)
				_ = json.NewEncoder(w).Encode(Token{
					AccessToken: "gho_token",
					TokenType:   "bearer",
					Scope:       "read:user",
				})
			case "/user":
				assert.Equal(t, "Bearer gho_token", r.Header.Get("Authorization"))
				_ = json.NewEncoder(w).Encode(UserInfo{ID: 7, Login: "octocat", Name: "Octo", Email: "octo@example.com"})
			case "/user/emails":
				assert.Equal(t, "Bearer gho_token", r.Header.Get("Authorization"))
				_ = json.NewEncoder(w).Encode([]Email{
					{Email: "octo@example.com", Primary: true, Verified: true},
				})
			default:
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		prevToken, prevUser, prevEmails := OauthTokenURL, OauthUserURL, OauthEmailsURL
		OauthTokenURL = server.URL + "/login/oauth/access_token"
		OauthUserURL = server.URL + "/user"
		OauthEmailsURL = server.URL + "/user/emails"
		t.Cleanup(func() {
			OauthTokenURL, OauthUserURL, OauthEmailsURL = prevToken, prevUser, prevEmails
		})

		oauth := NewOAuth(OAuthConfig{
			ClientID:     "cid",
			ClientSecret: "secret",
		})
		ctx := context.Background()

		token, err := oauth.Exchange(ctx, "code", "state", "state")
		assert.NoError(t, err)
		assert.Equal(t, "gho_token", token.AccessToken)

		user, err := oauth.User(ctx, token.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, "octocat", user.Login)

		emails, err := oauth.Emails(ctx, token.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, "octo@example.com", PrimaryEmail(emails, ""))
	})
}
