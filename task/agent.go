// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package task

import (
	"context"

	"github.com/actx0/ziee/db"
)

// HandleAgentDelete removes all session message and memory vectors for an agent.
func (h *Handlers) HandleAgentDelete(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.DeleteAgent(ctx, db.Id(payload["agentId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "deleted",
	})
}

// HandleAgentSessionDelete removes all message and memory vectors for a session.
func (h *Handlers) HandleAgentSessionDelete(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.DeleteAgentSession(ctx, db.Id(payload["agentId"]), db.Id(payload["sessionId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "deleted",
	})
}

// HandleAgentSessionMessagesDelete removes all message vectors for a session.
func (h *Handlers) HandleAgentSessionMessagesDelete(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.DeleteSessionMessages(ctx, db.Id(payload["agentId"]), db.Id(payload["sessionId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "deleted",
	})
}

// HandleAgentSessionMessageDelete removes vectors for one message.
func (h *Handlers) HandleAgentSessionMessageDelete(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.DeleteSessionMessage(ctx, db.Id(payload["messageId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "deleted",
	})
}

// HandleAgentSessionMemoriesDelete removes all memory vectors for a session.
func (h *Handlers) HandleAgentSessionMemoriesDelete(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.DeleteSessionMemories(ctx, db.Id(payload["agentId"]), db.Id(payload["sessionId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "deleted",
	})
}

// HandleAgentSessionMemoryDelete removes vectors for one memory item.
func (h *Handlers) HandleAgentSessionMemoryDelete(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.DeleteSessionMemory(ctx, db.Id(payload["memoryId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "deleted",
	})
}

// HandleAgentSessionMessageIndex embeds a session message into Qdrant.
func (h *Handlers) HandleAgentSessionMessageIndex(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.IndexSessionMessage(ctx, db.Id(payload["messageId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "indexed",
	})
}

// HandleAgentSessionMemoryIndex embeds a session memory into Qdrant.
func (h *Handlers) HandleAgentSessionMemoryIndex(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.IndexSessionMemory(ctx, db.Id(payload["memoryId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "indexed",
	})
}

// HandleAgentSessionMemoryGenerate extracts memory from a session message.
func (h *Handlers) HandleAgentSessionMemoryGenerate(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.agent.GenerateMemoryFromMessage(ctx, db.Id(payload["messageId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "generated",
	})
}
