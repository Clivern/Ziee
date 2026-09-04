// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/locale"
	"github.com/clivern/ziee/middleware"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/util"

	"github.com/rs/zerolog/log"
)

// GetProfileAction returns the current user's profile.
func GetProfileAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Getting profile")

	pm := module.NewProfile(db.NewUserRepository(db.GetDB()))
	profile, err := pm.GetProfile(user.Id)

	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to get profile")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_get_profile"),
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"id":              profile.Id,
			"email":           profile.Email,
			"role":            profile.Role,
			"isActive":        profile.IsActive,
			"isEmailVerified": profile.IsEmailVerified,
			"name":            profile.Name,
			"avatar":          util.GravatarURL(profile.Email, 80),
			"language":        profile.Language,
			"theme":           profile.Theme,
			"lastLoginAt":     profile.LastLoginAt.UTC().Format(time.RFC3339),
			"createdAt":       profile.CreatedAt.UTC().Format(time.RFC3339),
			"updatedAt":       profile.UpdatedAt.UTC().Format(time.RFC3339),
		},
	})
}

// UpdateProfileAction updates the current user's profile.
func UpdateProfileAction(w http.ResponseWriter, r *http.Request) {
	var req module.UpdateProfileRequest

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Updating profile")

	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	pm := module.NewProfile(db.NewUserRepository(db.GetDB()))
	updated, err := pm.UpdateProfile(user.Id, &req)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrFailedUpdateProfile), errors.Is(err, module.ErrUserNotFound):
			log.Error().
				Err(err).
				Str("userId", user.Id.String()).
				Msg("Failed to update profile")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_profile"),
			})
		default:
			log.Error().
				Err(err).
				Str("userId", user.Id.String()).
				Msg("Failed to update profile")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_profile"),
			})
		}
		return
	}

	log.Info().
		Str("userId", updated.Id.String()).
		Msg("Profile updated")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"id":              updated.Id,
			"email":           updated.Email,
			"role":            updated.Role,
			"isActive":        updated.IsActive,
			"isEmailVerified": updated.IsEmailVerified,
			"name":            updated.Name,
			"avatar":          util.GravatarURL(updated.Email, 80),
			"language":        updated.Language,
			"theme":           updated.Theme,
			"lastLoginAt":     updated.LastLoginAt.UTC().Format(time.RFC3339),
			"createdAt":       updated.CreatedAt.UTC().Format(time.RFC3339),
			"updatedAt":       updated.UpdatedAt.UTC().Format(time.RFC3339),
		},
		"successMessage": locale.TR(r, "profile_updated_successfully"),
	})
}
