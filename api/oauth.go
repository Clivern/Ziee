// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"
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

const OauthStateCookie = "_ziee_oauth_state"

// GitHubOAuthStartAction redirects the browser to GitHub's authorize URL.
func GitHubOAuthStartAction(w http.ResponseWriter, r *http.Request) {
	errorURL := util.AppURL("/login?oauth_error=github")

	oauth := github.NewOAuth(github.OAuthConfig{
		ClientID:     viper.GetString("app.oauth.github.client_id"),
		ClientSecret: viper.GetString("app.oauth.github.client_secret"),
		RedirectURL:  viper.GetString("app.oauth.github.redirect_url"),
		Scopes:       []string{"read:user", "user:email"},
		AllowSignup:  true,
	})

	state, err := util.GenerateSecureToken(24)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate oauth state")
		http.Redirect(w, r, errorURL, http.StatusFound)
		return
	}

	authorizeURL := oauth.AuthorizeURL(state)

	opts := lo.Ternary(
		strings.HasPrefix(util.AppURL(""), "https://"),
		util.SecureCookieOptions(),
		util.DefaultCookieOptions(),
	)
	opts.SameSite = http.SameSiteLaxMode
	opts.MaxAge = int((10 * time.Minute) / time.Second)
	util.SetCookie(w, OauthStateCookie, state, opts)

	http.Redirect(w, r, authorizeURL, http.StatusFound)
}

// GitHubOAuthCallbackAction completes GitHub OAuth and creates a session.
func GitHubOAuthCallbackAction(w http.ResponseWriter, r *http.Request) {
	errorURL := util.AppURL("/login?oauth_error=github")

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	expectedState := util.GetCookie(r, OauthStateCookie)

	util.DeleteCookie(w, OauthStateCookie)

	oauth := github.NewOAuth(github.OAuthConfig{
		ClientID:     viper.GetString("app.oauth.github.client_id"),
		ClientSecret: viper.GetString("app.oauth.github.client_secret"),
		RedirectURL:  viper.GetString("app.oauth.github.redirect_url"),
		Scopes:       []string{"read:user", "user:email"},
		AllowSignup:  true,
	})

	token, err := oauth.Exchange(r.Context(), code, state, expectedState)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth exchange failed")
		http.Redirect(w, r, errorURL, http.StatusFound)
		return
	}

	ghUser, err := oauth.User(r.Context(), token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth user fetch failed")
		http.Redirect(w, r, errorURL, http.StatusFound)
		return
	}

	emails, err := oauth.Emails(r.Context(), token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth emails fetch failed")
		http.Redirect(w, r, errorURL, http.StatusFound)
		return
	}

	email := GetPrimaryGitHubEmail(emails, ghUser.Email)
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
		http.Redirect(w, r, errorURL, http.StatusFound)
		return
	}

	util.SetCookie(w, "_ziee_session", result.Session.Token, result.CookieOptions)
	http.Redirect(w, r, util.AppURL("/login?oauth=github"), http.StatusFound)
}

func GetPrimaryGitHubEmail(emails []github.Email, fallback string) string {
	verified := lo.Filter(emails, func(e github.Email, _ int) bool {
		return e.Verified
	})

	if e, ok := lo.Find(verified, func(e github.Email) bool {
		return e.Primary
	}); ok {
		return e.Email
	}

	return lo.CoalesceOrEmpty(lo.FirstOrEmpty(verified).Email, fallback)
}
