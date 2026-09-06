// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"context"
	"encoding/json"

	"github.com/clivern/ziee/pkg/event"
)

const (
	EventReceived                 = "github.webhook"
	EventInstallation             = "github.installation"
	EventInstallationRepositories = "github.installation_repositories"
)

var (
	// Received is emitted for every verified GitHub webhook delivery.
	Received = event.New[Delivery](EventReceived)
	// Installation is emitted for GitHub App installation events.
	Installation = event.New[InstallationEvent](EventInstallation)
	// InstallationRepositories is emitted for installation repository access changes.
	InstallationRepositories = event.New[InstallationRepositoriesEvent](EventInstallationRepositories)
)

// Dispatch emits typed events for a verified webhook delivery.
func (d Delivery) Dispatch(ctx context.Context) {
	Received.Emit(ctx, d)

	switch d.Event {
	case "installation":
		var payload InstallationEvent
		json.Unmarshal(d.Body, &payload)
		Installation.Emit(ctx, payload)
	case "installation_repositories":
		var payload InstallationRepositoriesEvent
		json.Unmarshal(d.Body, &payload)
		InstallationRepositories.Emit(ctx, payload)
	}
}
