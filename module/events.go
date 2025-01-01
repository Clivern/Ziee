// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/pkg/event"

	"github.com/rs/zerolog/log"
)

const (
	EventUserLoggedIn  = "user.logged_in"
	EventUserLoggedOut = "user.logged_out"
)

// UserLoggedInData is emitted after a successful login.
type UserLoggedInData struct {
	UserID    db.Id
	SessionID db.Id
}

// UserLoggedOutData is emitted after a successful logout.
type UserLoggedOutData struct {
	UserID db.Id
}

var (
	UserLoggedIn  = event.New[UserLoggedInData](EventUserLoggedIn)
	UserLoggedOut = event.New[UserLoggedOutData](EventUserLoggedOut)
)

// RegisterEventListeners wires application event handlers.
func RegisterEventListeners() {
	UserLoggedIn.On(func(_ context.Context, e UserLoggedInData) {
		log.Debug().
			Str("event", EventUserLoggedIn).
			Str("userId", e.UserID.String()).
			Str("sessionId", e.SessionID.String()).
			Msg("Event emitted")
	})

	UserLoggedOut.On(func(_ context.Context, e UserLoggedOutData) {
		log.Debug().
			Str("event", EventUserLoggedOut).
			Str("userId", e.UserID.String()).
			Msg("Event emitted")
	})
}
