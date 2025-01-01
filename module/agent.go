// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/pkg/util"

	"github.com/samber/lo"
)

var (
	ErrAgentNotFound     = errors.New("agent not found")
	ErrFailedCreateAgent = errors.New("failed create agent")
	ErrFailedListAgents  = errors.New("failed list agents")
	ErrFailedGetAgent    = errors.New("failed get agent")
	ErrFailedUpdateAgent = errors.New("failed update agent")
	ErrFailedDeleteAgent = errors.New("failed delete agent")
)

// Agent is the module for workspace agent CRUD.
type Agent struct {
	AgentRepository     db.AgentRepository
	WorkspaceRepository db.WorkspaceRepository
}

// NewAgent creates an agent module with the given repositories.
func NewAgent(agents db.AgentRepository, workspaces db.WorkspaceRepository) *Agent {
	return &Agent{
		AgentRepository:     agents,
		WorkspaceRepository: workspaces,
	}
}

// CreateAgentRequest is the body for creating an agent.
type CreateAgentRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100" label:"Name"`
	Description string `json:"description" validate:"required,min=1,max=2000" label:"Description"`
}

// UpdateAgentRequest is the body for updating an agent.
type UpdateAgentRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100" label:"Name"`
	Description string `json:"description" validate:"required,min=1,max=2000" label:"Description"`
}

// AgentResponse is an agent shaped for API responses.
type AgentResponse struct {
	Id          db.Id             `json:"id"`
	WorkspaceId db.Id             `json:"workspaceId"`
	Name        string            `json:"name"`
	Kind        string            `json:"kind"`
	PromptId    *db.Id            `json:"promptId,omitempty"`
	KbLabels    map[string]string `json:"kbLabels,omitempty"`
	Handle      string            `json:"handle"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
}

// ListAgentsResponse is returned when listing agents.
type ListAgentsResponse struct {
	Agents []*AgentResponse
	Total  int64
}

// CreateAgent creates a new agent in a workspace.
func (a *Agent) CreateAgent(workspaceId db.Id, req *CreateAgentRequest) (*AgentResponse, error) {
	workspace, err := a.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	handle, err := a.CreateAgentHandle(4, 8)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateAgent, err)
	}

	agent := &db.Agent{
		WorkspaceId: workspaceId,
		Name:        req.Name,
		Handle:      handle,
		Description: req.Description,
	}

	err = a.AgentRepository.Create(agent)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateAgent, err)
	}

	return &AgentResponse{
		Id:          agent.Id,
		WorkspaceId: agent.WorkspaceId,
		Name:        agent.Name,
		Kind:        agent.Kind,
		PromptId:    agent.PromptId,
		Handle:      agent.Handle,
		Description: agent.Description,
		Status:      agent.Status,
		CreatedAt:   agent.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   agent.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// ListAgents returns paginated agents for a workspace.
func (a *Agent) ListAgents(workspaceId db.Id, limit, offset int) (*ListAgentsResponse, error) {
	workspace, err := a.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	total, err := a.AgentRepository.CountByWorkspaceId(workspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListAgents, err)
	}

	agents, err := a.AgentRepository.ListByWorkspaceId(workspaceId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListAgents, err)
	}

	list := make([]*AgentResponse, 0, len(agents))
	for _, item := range agents {
		var kbLabels map[string]string
		if item.KbLabels != nil && lo.IsNotEmpty(*item.KbLabels) {
			_ = json.Unmarshal([]byte(*item.KbLabels), &kbLabels)
		}

		list = append(list, &AgentResponse{
			Id:          item.Id,
			WorkspaceId: item.WorkspaceId,
			Name:        item.Name,
			Kind:        item.Kind,
			PromptId:    item.PromptId,
			KbLabels:    kbLabels,
			Handle:      item.Handle,
			Description: item.Description,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &ListAgentsResponse{Agents: list, Total: total}, nil
}

// GetAgent returns one agent by Id.
func (a *Agent) GetAgent(workspaceId, agentId db.Id) (*AgentResponse, error) {
	workspace, err := a.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := a.AgentRepository.GetById(agentId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetAgent, err)
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return nil, ErrAgentNotFound
	}

	var kbLabels map[string]string
	if agent.KbLabels != nil && lo.IsNotEmpty(*agent.KbLabels) {
		_ = json.Unmarshal([]byte(*agent.KbLabels), &kbLabels)
	}

	return &AgentResponse{
		Id:          agent.Id,
		WorkspaceId: agent.WorkspaceId,
		Name:        agent.Name,
		Kind:        agent.Kind,
		PromptId:    agent.PromptId,
		KbLabels:    kbLabels,
		Handle:      agent.Handle,
		Description: agent.Description,
		Status:      agent.Status,
		CreatedAt:   agent.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   agent.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// UpdateAgent updates an agent by Id.
func (a *Agent) UpdateAgent(workspaceId, agentId db.Id, req *UpdateAgentRequest) (*AgentResponse, error) {
	workspace, err := a.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := a.AgentRepository.GetById(agentId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateAgent, err)
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return nil, ErrAgentNotFound
	}

	agent.Name = req.Name
	agent.Description = req.Description

	err = a.AgentRepository.Update(agent)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateAgent, err)
	}

	var kbLabels map[string]string
	if agent.KbLabels != nil && lo.IsNotEmpty(*agent.KbLabels) {
		_ = json.Unmarshal([]byte(*agent.KbLabels), &kbLabels)
	}

	return &AgentResponse{
		Id:          agent.Id,
		WorkspaceId: agent.WorkspaceId,
		Name:        agent.Name,
		Kind:        agent.Kind,
		PromptId:    agent.PromptId,
		KbLabels:    kbLabels,
		Handle:      agent.Handle,
		Description: agent.Description,
		Status:      agent.Status,
		CreatedAt:   agent.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   agent.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// DeleteAgent deletes an agent by Id.
func (a *Agent) DeleteAgent(workspaceId, agentId db.Id) error {
	workspace, err := a.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	agent, err := a.AgentRepository.GetById(agentId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteAgent, err)
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return ErrAgentNotFound
	}

	err = a.AgentRepository.Delete(agentId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteAgent, err)
	}

	return a.EnqueueDeleteTask(workspaceId, agentId)
}

// EnqueueDeleteTask enqueues an async task to delete a resource.
func (a *Agent) EnqueueDeleteTask(workspaceId, agentId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.delete",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"agentId": agentId.String(),
		},
	})

	return err
}

// CreateAgentHandle creates a unique handle like abcd-1234 (min–max letters, dash, min–max digits).
func (a *Agent) CreateAgentHandle(min, max int) (string, error) {
	for {
		handle, err := util.RandomHandle(min, max)
		if err != nil {
			return "", err
		}
		existing, err := a.AgentRepository.GetByHandle(handle)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return handle, nil
		}
	}
}
