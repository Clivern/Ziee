// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package agent

import (
	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/pkg/ai"
	"github.com/actx0/ziee/service/knowledge"
)

const (
	MemoryKindSummary    = "summary"
	MemoryKindFact       = "fact"
	MemoryKindPreference = "preference"
)

// RunRequest executes one agent turn.
type RunRequest struct {
	AgentId    db.Id
	ExternalId string
	Input      string
	Labels     map[string]string
	UseLarge   bool
}

// RunResponse is the result of an agent turn.
type RunResponse struct {
	Reply      string
	SessionId  db.Id
	ExternalId string
	MessageId  db.Id
}

// StoreMemoryRequest persists a session memory item.
type StoreMemoryRequest struct {
	ExternalId string
	AgentId    db.Id
	Kind       string
	Content    string
	Meta       *string
}

// ConsolidateMemoryRequest triggers async session memory consolidation.
type ConsolidateMemoryRequest struct {
	AgentId    db.Id
	ExternalId string
}

// RecallMemoryRequest finds memories relevant to a query.
type RecallMemoryRequest struct {
	ExternalId string
	AgentId    db.Id
	Query      string
	Kinds      []string
	Limit      uint64
}

// MemoryHit is a scored memory item.
type MemoryHit struct {
	MemoryId db.Id
	Kind     string
	Content  string
	Score    float32
}

// contextBundle is the assembled prompt context for a model call.
type contextBundle struct {
	Messages     []ai.Message
	Knowledge    []knowledge.SearchHit
	Memories     []MemoryHit
	SystemPrompt string
	HistoryCount int
}
