// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/resend"
	"github.com/actx0/ziee/pkg/util"

	"github.com/rs/zerolog/log"
)

// ForgotPasswordAction sends a password reset email.
func ForgotPasswordAction(w http.ResponseWriter, r *http.Request) {
	var req module.ForgotPasswordRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	log.Info().
		Str("email", req.Email).
		Msg("New forgot password request")

	au := module.NewAuth(
		db.NewUserRepository(db.GetDB()),
		db.NewSessionRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewPasswordResetTokenRepository(db.GetDB()),
		resend.NewMailer(),
	)

	err = au.ForgotPassword(r.Context(), &req)
	if err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("Failed to process forgot password request")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_process_request"),
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"successMessage": locale.TR(r, "forgot_password_success"),
	})
}

// ResetPasswordAction resets the password using the token from the email.
func ResetPasswordAction(w http.ResponseWriter, r *http.Request) {
	var req module.ResetPasswordRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	log.Info().Msg("New password reset request")

	au := module.NewAuth(
		db.NewUserRepository(db.GetDB()),
		db.NewSessionRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewPasswordResetTokenRepository(db.GetDB()),
		resend.NewMailer(),
	)

	err = au.ResetPassword(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrResetTokenInvalid):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "reset_link_invalid_or_expired"),
			})
		default:
			log.Error().Err(err).Msg("Failed to reset password")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_reset_password"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"successMessage": locale.TR(r, "password_reset_success"),
	})
}
