// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package task

import (
	"context"
	"fmt"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/module"
)

// Knowledge indexes and deletes workspace documents.
type Knowledge interface {
	Index(ctx context.Context, documentId db.Id) error
	Delete(ctx context.Context, documentId db.Id, internalId string) error
}

// Agent removes agent session vectors from Qdrant and indexes session content.
type Agent interface {
	DeleteAgent(ctx context.Context, agentId db.Id) error
	DeleteAgentSession(ctx context.Context, agentId, sessionId db.Id) error
	DeleteSessionMessages(ctx context.Context, agentId, sessionId db.Id) error
	DeleteSessionMessage(ctx context.Context, messageId db.Id) error
	DeleteSessionMemories(ctx context.Context, agentId, sessionId db.Id) error
	DeleteSessionMemory(ctx context.Context, memoryId db.Id) error
	IndexSessionMessage(ctx context.Context, messageId db.Id) error
	IndexSessionMemory(ctx context.Context, memoryId db.Id) error
	GenerateMemoryFromMessage(ctx context.Context, messageId db.Id) error
}

// Dependencies are the services required by async task handlers.
type Dependencies struct {
	Knowledge Knowledge
	Agent     Agent
}

// Handlers registers and runs async task handlers.
type Handlers struct {
	knowledge Knowledge
	agent     Agent
}

// New returns task handlers backed by services only.
func New(deps Dependencies) *Handlers {
	return &Handlers{
		knowledge: deps.Knowledge,
		agent:     deps.Agent,
	}
}

// Register attaches all task handlers to the async manager.
func Register(am *module.Async, deps Dependencies) {
	New(deps).Register(am)
}

// Register attaches handlers to the async module.
func (h *Handlers) Register(am *module.Async) {
	am.RegisterHandler("task.doc.index", h.HandleDocumentIndex)
	am.RegisterHandler("task.doc.delete", h.HandleDocumentDelete)
	am.RegisterHandler("task.agent.delete", h.HandleAgentDelete)
	am.RegisterHandler("task.agent.session.delete", h.HandleAgentSessionDelete)
	am.RegisterHandler("task.agent.session.messages.delete", h.HandleAgentSessionMessagesDelete)
	am.RegisterHandler("task.agent.session.message.delete", h.HandleAgentSessionMessageDelete)
	am.RegisterHandler("task.agent.session.memories.delete", h.HandleAgentSessionMemoriesDelete)
	am.RegisterHandler("task.agent.session.memory.delete", h.HandleAgentSessionMemoryDelete)
	am.RegisterHandler("task.agent.session.message.index", h.HandleAgentSessionMessageIndex)
	am.RegisterHandler("task.agent.session.memory.index", h.HandleAgentSessionMemoryIndex)
	am.RegisterHandler("task.agent.session.memory.generate", h.HandleAgentSessionMemoryGenerate)
}

// Run executes a handler synchronously without HTTP or the worker pool.
func (h *Handlers) Run(ctx context.Context, taskType string, task *db.AsyncTask) (string, error) {
	switch taskType {
	case "task.doc.index":
		return h.HandleDocumentIndex(ctx, task)
	case "task.doc.delete":
		return h.HandleDocumentDelete(ctx, task)
	case "task.agent.delete":
		return h.HandleAgentDelete(ctx, task)
	case "task.agent.session.delete":
		return h.HandleAgentSessionDelete(ctx, task)
	case "task.agent.session.messages.delete":
		return h.HandleAgentSessionMessagesDelete(ctx, task)
	case "task.agent.session.message.delete":
		return h.HandleAgentSessionMessageDelete(ctx, task)
	case "task.agent.session.memories.delete":
		return h.HandleAgentSessionMemoriesDelete(ctx, task)
	case "task.agent.session.memory.delete":
		return h.HandleAgentSessionMemoryDelete(ctx, task)
	case "task.agent.session.message.index":
		return h.HandleAgentSessionMessageIndex(ctx, task)
	case "task.agent.session.memory.index":
		return h.HandleAgentSessionMemoryIndex(ctx, task)
	case "task.agent.session.memory.generate":
		return h.HandleAgentSessionMemoryGenerate(ctx, task)
	default:
		return "", fmt.Errorf("unknown task type: %s", taskType)
	}
}
