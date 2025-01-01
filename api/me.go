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
	"github.com/samber/lo"
)

// GetMeAction returns the authenticated API key or access key principal.
func GetMeAction(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	accessKey := r.Header.Get("X-Access-Key")

	if lo.IsEmpty(apiKey) && lo.IsEmpty(accessKey) {
		util.WriteJSON(w, http.StatusForbidden, map[string]any{
			"errorMessage": locale.TR(r, "me_requires_key_header"),
		})
		return
	}

	mm := module.NewMe(
		db.NewAPIKeyRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewWorkspaceAccessKeyRepository(db.GetDB()),
	)

	// Check if the request is for an API key
	if lo.IsNotEmpty(apiKey) {
		me, err := mm.GetByAPIKey(apiKey)

		if err != nil {
			switch {
			case errors.Is(err, module.ErrAPIKeyNotFound):
				util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
					"errorMessage": locale.TR(r, "invalid_api_key"),
				})
			default:
				log.Error().Err(err).Msg("Failed to resolve API key for /me")
				util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
					"errorMessage": locale.TR(r, "failed_get_me"),
				})
			}
			return
		}

		util.WriteJSON(w, http.StatusOK, me)
		return
	}

	// Check if the request is for an access key
	me, err := mm.GetByAccessKey(accessKey)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrAccessKeyNotFound):
			util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"errorMessage": locale.TR(r, "invalid_access_key"),
			})
		default:
			log.Error().Err(err).Msg("Failed to resolve access key for /me")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_me"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, me)
}
