// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.
//
// Agent service: session turns, context assembly, and session memory.
//
//	Run
//	  → buildContext (history + knowledge.Search + RecallMemory)
//	  → ai.ChatClient.Complete

package agent

import (
	"context"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/pkg/ai"
	"github.com/actx0/ziee/pkg/qdrant"
)

// Dependencies are the collaborators required by the agent service.
type Dependencies struct {
	Embed    *ai.EmbedClient
	Vectors  *qdrant.Client
	LiteLLM  *ai.ChatClient
	LargeLLM *ai.ChatClient
	Agents   db.AgentRepository
	Sessions db.AgentSessionRepository
	Messages db.SessionMessageRepository
	Memories db.SessionMemoryRepository
}

// Service orchestrates agent turns, context, and session memory.
type Service struct {
	embed    *ai.EmbedClient
	vectors  *qdrant.Client
	liteLLM  *ai.ChatClient
	largeLLM *ai.ChatClient
	agents   db.AgentRepository
	sessions db.AgentSessionRepository
	messages db.SessionMessageRepository
	memories db.SessionMemoryRepository
}

// New returns an agent service.
func New(deps Dependencies) *Service {
	return &Service{
		embed:    deps.Embed,
		vectors:  deps.Vectors,
		liteLLM:  deps.LiteLLM,
		largeLLM: deps.LargeLLM,
		agents:   deps.Agents,
		sessions: deps.Sessions,
		messages: deps.Messages,
		memories: deps.Memories,
	}
}

// Run executes one agent turn and returns the assistant reply.
func (s *Service) Run(ctx context.Context, req RunRequest) (*RunResponse, error) {
	return nil, ErrNotImplemented
}
