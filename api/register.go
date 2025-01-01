// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/util"

	"github.com/rs/zerolog/log"
)

// RegisterAction creates an account directly.
func RegisterAction(w http.ResponseWriter, r *http.Request) {
	var req module.RegisterRequest

	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	log.Info().
		Str("email", req.Email).
		Msg("New registration request")

	rm := module.NewRegister(db.NewUserRepository(db.GetDB()))

	result, err := rm.Register(&req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrUserEmailAlreadyExists):
			util.WriteJSON(w, http.StatusConflict, map[string]any{
				"errorMessage": locale.TR(r, "account_with_email_already_exists"),
			})
		case errors.Is(err, module.ErrFailedCompleteRegistration):
			log.Error().
				Err(err).
				Str("email", req.Email).
				Msg("Failed to complete registration")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_complete_registration"),
			})
		default:
			log.Error().
				Err(err).
				Str("email", req.Email).
				Msg("Registration failed")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_complete_registration"),
			})
		}
		return
	}

	log.Info().
		Str("userId", result.User.Id.String()).
		Str("email", result.User.Email).
		Msg("Registration completed")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"successMessage": locale.TR(r, "registration_successful"),
	})
}
