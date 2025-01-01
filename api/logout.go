// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/middleware"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/util"

	"github.com/rs/zerolog/log"
)

// LogoutAction logs the user out and revokes their session.
func LogoutAction(w http.ResponseWriter, r *http.Request) {
	util.DeleteCookie(w, "_ziee_session")

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Info().Msg("New logout request")
		util.WriteJSON(w, http.StatusOK, map[string]any{
			"successMessage": locale.TR(r, "logout_successful"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("New logout request")

	am := module.NewAuth(
		db.NewUserRepository(db.GetDB()),
		db.NewSessionRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
	)

	err := am.Logout(r.Context(), user.Id)
	if err != nil {
		log.Error().
			Str("userId", user.Id.String()).
			Err(err).
			Msg("Failed to revoke session")
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"successMessage": locale.TR(r, "logout_successful"),
	})
}
