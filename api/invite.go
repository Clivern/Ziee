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
	"github.com/clivern/ziee/pkg/resend"
	"github.com/clivern/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// CreateInviteAction creates a new user invite.
func CreateInviteAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	workspaceId := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(workspaceId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Str("workspaceId", workspaceId).
		Msg("New invite request")

	var req module.CreateInviteRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	invite, err := im.CreateInvite(db.Id(workspaceId), &req, user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrUserAlreadyInWorkspace):
			util.WriteJSON(w, http.StatusConflict, map[string]any{
				"errorMessage": locale.TR(r, "user_already_workspace_member"),
			})
		case errors.Is(err, module.ErrPendingInviteExists):
			util.WriteJSON(w, http.StatusConflict, map[string]any{
				"errorMessage": locale.TR(r, "pending_invite_already_exists"),
			})
		default:
			log.Error().
				Err(err).
				Str("userId", user.Id.String()).
				Str("workspaceId", workspaceId).
				Msg("Failed to create invite")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_invite"),
			})
		}
		return
	}

	log.Info().
		Str("inviteId", invite.Id.String()).
		Str("userId", user.Id.String()).
		Str("workspaceId", workspaceId).
		Str("email", invite.Email).
		Msg("Invite created")

	util.WriteJSON(w, http.StatusCreated, invite)
}

// ListInvitesAction returns invites for a workspace.
func ListInvitesAction(w http.ResponseWriter, r *http.Request) {
	workspaceId := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(workspaceId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	log.Info().Str("workspaceId", workspaceId).Msg("Listing invites")

	limit, offset := util.ParsePagination(r)

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	result, err := im.ListInvites(db.Id(workspaceId), limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", workspaceId).Msg("Failed to list invites")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_invites"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"invites": result.Invites,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetInviteAction returns one invite by Id.
func GetInviteAction(w http.ResponseWriter, r *http.Request) {
	workspaceId := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(workspaceId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	inviteId := chi.URLParam(r, "inviteId")
	if lo.IsEmpty(inviteId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_invite_id"),
		})
		return
	}

	log.Info().
		Str("inviteId", inviteId).
		Str("workspaceId", workspaceId).
		Msg("Getting invite")

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	invite, err := im.GetInvite(db.Id(workspaceId), db.Id(inviteId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrInviteNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "invite_not_found"),
			})
		default:
			log.Error().
				Err(err).
				Str("inviteId", inviteId).
				Str("workspaceId", workspaceId).
				Msg("Failed to get invite")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_invite"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, invite)
}

// DeleteInviteAction deletes an invite by Id.
func DeleteInviteAction(w http.ResponseWriter, r *http.Request) {
	workspaceId := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(workspaceId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	inviteId := chi.URLParam(r, "inviteId")
	if lo.IsEmpty(inviteId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_invite_id"),
		})
		return
	}

	log.Info().
		Str("inviteId", inviteId).
		Str("workspaceId", workspaceId).
		Msg("Deleting invite")

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	err := im.DeleteInvite(db.Id(workspaceId), db.Id(inviteId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrInviteNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "invite_not_found"),
			})
		default:
			log.Error().Err(err).
				Str("inviteId", inviteId).
				Str("workspaceId", workspaceId).
				Msg("Failed to delete invite")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_invite"),
			})
		}
		return
	}

	log.Info().
		Str("inviteId", inviteId).
		Str("workspaceId", workspaceId).
		Msg("Invite deleted")

	w.WriteHeader(http.StatusNoContent)
}
