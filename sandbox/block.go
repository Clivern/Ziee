// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"sync"

	"github.com/rs/zerolog/log"
)

var blocked sync.Map

// Block marks a sandbox id as blocked.
func Block(id string) {
	blocked.Store(id, struct{}{})

	log.Info().
		Str("id", id).
		Msg("Sandbox id blocked")
}

// IsBlocked reports whether a sandbox id is blocked.
func IsBlocked(id string) bool {
	_, ok := blocked.Load(id)

	return ok
}
