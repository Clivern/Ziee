// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/github"
	"github.com/actx0/ziee/pkg/util"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

const oauthStateCookie = "_ziee_oauth_state"

// GitHubOAuthStartAction redirects the browser to GitHub's authorize URL.
func GitHubOAuthStartAction(w http.ResponseWriter, r *http.Request) {
	oauth := newGitHubOAuth()
	state, err := util.GenerateSecureToken(24)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate oauth state")
		http.Redirect(w, r, oauthLoginErrorURL(), http.StatusFound)
		return
	}

	authorizeURL, err := oauth.AuthorizeURL(state)
	if err != nil {
		log.Error().Err(err).Msg("Failed to build github authorize url")
		http.Redirect(w, r, oauthLoginErrorURL(), http.StatusFound)
		return
	}

	opts := oauthStateCookieOptions()
	opts.MaxAge = int((10 * time.Minute) / time.Second)
	util.SetCookie(w, oauthStateCookie, state, opts)

	http.Redirect(w, r, authorizeURL, http.StatusFound)
}

// GitHubOAuthCallbackAction completes GitHub OAuth and creates a session.
func GitHubOAuthCallbackAction(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	expectedState := util.GetCookie(r, oauthStateCookie)

	util.DeleteCookie(w, oauthStateCookie)

	oauth := newGitHubOAuth()
	token, err := oauth.Exchange(r.Context(), code, state, expectedState)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth exchange failed")
		http.Redirect(w, r, oauthLoginErrorURL(), http.StatusFound)
		return
	}

	ghUser, err := oauth.User(r.Context(), token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth user fetch failed")
		http.Redirect(w, r, oauthLoginErrorURL(), http.StatusFound)
		return
	}

	emails, err := oauth.Emails(r.Context(), token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth emails fetch failed")
		http.Redirect(w, r, oauthLoginErrorURL(), http.StatusFound)
		return
	}

	email := primaryGitHubEmail(emails, ghUser.Email)
	name := module.OAuthDisplayName(ghUser.Name, ghUser.Login)

	auth := module.NewAuth(
		db.NewUserRepository(db.GetDB()),
		db.NewSessionRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
	)

	result, err := auth.LoginWithOAuth(r.Context(), &module.OAuthIdentity{
		Provider:       db.UserProviderGithub,
		ProviderUserID: module.OAuthProviderUserID(ghUser.ID),
		Email:          email,
		Name:           name,
	})
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth login failed")
		http.Redirect(w, r, oauthLoginErrorURL(), http.StatusFound)
		return
	}

	util.SetCookie(w, "_ziee_session", result.Session.Token, result.CookieOptions)
	http.Redirect(w, r, oauthLoginSuccessURL(), http.StatusFound)
}

func newGitHubOAuth() *github.OAuth {
	redirectURL := viper.GetString("app.oauth.github.redirect_url")
	if lo.IsEmpty(redirectURL) {
		redirectURL = strings.TrimRight(viper.GetString("app.url"), "/") + "/api/v1/public/action/oauth/github/callback"
	}

	return github.NewOAuth(github.OAuthConfig{
		ClientID:     viper.GetString("app.oauth.github.client_id"),
		ClientSecret: viper.GetString("app.oauth.github.client_secret"),
		RedirectURL:  redirectURL,
		Scopes:       []string{"read:user", "user:email"},
		AllowSignup:  true,
	}, nil)
}

func primaryGitHubEmail(emails []github.Email, fallback string) string {
	var firstVerified string
	for _, e := range emails {
		if !e.Verified {
			continue
		}
		if e.Primary {
			return e.Email
		}
		if lo.IsEmpty(firstVerified) {
			firstVerified = e.Email
		}
	}
	if lo.IsNotEmpty(firstVerified) {
		return firstVerified
	}

	return fallback
}

func oauthStateCookieOptions() *util.CookieOptions {
	opts := util.DefaultCookieOptions()
	if strings.HasPrefix(viper.GetString("app.url"), "https://") {
		opts.Secure = true
	}
	opts.SameSite = http.SameSiteLaxMode

	return opts
}

func oauthLoginSuccessURL() string {
	return strings.TrimRight(viper.GetString("app.url"), "/") + "/login?oauth=github"
}

func oauthLoginErrorURL() string {
	return strings.TrimRight(viper.GetString("app.url"), "/") + "/login?" + url.Values{
		"oauth_error": {"github"},
	}.Encode()
}
