// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/broker"

	"github.com/rs/zerolog/log"
)

// HandleRepositoryMerge merges queued pull requests for a repository.
func (h *handlers) HandleRepositoryMerge(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])
	repoId := db.Id(payload["repositoryId"])

	h.Tasks.MarkRunning(taskId)

	repo, err := h.Repository.RepoRepository.GetById(repoId)
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	queued, err := h.PQueue.ListQueued(repoId, 1000)
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	snapshot := module.SnapshotWorkingQueue(queued)

	log.Info().
		Str("taskId", taskId.String()).
		Str("repositoryId", repoId.String()).
		Str("owner", repo.Owner).
		Str("repo", repo.Name).
		Int("queued", len(snapshot)).
		Msg("Repository merge started")

	result, err := json.Marshal(map[string]any{
		"repositoryId": repoId.String(),
		"queued":       len(snapshot),
	})
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	return h.Tasks.Complete(taskId, string(result))
}
