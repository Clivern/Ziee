// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/pkg/util"

	"github.com/rs/zerolog/log"
)

// ReadyAction returns whether the app is ready (e.g. DB reachable).
func ReadyAction(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("Readiness check")

	err := db.GetDB().Ping()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Database ping failed during readiness check")

		util.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "not_ok",
		})
		return
	}

	log.Debug().Msg("Readiness check passed")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}
