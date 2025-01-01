// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// GetWorkspaceStatsAction returns dashboard metrics for a workspace.
func GetWorkspaceStatsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	sm := module.NewStats(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceStatsRepository(db.GetDB()),
	)

	stats, err := sm.GetWorkspaceStats(db.Id(wid))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().
				Err(err).
				Str("workspaceId", wid).
				Msg("Failed to get workspace stats")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_workspace_stats"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, stats)
}
