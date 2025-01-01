// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
)

// SetupAction runs the initial platform setup.
func SetupAction(w http.ResponseWriter, r *http.Request) {
	var req module.SetupRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	log.Info().
		Str("adminEmail", req.AdminEmail).
		Msg("New setup request")

	sm := module.NewSetup(
		db.NewConfigRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	err = sm.Install(&req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrPlatformAlreadyInstalled):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "platform_already_installed"),
			})
		case errors.Is(err, module.ErrFailedCompleteSetup):
			log.Error().
				Err(err).
				Str("adminEmail", req.AdminEmail).
				Msg("Failed to complete setup")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_complete_setup"),
			})
		default:
			log.Error().
				Err(err).
				Str("adminEmail", req.AdminEmail).
				Msg("Setup failed")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_complete_setup"),
			})
		}
		return
	}

	log.Info().
		Str("platformEmail", req.PlatformEmail).
		Str("adminEmail", req.AdminEmail).
		Msg("Platform setup completed")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"successMessage": locale.TR(r, "setup_completed_successfully"),
	})
}

// SetupStatusAction returns whether the platform is already installed.
func SetupStatusAction(w http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("Setup status request")

	sm := module.NewSetup(
		db.NewConfigRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"installed": sm.IsInstalled(),
	})
}
