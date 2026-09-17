// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/broker"
)

// HandleGitHubPullRequest evaluates a pull request webhook.
func (h *handlers) HandleGitHubPullRequest(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])
	h.tasks.MarkRunning(taskId)

	// TODO: load .ziee.yml and evaluate the pull request with eval.Run
	// (spam / pr_triage via EvaluatePROpened when wired)

	return h.tasks.Complete(taskId, "")
}
