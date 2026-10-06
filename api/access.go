// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/locale"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// CreateAccessKeyAction creates a workspace access key; raw key is returned only once.
func (a *API) CreateAccessKeyAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if wid == "" {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Msg("New access key request")

	var req module.CreateAccessKeyRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	key, err := a.Access.CreateAccessKey(db.Id(wid), &req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrInvalidExpiresAt):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_expires_at_format"),
			})
		case errors.Is(err, module.ErrInvalidAccessKeyPermissions):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_access_key_permissions"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Msg("Failed to create workspace access key")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_access_key"),
			})
		}
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("accessKeyId", key.Id.String()).
		Msg("Access key created")

	util.WriteJSON(w, http.StatusCreated, key)
}

// ListAccessKeysAction lists workspace access keys (metadata only, never the secret).
func (a *API) ListAccessKeysAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if wid == "" {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Msg("Listing access keys")

	limit, offset := util.ParsePagination(r)

	result, err := a.Access.ListAccessKeys(db.Id(wid), limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Msg("Failed to list workspace access keys")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_access_keys"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"keys": result.Keys,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetAccessKeyAction returns one workspace access key (never the secret).
func (a *API) GetAccessKeyAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if wid == "" {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	keyId := chi.URLParam(r, "keyId")
	if keyId == "" {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_access_key_id"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("accessKeyId", keyId).
		Msg("Getting access key")

	key, err := a.Access.GetAccessKey(db.Id(wid), db.Id(keyId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAccessKeyNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "access_key_not_found"),
			})
		default:
			log.Error().Err(err).Str("accessKeyId", keyId).Msg("Failed to get workspace access key")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_access_key"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, key)
}

// DeleteAccessKeyAction deletes a workspace access key.
func (a *API) DeleteAccessKeyAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if wid == "" {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	keyId := chi.URLParam(r, "keyId")
	if keyId == "" {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_access_key_id"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("accessKeyId", keyId).
		Msg("Deleting access key")

	err := a.Access.DeleteAccessKey(db.Id(wid), db.Id(keyId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAccessKeyNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "access_key_not_found"),
			})
		default:
			log.Error().Err(err).Str("accessKeyId", keyId).Msg("Failed to delete workspace access key")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_access_key"),
			})
		}
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("accessKeyId", keyId).
		Msg("Access key deleted")

	w.WriteHeader(http.StatusNoContent)
}
