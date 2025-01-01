// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/middleware"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/resend"
	"github.com/actx0/ziee/pkg/util"

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
		case errors.Is(err, module.ErrUserWithEmailRegistered):
			util.WriteJSON(w, http.StatusConflict, map[string]any{
				"errorMessage": locale.TR(r, "user_with_email_already_registered"),
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

// GetUserInviteByTokenAction returns invite details by token (used on the signup page).
func GetUserInviteByTokenAction(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if lo.IsEmpty(token) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_token"),
		})
		return
	}

	log.Info().
		Str("token", token).
		Msg("Getting invite by token")

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	invite, err := im.GetUserInviteByToken(token)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrInviteNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "invite_not_found"),
			})
		case errors.Is(err, module.ErrInviteNoLongerValid):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invite_no_longer_valid"),
			})
		default:
			log.Error().
				Err(err).
				Str("token", token).
				Msg("Failed to get invite by token")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_invite"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, invite)
}

// GetAuthenticatedUserInviteByTokenAction returns invite details for the current user.
func GetAuthenticatedUserInviteByTokenAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	token := chi.URLParam(r, "token")
	if lo.IsEmpty(token) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_token"),
		})
		return
	}

	log.Info().
		Str("token", token).
		Str("userId", user.Id.String()).
		Msg("Getting authenticated invite by token")

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	invite, err := im.GetAuthenticatedUserInviteByToken(token, user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrInviteNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "invite_not_found"),
			})
		case errors.Is(err, module.ErrInviteNoLongerValid):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invite_no_longer_valid"),
			})
		default:
			log.Error().
				Err(err).
				Str("token", token).
				Str("userId", user.Id.String()).
				Msg("Failed to get authenticated invite by token")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_invite"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, invite)
}

// ListUserInvitesAction returns invites addressed to the current user.
func ListUserInvitesAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Listing user invites")

	limit, offset := util.ParsePagination(r)

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	result, err := im.ListUserInvites(user, limit, offset)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to list invites")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_list_invites"),
		})
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

// AcceptUserInviteByTokenAction accepts an invite by token.
func AcceptUserInviteByTokenAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	token := chi.URLParam(r, "token")
	if lo.IsEmpty(token) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_token"),
		})
		return
	}

	log.Info().
		Str("token", token).
		Str("userId", user.Id.String()).
		Msg("Accepting invite")

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	invite, err := im.AcceptUserInviteByToken(token, user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrInviteNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "invite_not_found"),
			})
		case errors.Is(err, module.ErrInviteExpired):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invite_expired"),
			})
		case errors.Is(err, module.ErrInviteNoLongerValid):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invite_no_longer_valid"),
			})
		case errors.Is(err, module.ErrInviteEmailMismatch):
			util.WriteJSON(w, http.StatusForbidden, map[string]any{
				"errorMessage": locale.TR(r, "invite_email_mismatch"),
			})
		default:
			log.Error().
				Err(err).
				Str("token", token).
				Str("userId", user.Id.String()).
				Msg("Failed to accept invite")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_accept_invite"),
			})
		}
		return
	}

	log.Info().
		Str("inviteId", invite.Id.String()).
		Str("workspaceId", invite.WorkspaceId.String()).
		Str("userId", user.Id.String()).
		Msg("Invite accepted")

	util.WriteJSON(w, http.StatusOK, invite)
}

// RejectUserInviteByTokenAction rejects an invite by token.
func RejectUserInviteByTokenAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	token := chi.URLParam(r, "token")
	if lo.IsEmpty(token) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_token"),
		})
		return
	}

	log.Info().
		Str("token", token).
		Str("userId", user.Id.String()).
		Msg("Rejecting invite")

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	invite, err := im.RejectUserInviteByToken(token, user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrInviteNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "invite_not_found"),
			})
		case errors.Is(err, module.ErrInviteExpired):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invite_expired"),
			})
		case errors.Is(err, module.ErrInviteNoLongerValid):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invite_no_longer_valid"),
			})
		case errors.Is(err, module.ErrInviteEmailMismatch):
			util.WriteJSON(w, http.StatusForbidden, map[string]any{
				"errorMessage": locale.TR(r, "invite_email_mismatch"),
			})
		default:
			log.Error().
				Err(err).
				Str("token", token).
				Str("userId", user.Id.String()).
				Msg("Failed to reject invite")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_reject_invite"),
			})
		}
		return
	}

	log.Info().
		Str("inviteId", invite.Id.String()).
		Str("workspaceId", invite.WorkspaceId.String()).
		Str("userId", user.Id.String()).
		Msg("Invite rejected")

	util.WriteJSON(w, http.StatusOK, invite)
}
