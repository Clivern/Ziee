// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/broker"

	"github.com/rs/zerolog/log"
)

// HandleGitHubIssue evaluates an issue or issue comment webhook.
func (h *handlers) HandleGitHubIssue(_ context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])
	h.tasks.MarkRunning(taskId)

	log.Info().
		Str("deliveryId", payload["deliveryId"]).
		Str("event", payload["event"]).
		Str("action", payload["action"]).
		Str("owner", payload["owner"]).
		Str("repo", payload["repo"]).
		Str("number", payload["number"]).
		Msg("GitHub issue webhook queued")

	return h.tasks.Complete(taskId)
}
