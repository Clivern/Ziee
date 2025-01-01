// Copyright 2026 Clivern. All rights reserved.
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
)

// CreateWorkspaceAccessKeyAction creates a workspace access key; raw key is returned only once.
func CreateWorkspaceAccessKeyAction(w http.ResponseWriter, r *http.Request) {
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

	am := module.NewAccess(
		db.NewWorkspaceAccessKeyRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)

	key, err := am.CreateAccessKey(db.Id(wid), &req)
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

	util.WriteJSON(w, http.StatusCreated, key)
}

// ListWorkspaceAccessKeysAction lists workspace access keys (metadata only, never the secret).
func ListWorkspaceAccessKeysAction(w http.ResponseWriter, r *http.Request) {
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

	am := module.NewAccess(
		db.NewWorkspaceAccessKeyRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)

	result, err := am.ListAccessKeys(db.Id(wid), limit, offset)
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

// GetWorkspaceAccessKeyAction returns one workspace access key (never the secret).
func GetWorkspaceAccessKeyAction(w http.ResponseWriter, r *http.Request) {
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

	am := module.NewAccess(
		db.NewWorkspaceAccessKeyRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)

	key, err := am.GetAccessKey(db.Id(wid), db.Id(keyId))
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

// DeleteWorkspaceAccessKeyAction deletes a workspace access key.
func DeleteWorkspaceAccessKeyAction(w http.ResponseWriter, r *http.Request) {
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

	am := module.NewAccess(
		db.NewWorkspaceAccessKeyRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)

	err := am.DeleteAccessKey(db.Id(wid), db.Id(keyId))
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

	w.WriteHeader(http.StatusNoContent)
}
