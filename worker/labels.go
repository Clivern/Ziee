// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/broker"
	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/pkg/github/policy/spec"

	"github.com/rs/zerolog/log"
)

// HandleRepositoryLabels creates GitHub labels defined in `.ziee.yml`.
func (h *handlers) HandleRepositoryLabels(ctx context.Context, msg *broker.Msg) error {
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
	path := payload["path"]
	sha := payload["sha"]

	log.Info().
		Str("taskId", taskId.String()).
		Int64("installationId", installationId).
		Str("owner", owner).
		Str("repo", repo).
		Str("path", path).
		Str("sha", sha).
		Msg("Repository labels sync started")

	check, err := app.Get().CreateCheckRun(ctx, installationId, owner, repo, conf.ConfigSyncCheckName, sha)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	data, err := app.Get().GetFile(ctx, installationId, owner, repo, path)
	if err != nil {
		app.Get().CompleteCheckRun(
			ctx, installationId, owner, repo, check.ID,
			"failure",
			"Config sync failed",
			err.Error(),
		)
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	file, err := spec.Parse(data)
	if err != nil {
		app.Get().CompleteCheckRun(
			ctx, installationId, owner, repo, check.ID,
			"failure",
			"Config sync failed",
			err.Error(),
		)
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	existing, err := app.Get().ListLabels(ctx, installationId, owner, repo)
	if err != nil {
		app.Get().CompleteCheckRun(
			ctx, installationId, owner, repo, check.ID,
			"failure",
			"Config sync failed",
			err.Error(),
		)
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	have := map[string]struct{}{}
	for _, label := range existing {
		have[strings.ToLower(label.Name)] = struct{}{}
	}

	created := 0
	for _, label := range file.Labels {
		if _, ok := have[strings.ToLower(label.Name)]; ok {
			continue
		}

		err = app.Get().CreateLabel(ctx, installationId, owner, repo, app.Label{
			Name:        label.Name,
			Color:       label.Color,
			Description: label.Description,
		})
		if err != nil {
			app.Get().CompleteCheckRun(
				ctx, installationId, owner, repo, check.ID,
				"failure",
				"Config sync failed",
				err.Error(),
			)
			h.tasks.Fail(taskId, err.Error())
			return err
		}

		created++

		log.Info().
			Str("owner", owner).
			Str("repo", repo).
			Str("label", label.Name).
			Msg("GitHub label created")
	}

	err = app.Get().CompleteCheckRun(
		ctx, installationId, owner, repo, check.ID,
		"success",
		"Config synced",
		fmt.Sprintf("Created %d labels from `%s`.", created, path),
	)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	result, err := json.Marshal(map[string]int{"created": created})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	log.Info().
		Str("owner", owner).
		Str("repo", repo).
		Int("created", created).
		Int("defined", len(file.Labels)).
		Msg("Repository labels synced")

	return h.tasks.Complete(taskId, string(result))
}
