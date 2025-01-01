// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/module"

	"github.com/rs/zerolog/log"
)

type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *StatusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// TrackWorkspaceAPICall increments workspace API usage after authenticated, authorized
// requests complete successfully. Place on workspace routes after Protect with Workspace: true.
func TrackWorkspaceAPICall() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &StatusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)

			if recorder.status < http.StatusOK || recorder.status >= http.StatusMultipleChoices {
				return
			}

			if _, _, ok := GetPrincipal(r); !ok {
				return
			}

			workspace, ok := GetWorkspaceFromContext(r.Context())
			if !ok {
				return
			}

			err := module.NewUsage().IncrementAPICalls(
				db.NewUsageRepository(db.GetDB(false)),
				workspace.Id,
			)
			if err != nil {
				log.Warn().
					Err(err).
					Str("workspaceId", workspace.Id.String()).
					Str("path", r.URL.Path).
					Msg("Failed to record workspace API call usage")
			}
		})
	}
}
