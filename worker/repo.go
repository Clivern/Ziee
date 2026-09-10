// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/broker"
)

// HandleRepositoryBootstrap creates a setup issue and .ziee.yml PR for a repository.
func (h *handlers) HandleRepositoryBootstrap(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])
	h.tasks.MarkRunning(taskId)

	installationId, err := strconv.ParseInt(payload["installationId"], 10, 64)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	owner, _, _ := strings.Cut(payload["fullName"], "/")
	repo := payload["name"]

	// TODO: Implement this
	//
	// if app.Get().FileExists(ctx, installationID, owner, repo, ".ziee.yml") {
	//     return h.tasks.Complete(taskId)
	// }
	//
	// app.Get().CreateIssue(ctx, installationID, owner, repo, {
	//     title: "Set up Ziee for this repository",
	//     body:  "...",
	// })
	//
	// app.Get().CreatePullRequest(ctx, installationID, owner, repo, {
	//     branch:  "ziee/init",
	//     path:    ".ziee.yml",
	//     content: defaultZieeYml,
	//     title:   "Add .ziee.yml",
	//     body:    "...",
	// })

	_ = installationId
	_ = owner
	_ = repo

	return h.tasks.Complete(taskId)
}
