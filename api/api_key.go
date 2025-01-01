// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/middleware"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// CreateUserAPIKeyAction creates an API key; raw key is returned only once.
func CreateUserAPIKeyAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("New API key request")

	var req module.CreateAPIKeyRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	am := module.NewAPIKey(db.NewAPIKeyRepository(db.GetDB()))
	apiKey, err := am.CreateAPIKey(&req, user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrInvalidExpiresAt):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_expires_at_format"),
			})
			return
		default:
			log.Error().
				Err(err).
				Str("userId", user.Id.String()).
				Msg("Failed to create API key")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_api_key"),
			})
			return
		}
	}

	log.Info().
		Str("userId", user.Id.String()).
		Str("apiKeyId", apiKey.Id.String()).
		Msg("API key created")

	util.WriteJSON(w, http.StatusCreated, apiKey)
}

// ListUserAPIKeysAction lists your API keys (metadata only, never the secret).
func ListUserAPIKeysAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Listing API keys")

	limit, offset := util.ParsePagination(r)

	am := module.NewAPIKey(db.NewAPIKeyRepository(db.GetDB()))
	result, err := am.ListAPIKeys(user, limit, offset)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to list API keys")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_list_api_keys"),
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"keys": result.APIKeys,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetUserAPIKeyAction returns one API key's metadata (never the key itself).
func GetUserAPIKeyAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	apiKeyId := chi.URLParam(r, "apiKeyId")
	if lo.IsEmpty(apiKeyId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_api_key_id"),
		})
		return
	}

	log.Info().
		Str("apiKeyId", apiKeyId).
		Str("userId", user.Id.String()).
		Msg("Getting API key")

	am := module.NewAPIKey(db.NewAPIKeyRepository(db.GetDB()))
	k, err := am.GetAPIKey(db.Id(apiKeyId), user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrAPIKeyNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "api_key_not_found"),
			})
			return
		default:
			log.Error().
				Err(err).
				Str("apiKeyId", apiKeyId).
				Str("userId", user.Id.String()).
				Msg("Failed to get API key")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_api_key"),
			})
			return
		}
	}

	util.WriteJSON(w, http.StatusOK, k)
}

// DeleteUserAPIKeyAction deletes one of your API keys.
func DeleteUserAPIKeyAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	apiKeyId := chi.URLParam(r, "apiKeyId")
	if lo.IsEmpty(apiKeyId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_api_key_id"),
		})
		return
	}

	log.Info().
		Str("apiKeyId", apiKeyId).
		Str("userId", user.Id.String()).
		Msg("Deleting API key")

	am := module.NewAPIKey(db.NewAPIKeyRepository(db.GetDB()))
	err := am.DeleteAPIKey(db.Id(apiKeyId), user)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrAPIKeyNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "api_key_not_found"),
			})
			return
		default:
			log.Error().
				Err(err).
				Str("apiKeyId", apiKeyId).
				Str("userId", user.Id.String()).
				Msg("Failed to delete API key")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_api_key"),
			})
			return
		}
	}

	log.Info().
		Str("userId", user.Id.String()).
		Str("apiKeyId", apiKeyId).
		Msg("API key deleted")

	w.WriteHeader(http.StatusNoContent)
}
