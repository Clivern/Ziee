// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/util"

	"github.com/samber/lo"
)

var (
	ErrSessionMemoryNotFound     = errors.New("session memory not found")
	ErrFailedCreateSessionMemory = errors.New("failed create session memory")
	ErrFailedListSessionMemories = errors.New("failed list session memories")
	ErrFailedGetSessionMemory    = errors.New("failed get session memory")
	ErrFailedUpdateSessionMemory = errors.New("failed update session memory")
	ErrFailedDeleteSessionMemory = errors.New("failed delete session memory")
)

// CreateSessionMemoryRequest is the body for creating a session memory.
type CreateSessionMemoryRequest struct {
	Kind    string  `json:"kind" validate:"required,oneof=summary fact preference" label:"Kind"`
	Content string  `json:"content" validate:"required" label:"Content"`
	Meta    *string `json:"meta" validate:"omitempty,json" label:"Meta"`
}

// BatchCreateSessionMemoriesRequest is the body for batch creating session memories.
type BatchCreateSessionMemoriesRequest struct {
	Memories []CreateSessionMemoryRequest `json:"memories" validate:"required,min=1,max=100,dive" label:"Memories"`
}

// BatchDeleteSessionMemoriesRequest is the body for batch deleting session memories.
type BatchDeleteSessionMemoriesRequest struct {
	IDs []string `json:"ids" validate:"required,min=1,max=100,dive,required" label:"IDs"`
}

// UpdateSessionMemoryRequest is the body for updating a session memory.
type UpdateSessionMemoryRequest struct {
	Kind    string  `json:"kind" validate:"omitempty,oneof=summary fact preference" label:"Kind"`
	Content string  `json:"content" validate:"required" label:"Content"`
	Meta    *string `json:"meta" validate:"omitempty,json" label:"Meta"`
}

// SessionMemoryResponse is a session memory shaped for API responses.
type SessionMemoryResponse struct {
	Id        db.Id           `json:"id"`
	SessionId db.Id           `json:"sessionId"`
	Kind      string          `json:"kind"`
	Content   string          `json:"content"`
	Meta      json.RawMessage `json:"meta,omitempty"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

// ListSessionMemoriesResponse is returned when listing session memories.
type ListSessionMemoriesResponse struct {
	Memories []*SessionMemoryResponse
	Total    int64
}

// Memory is the module for session memory CRUD.
type Memory struct {
	MemoryRepository db.SessionMemoryRepository
	AgentSession     *AgentSession
}

// NewMemory creates a session memory module with the given repositories.
func NewMemory(memories db.SessionMemoryRepository, agentSessions *AgentSession) *Memory {
	return &Memory{
		MemoryRepository: memories,
		AgentSession:     agentSessions,
	}
}

// CreateSessionMemory creates a memory in a session.
func (m *Memory) CreateSessionMemory(workspaceId, agentId, sessionId db.Id, req *CreateSessionMemoryRequest) (*SessionMemoryResponse, error) {
	session, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	memory := &db.SessionMemory{
		SessionId: session.Id,
		Kind:      req.Kind,
		Content:   req.Content,
		Meta:      req.Meta,
	}

	err = m.MemoryRepository.Create(memory)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSessionMemory, err)
	}

	err = m.EnqueueIndexTask(workspaceId, memory.Id)
	if err != nil {
		return nil, err
	}

	return &SessionMemoryResponse{
		Id:        memory.Id,
		SessionId: memory.SessionId,
		Kind:      memory.Kind,
		Content:   memory.Content,
		Meta:      util.JSONRawFromString(memory.Meta),
		CreatedAt: memory.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: memory.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// BatchCreateSessionMemories creates multiple memories in a session.
func (m *Memory) BatchCreateSessionMemories(workspaceId, agentId, sessionId db.Id, req *BatchCreateSessionMemoriesRequest) ([]*SessionMemoryResponse, error) {
	session, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	responses := make([]*SessionMemoryResponse, 0, len(req.Memories))
	for _, item := range req.Memories {
		memory := &db.SessionMemory{
			SessionId: session.Id,
			Kind:      item.Kind,
			Content:   item.Content,
			Meta:      item.Meta,
		}

		err = m.MemoryRepository.Create(memory)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedCreateSessionMemory, err)
		}

		err = m.EnqueueIndexTask(workspaceId, memory.Id)
		if err != nil {
			return nil, err
		}

		responses = append(responses, &SessionMemoryResponse{
			Id:        memory.Id,
			SessionId: memory.SessionId,
			Kind:      memory.Kind,
			Content:   memory.Content,
			Meta:      util.JSONRawFromString(memory.Meta),
			CreatedAt: memory.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: memory.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return responses, nil
}

// GetSessionMemory returns one memory by id.
func (m *Memory) GetSessionMemory(workspaceId, agentId, sessionId, memoryId db.Id) (*SessionMemoryResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	memory, err := m.MemoryRepository.GetById(memoryId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetSessionMemory, err)
	}
	if memory == nil || memory.SessionId != sessionId {
		return nil, ErrSessionMemoryNotFound
	}

	return &SessionMemoryResponse{
		Id:        memory.Id,
		SessionId: memory.SessionId,
		Kind:      memory.Kind,
		Content:   memory.Content,
		Meta:      util.JSONRawFromString(memory.Meta),
		CreatedAt: memory.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: memory.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// UpdateSessionMemory updates a memory by id.
func (m *Memory) UpdateSessionMemory(workspaceId, agentId, sessionId, memoryId db.Id, req *UpdateSessionMemoryRequest) (*SessionMemoryResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	memory, err := m.MemoryRepository.GetById(memoryId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateSessionMemory, err)
	}
	if memory == nil || memory.SessionId != sessionId {
		return nil, ErrSessionMemoryNotFound
	}

	if lo.IsNotEmpty(req.Kind) {
		memory.Kind = req.Kind
	}
	memory.Content = req.Content
	memory.Meta = req.Meta

	err = m.MemoryRepository.Update(memory)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateSessionMemory, err)
	}

	return &SessionMemoryResponse{
		Id:        memory.Id,
		SessionId: memory.SessionId,
		Kind:      memory.Kind,
		Content:   memory.Content,
		Meta:      util.JSONRawFromString(memory.Meta),
		CreatedAt: memory.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: memory.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// DeleteSessionMemory deletes a memory by id.
func (m *Memory) DeleteSessionMemory(workspaceId, agentId, sessionId, memoryId db.Id) error {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return err
	}

	memory, err := m.MemoryRepository.GetById(memoryId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMemory, err)
	}
	if memory == nil || memory.SessionId != sessionId {
		return ErrSessionMemoryNotFound
	}

	err = m.MemoryRepository.Delete(memoryId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMemory, err)
	}

	return m.EnqueueDeleteTask(workspaceId, memoryId)
}

// BatchDeleteSessionMemories deletes multiple memories by id.
func (m *Memory) BatchDeleteSessionMemories(workspaceId, agentId, sessionId db.Id, req *BatchDeleteSessionMemoriesRequest) error {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return err
	}

	for _, id := range req.IDs {
		memoryId := db.Id(id)

		memory, err := m.MemoryRepository.GetById(memoryId)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMemory, err)
		}
		if memory == nil || memory.SessionId != sessionId {
			return ErrSessionMemoryNotFound
		}

		err = m.MemoryRepository.Delete(memoryId)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMemory, err)
		}

		err = m.EnqueueDeleteTask(workspaceId, memoryId)
		if err != nil {
			return err
		}
	}

	return nil
}

// ListSessionMemories returns paginated memories for a session.
func (m *Memory) ListSessionMemories(workspaceId, agentId, sessionId db.Id, limit, offset int) (*ListSessionMemoriesResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	total, err := m.MemoryRepository.CountBySessionId(sessionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListSessionMemories, err)
	}

	memories, err := m.MemoryRepository.ListBySessionId(sessionId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListSessionMemories, err)
	}

	list := make([]*SessionMemoryResponse, 0, len(memories))
	for _, item := range memories {
		list = append(list, &SessionMemoryResponse{
			Id:        item.Id,
			SessionId: item.SessionId,
			Kind:      item.Kind,
			Content:   item.Content,
			Meta:      util.JSONRawFromString(item.Meta),
			CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: item.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &ListSessionMemoriesResponse{Memories: list, Total: total}, nil
}

// EnqueueIndexTask enqueues an async task to index a session memory.
func (m *Memory) EnqueueIndexTask(workspaceId, memoryId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.memory.index",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"memoryId": memoryId.String(),
		},
	})

	return err
}

// EnqueueDeleteTask enqueues an async task to delete session memory vectors.
func (m *Memory) EnqueueDeleteTask(workspaceId, memoryId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.memory.delete",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"memoryId": memoryId.String(),
		},
	})

	return err
}
