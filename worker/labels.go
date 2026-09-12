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
	"github.com/clivern/ziee/pkg/github/policy/spec"
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

	data, err := app.Get().GetFile(ctx, installationId, owner, repo, ".ziee.yml")
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	file, err := spec.Parse(data)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	existing, err := app.Get().ListLabels(ctx, installationId, owner, repo)
	if err != nil {
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
			h.tasks.Fail(taskId, err.Error())
			return err
		}

		created++
	}

	result, err := json.Marshal(map[string]int{"created": created})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	return h.tasks.Complete(taskId, string(result))
}
