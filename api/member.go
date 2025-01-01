// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/middleware"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// ListWorkspaceMembersAction returns workspace members for managers.
func ListWorkspaceMembersAction(w http.ResponseWriter, r *http.Request) {
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
		Msg("Listing workspace members")

	limit, offset := util.ParsePagination(r)

	wm := module.NewWorkspace(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	result, err := wm.ListWorkspaceMembers(
		db.Id(wid),
		limit,
		offset,
	)
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
				Str("userId", user.Id.String()).
				Msg("Failed to list workspace members")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_workspace_members"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"members": result.Members,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// UpdateWorkspaceMemberRoleAction updates a workspace member role.
func UpdateWorkspaceMemberRoleAction(w http.ResponseWriter, r *http.Request) {
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

	memberUserId := chi.URLParam(r, "memberUserId")
	if lo.IsEmpty(memberUserId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_user_id"),
		})
		return
	}

	var req module.UpdateWorkspaceMemberRoleRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("memberUserId", memberUserId).
		Str("userId", user.Id.String()).
		Str("role", req.Role).
		Msg("Updating workspace member role")

	wm := module.NewWorkspace(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	member, err := wm.UpdateWorkspaceMemberRole(
		db.Id(wid),
		db.Id(memberUserId),
		req.Role,
	)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrWorkspaceUserNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_member_not_found"),
			})
		default:
			log.Error().
				Err(err).
				Str("workspaceId", wid).
				Str("memberUserId", memberUserId).
				Str("userId", user.Id.String()).
				Msg("Failed to update workspace member role")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_workspace_member"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, member)
}

// DeleteWorkspaceMemberAction removes a user from a workspace.
func DeleteWorkspaceMemberAction(w http.ResponseWriter, r *http.Request) {
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

	memberUserId := chi.URLParam(r, "memberUserId")
	if lo.IsEmpty(memberUserId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_user_id"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("memberUserId", memberUserId).
		Str("userId", user.Id.String()).
		Msg("Removing workspace member")

	wm := module.NewWorkspace(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	err := wm.DeleteWorkspaceMember(
		db.Id(wid),
		db.Id(memberUserId),
	)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrWorkspaceUserNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_member_not_found"),
			})
		default:
			log.Error().
				Err(err).
				Str("workspaceId", wid).
				Str("memberUserId", memberUserId).
				Str("userId", user.Id.String()).
				Msg("Failed to delete workspace member")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_workspace_member"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
