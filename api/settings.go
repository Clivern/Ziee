// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/locale"
	"github.com/clivern/ziee/middleware"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/util"

	"github.com/rs/zerolog/log"
)

// UpdateSettingsRequest is the body for updating app settings.
type UpdateSettingsRequest struct {
	PlatformEmail   string `json:"platformEmail" validate:"required,email,max=255" label:"Platform Email"`
	MaintenanceMode bool   `json:"maintenanceMode" label:"Maintenance Mode"`
}

// UpdateSettingsAction updates app settings (admin only).
func UpdateSettingsAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Updating settings")

	var req UpdateSettingsRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mod := module.NewSettings(db.NewConfigRepository(db.GetDB()))
	err = mod.Update(req.PlatformEmail, req.MaintenanceMode)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to update settings")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_update_settings"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Settings updated")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"successMessage": locale.TR(r, "settings_updated_successfully"),
	})
}

// GetSettingsAction returns current app settings.
func GetSettingsAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Getting settings")

	mod := module.NewSettings(db.NewConfigRepository(db.GetDB()))
	settings, err := mod.GetSettings()

	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to get settings")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_get_settings"),
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"settings": settings,
	})
}
