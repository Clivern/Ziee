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
	"github.com/clivern/ziee/pkg/github/app"
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

	exists, err := app.Get().FileExists(ctx, installationId, owner, repo, ".ziee.yml")
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}
	if exists {
		return h.tasks.Complete(taskId)
	}

	_, err = app.Get().CreateIssue(ctx, installationId, owner, repo,
		"Set up Ziee for this repository",
		"Ziee is installed on this repository. Merge the `.ziee.yml` pull request to start enforcing policy.",
	)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	_, err = app.Get().CreatePullRequest(ctx, installationId, owner, repo, app.NewPullRequest{
		Branch:  "ziee/init",
		Path:    ".ziee.yml",
		Content: app.DefaultZieeYML,
		Title:   "Add .ziee.yml",
		Body:    "Adds a starter `.ziee.yml` so Ziee can enforce repository policy.",
	})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	return h.tasks.Complete(taskId)
}
