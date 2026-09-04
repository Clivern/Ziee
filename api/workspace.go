// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/locale"
	"github.com/clivern/ziee/middleware"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// CreateWorkspaceAction creates a new workspace.
func CreateWorkspaceAction(w http.ResponseWriter, r *http.Request) {
	var req module.CreateWorkspaceRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("New workspace request")

	wm := module.NewWorkspace(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)
	workspace, err := wm.CreateWorkspace(&req, user)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to create workspace")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_create_workspace"),
		})
		return
	}

	log.Info().
		Str("workspaceId", workspace.Id.String()).
		Str("userId", user.Id.String()).
		Msg("Workspace created")

	util.WriteJSON(w, http.StatusCreated, workspace)
}

// ListWorkspacesAction returns workspaces the user is a member of (paginated).
func ListWorkspacesAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Listing workspaces")

	limit, offset := util.ParsePagination(r)

	wm := module.NewWorkspace(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	result, err := wm.ListWorkspaces(user, limit, offset)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to list workspaces")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_list_workspaces"),
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"workspaces": result.Workspaces,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetWorkspaceAction returns one workspace by Id.
func GetWorkspaceAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("userId", user.Id.String()).
		Msg("Getting workspace")

	wm := module.NewWorkspace(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	workspace, err := wm.GetWorkspace(db.Id(wid), user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
			return
		default:
			log.Error().
				Err(err).
				Str("workspaceId", wid).
				Str("userId", user.Id.String()).
				Msg("Failed to get workspace")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_workspace"),
			})
			return
		}
	}

	util.WriteJSON(w, http.StatusOK, workspace)
}

// UpdateWorkspaceAction updates a workspace.
func UpdateWorkspaceAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("userId", user.Id.String()).
		Msg("Updating workspace")

	var req module.UpdateWorkspaceRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	wm := module.NewWorkspace(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	workspace, err := wm.UpdateWorkspace(db.Id(wid), &req, user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
			return
		default:
			log.Error().
				Err(err).
				Str("workspaceId", wid).
				Str("userId", user.Id.String()).
				Msg("Failed to update workspace")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_workspace"),
			})
			return
		}
	}

	log.Info().
		Str("workspaceId", workspace.Id.String()).
		Str("userId", user.Id.String()).
		Msg("Workspace updated")

	util.WriteJSON(w, http.StatusOK, workspace)
}

// DeleteWorkspaceAction deletes a workspace.
func DeleteWorkspaceAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Msg("Deleting workspace")

	wm := module.NewWorkspace(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	err := wm.DeleteWorkspace(db.Id(wid))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
			return
		default:
			log.Error().
				Err(err).
				Str("workspaceId", wid).
				Msg("Failed to delete workspace")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_workspace"),
			})
			return
		}
	}

	log.Info().
		Str("workspaceId", wid).
		Msg("Workspace deleted")

	util.WriteJSON(w, http.StatusNoContent, map[string]any{})
}
