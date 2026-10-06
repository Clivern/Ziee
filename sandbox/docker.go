// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"context"
	"os/exec"

	"github.com/rs/zerolog/log"
)

// IsDockerRunning reports whether the Docker daemon is reachable.
func IsDockerRunning(ctx context.Context) bool {
	err := exec.CommandContext(ctx, "docker", "info").Run()
	if err != nil {
		log.Debug().
			Err(err).
			Msg("Docker info failed")
		return false
	}

	return true
}
