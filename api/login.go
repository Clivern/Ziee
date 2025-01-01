// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/resend"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
)

// LoginAction logs the user in and creates a session.
func LoginAction(w http.ResponseWriter, r *http.Request) {
	var req module.LoginRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	log.Info().
		Str("email", req.Email).
		Msg("New login request")

	lm := module.NewAuth(
		db.NewUserRepository(db.GetDB()),
		db.NewSessionRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewPasswordResetTokenRepository(db.GetDB()),
		resend.NewMailer(),
	)

	result, err := lm.Login(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrInvalidCredentials):
			util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"errorMessage": locale.TR(r, "invalid_username_or_password"),
			})
		case errors.Is(err, module.ErrAccountNotActive):
			util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"errorMessage": locale.TR(r, "account_not_active"),
			})
		case errors.Is(err, module.ErrMaintenanceModeEnabled):
			util.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
				"errorMessage": locale.TR(r, "maintenance_mode_enabled"),
			})
		case errors.Is(err, module.ErrFailedCheckMaintenance):
			log.Error().Err(err).Msg("Failed to get maintenance_mode config")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_check_maintenance"),
			})
		case errors.Is(err, module.ErrFailedCreateSession):
			log.Error().Err(err).Msg("Failed to create session")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_session"),
			})
		default:
			log.Error().Err(err).Msg("Login failed")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_session"),
			})
		}
		return
	}

	util.SetCookie(w, "_actx0_session", result.Session.Token, result.CookieOptions)
	util.WriteJSON(w, http.StatusOK, map[string]any{
		"successMessage": locale.TR(r, "login_successful"),
		"user": map[string]any{
			"id":              result.User.Id,
			"email":           result.User.Email,
			"role":            result.User.Role,
			"isActive":        result.User.IsActive,
			"isEmailVerified": result.User.IsEmailVerified,
			"name":            result.User.Name,
			"avatar":          util.GravatarURL(result.User.Email, 80),
			"language":        result.User.Language,
			"theme":           result.User.Theme,
			"lastLoginAt":     result.User.LastLoginAt.UTC().Format(time.RFC3339),
			"createdAt":       result.User.CreatedAt.UTC().Format(time.RFC3339),
			"updatedAt":       result.User.UpdatedAt.UTC().Format(time.RFC3339),
		},
	})
}
