// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package mcp

import (
	"context"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type MemoryMessage struct {
	Role    string `json:"role" jsonschema:"Message role (user or assistant)"`
	Content string `json:"content" jsonschema:"Message content"`
}

type MemoryScopeInput struct {
	UserID  string `json:"user_id,omitempty" jsonschema:"User identifier"`
	AgentID string `json:"agent_id,omitempty" jsonschema:"Agent identifier"`
	AppID   string `json:"app_id,omitempty" jsonschema:"App identifier"`
	RunID   string `json:"run_id,omitempty" jsonschema:"Run or session identifier"`
}

type AddMemoryInput struct {
	MemoryScopeInput
	Content  string          `json:"content,omitempty" jsonschema:"Text content to store"`
	Messages []MemoryMessage `json:"messages,omitempty" jsonschema:"Conversation messages to store"`
	Metadata map[string]any  `json:"metadata,omitempty" jsonschema:"Optional metadata"`
}

type AddMemoryOutput struct {
	EventID  string `json:"event_id"`
	MemoryID string `json:"memory_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type SearchMemoriesInput struct {
	MemoryScopeInput
	Query   string         `json:"query" jsonschema:"required,Semantic search query"`
	TopK    int            `json:"top_k,omitempty" jsonschema:"Maximum number of results"`
	Filters map[string]any `json:"filters,omitempty" jsonschema:"Additional filters"`
}

type MemoryRecord struct {
	MemoryID  string         `json:"memory_id"`
	Memory    string         `json:"memory"`
	UserID    string         `json:"user_id,omitempty"`
	AgentID   string         `json:"agent_id,omitempty"`
	AppID     string         `json:"app_id,omitempty"`
	RunID     string         `json:"run_id,omitempty"`
	Score     float32        `json:"score,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

type SearchMemoriesOutput struct {
	Results []MemoryRecord `json:"results"`
}

type GetMemoriesInput struct {
	MemoryScopeInput
	Limit  int `json:"limit,omitempty" jsonschema:"Page size"`
	Offset int `json:"offset,omitempty" jsonschema:"Page offset"`
}

type GetMemoriesOutput struct {
	Memories []MemoryRecord `json:"memories"`
	Total    int            `json:"total"`
	Limit    int            `json:"limit"`
	Offset   int            `json:"offset"`
}

type GetMemoryInput struct {
	MemoryID string `json:"memory_id" jsonschema:"required,Memory identifier"`
}

type UpdateMemoryInput struct {
	MemoryID string `json:"memory_id" jsonschema:"required,Memory identifier"`
	Content  string `json:"content" jsonschema:"required,Updated memory text"`
}

type UpdateMemoryOutput struct {
	MemoryID string `json:"memory_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type DeleteMemoryInput struct {
	MemoryID string `json:"memory_id" jsonschema:"required,Memory identifier"`
}

type DeleteMemoryOutput struct {
	MemoryID string `json:"memory_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type DeleteAllMemoriesInput struct {
	MemoryScopeInput
}

type DeleteAllMemoriesOutput struct {
	Deleted int    `json:"deleted"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DeleteEntitiesInput struct {
	MemoryScopeInput
	EntityType string `json:"entity_type,omitempty" jsonschema:"Entity type: user, agent, app, or run"`
}

type DeleteEntitiesOutput struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ListEntitiesOutput struct {
	Users  []string `json:"users"`
	Agents []string `json:"agents"`
	Apps   []string `json:"apps"`
	Runs   []string `json:"runs"`
}

type ListEventsInput struct {
	MemoryScopeInput
	Limit  int `json:"limit,omitempty" jsonschema:"Page size"`
	Offset int `json:"offset,omitempty" jsonschema:"Page offset"`
}

type MemoryEvent struct {
	EventID   string `json:"event_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ListEventsOutput struct {
	Events []MemoryEvent `json:"events"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

type GetEventStatusInput struct {
	EventID string `json:"event_id" jsonschema:"required,Async operation event identifier"`
}

type GetEventStatusOutput struct {
	EventID   string `json:"event_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	UpdatedAt string `json:"updated_at"`
}

func RegisterMemoryTools(server *mcpsdk.Server) {
	now := time.Now().UTC().Format(time.RFC3339)

	sampleMemory := MemoryRecord{
		MemoryID:  "mem_fake_001",
		Memory:    "User prefers concise answers and morning standups at 9am.",
		UserID:    "user_demo",
		AgentID:   "agent_demo",
		AppID:     "app_demo",
		RunID:     "run_demo",
		Score:     0.92,
		Metadata:  map[string]any{"source": "stub"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "add_memory",
		Description: "Save text or conversation history for a user/agent",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, _ AddMemoryInput) (*mcpsdk.CallToolResult, AddMemoryOutput, error) {
		return nil, AddMemoryOutput{
			EventID:  "evt_fake_add_001",
			MemoryID: "mem_fake_001",
			Status:   "completed",
			Message:  "Stub memory saved (not persisted)",
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "search_memories",
		Description: "Semantic search across existing memories with filters",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, _ SearchMemoriesInput) (*mcpsdk.CallToolResult, SearchMemoriesOutput, error) {
		return nil, SearchMemoriesOutput{Results: []MemoryRecord{sampleMemory}}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "get_memories",
		Description: "List memories with structured filters and pagination",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, input GetMemoriesInput) (*mcpsdk.CallToolResult, GetMemoriesOutput, error) {
		return nil, GetMemoriesOutput{
			Memories: []MemoryRecord{sampleMemory},
			Total:    1,
			Limit:    input.Limit,
			Offset:   input.Offset,
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "get_memory",
		Description: "Retrieve one memory by its memory_id",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, input GetMemoryInput) (*mcpsdk.CallToolResult, MemoryRecord, error) {
		record := sampleMemory
		record.MemoryID = input.MemoryID
		return nil, record, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "update_memory",
		Description: "Overwrite a memory's text after confirming the ID",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, input UpdateMemoryInput) (*mcpsdk.CallToolResult, UpdateMemoryOutput, error) {
		return nil, UpdateMemoryOutput{
			MemoryID: input.MemoryID,
			Status:   "completed",
			Message:  "Stub memory updated (not persisted)",
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "delete_memory",
		Description: "Delete a single memory by memory_id",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, input DeleteMemoryInput) (*mcpsdk.CallToolResult, DeleteMemoryOutput, error) {
		return nil, DeleteMemoryOutput{
			MemoryID: input.MemoryID,
			Status:   "completed",
			Message:  "Stub memory deleted (not persisted)",
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "delete_all_memories",
		Description: "Bulk delete all memories in scope",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, _ DeleteAllMemoriesInput) (*mcpsdk.CallToolResult, DeleteAllMemoriesOutput, error) {
		return nil, DeleteAllMemoriesOutput{
			Deleted: 1,
			Status:  "completed",
			Message: "Stub memories deleted (not persisted)",
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "delete_entities",
		Description: "Delete a user/agent/app/run entity and its memories",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, _ DeleteEntitiesInput) (*mcpsdk.CallToolResult, DeleteEntitiesOutput, error) {
		return nil, DeleteEntitiesOutput{
			Status:  "completed",
			Message: "Stub entity deleted (not persisted)",
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "list_entities",
		Description: "Enumerate users/agents/apps/runs stored in Mem0",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, _ struct{}) (*mcpsdk.CallToolResult, ListEntitiesOutput, error) {
		return nil, ListEntitiesOutput{
			Users:  []string{"user_demo"},
			Agents: []string{"agent_demo"},
			Apps:   []string{"app_demo"},
			Runs:   []string{"run_demo"},
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "list_events",
		Description: "List memory operation events with filters and pagination",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, input ListEventsInput) (*mcpsdk.CallToolResult, ListEventsOutput, error) {
		return nil, ListEventsOutput{
			Events: []MemoryEvent{{
				EventID:   "evt_fake_add_001",
				Type:      "add_memory",
				Status:    "completed",
				CreatedAt: now,
			}},
			Total:  1,
			Limit:  input.Limit,
			Offset: input.Offset,
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "get_event_status",
		Description: "Check the status of an async memory operation by event_id",
	}, func(_ context.Context, _ *mcpsdk.CallToolRequest, input GetEventStatusInput) (*mcpsdk.CallToolResult, GetEventStatusOutput, error) {
		return nil, GetEventStatusOutput{
			EventID:   input.EventID,
			Status:    "completed",
			Message:   "Stub event completed (not persisted)",
			UpdatedAt: now,
		}, nil
	})
}
