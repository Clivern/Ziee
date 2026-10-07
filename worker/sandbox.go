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

// HandleSandboxStart starts a sandbox for a repository.
func (h *handlers) HandleSandboxStart(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	// TODO: implement
	log.Info().
		Str("taskId", taskId.String()).
		Str("sandboxId", payload["sandboxId"]).
		Msg("Sandbox start task received")

	return h.Tasks.Complete(taskId, "")
}

// HandleSandboxQuery sends a query to a running sandbox.
func (h *handlers) HandleSandboxQuery(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	// TODO: implement
	log.Info().
		Str("taskId", taskId.String()).
		Str("sandboxId", payload["sandboxId"]).
		Msg("Sandbox query task received")

	return h.Tasks.Complete(taskId, "")
}

// HandleSandboxStop stops a sandbox container.
func (h *handlers) HandleSandboxStop(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	// TODO: implement
	log.Info().
		Str("taskId", taskId.String()).
		Str("sandboxId", payload["sandboxId"]).
		Msg("Sandbox stop task received")

	return h.Tasks.Complete(taskId, "")
}

// HandleSandboxRemove removes a sandbox container and its clone.
func (h *handlers) HandleSandboxRemove(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	// TODO: implement
	log.Info().
		Str("taskId", taskId.String()).
		Str("sandboxId", payload["sandboxId"]).
		Msg("Sandbox remove task received")

	return h.Tasks.Complete(taskId, "")
}

// HandleSandboxDeleteExpired deletes expired sandboxes.
func (h *handlers) HandleSandboxDeleteExpired(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	// TODO: implement
	log.Info().
		Str("taskId", taskId.String()).
		Msg("Sandbox delete expired task received")

	return h.Tasks.Complete(taskId, "")
}

// HandleSandboxDeleteRemoved deletes removed sandboxes.
func (h *handlers) HandleSandboxDeleteRemoved(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	// TODO: implement
	log.Info().
		Str("taskId", taskId.String()).
		Msg("Sandbox delete removed task received")

	return h.Tasks.Complete(taskId, "")
}
