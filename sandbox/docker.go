// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"context"
	"os/exec"
)

// IsDockerRunning reports whether the Docker daemon is reachable.
func IsDockerRunning(ctx context.Context) bool {
	return exec.CommandContext(ctx, "docker", "info").Run() == nil
}
