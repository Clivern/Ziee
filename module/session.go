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

	"github.com/google/uuid"
	"github.com/samber/lo"
)

var (
	ErrAgentSessionNotFound      = errors.New("agent session not found")
	ErrAgentSessionAlreadyExists = errors.New("agent session already exists")
	ErrInvalidSessionLabels      = errors.New("invalid session labels")
	ErrInvalidSessionReference   = errors.New("invalid session reference")
	ErrFailedCreateAgentSession  = errors.New("failed create agent session")
	ErrFailedListAgentSessions   = errors.New("failed list agent sessions")
	ErrFailedGetAgentSession     = errors.New("failed get agent session")
	ErrFailedUpdateAgentSession  = errors.New("failed update agent session")
	ErrFailedDeleteAgentSession  = errors.New("failed delete agent session")
)

// CreateAgentSessionRequest is the optional body for creating an agent session.
type CreateAgentSessionRequest struct {
	Title *string `json:"title" validate:"omitempty,max=255" label:"Title"`
}

// UpdateAgentSessionRequest is the body for updating an agent session.
type UpdateAgentSessionRequest struct {
	Title  *string           `json:"title" validate:"omitempty,max=255" label:"Title"`
	Labels map[string]string `json:"labels" validate:"omitempty" label:"Labels"`
}

// AgentSessionResponse is an agent session shaped for API responses.
type AgentSessionResponse struct {
	Id             db.Id             `json:"id"`
	ExternalId     string            `json:"externalId"`
	WorkspaceId    db.Id             `json:"workspaceId"`
	AgentId        db.Id             `json:"agentId"`
	Title          *string           `json:"title,omitempty"`
	Status         string            `json:"status"`
	Labels         map[string]string `json:"labels,omitempty"`
	Meta           json.RawMessage   `json:"meta,omitempty"`
	LastActivityAt string            `json:"lastActivityAt,omitempty"`
	CreatedAt      string            `json:"createdAt"`
	UpdatedAt      string            `json:"updatedAt"`
}

// ListAgentSessionsResponse is returned when listing agent sessions.
type ListAgentSessionsResponse struct {
	Sessions []*AgentSessionResponse
	Total    int64
}

// AgentSession is the module for agent conversation session CRUD.
type AgentSession struct {
	SessionRepository   db.AgentSessionRepository
	AgentRepository     db.AgentRepository
	WorkspaceRepository db.WorkspaceRepository
}

// NewAgentSession creates an agent session module with the given repositories.
func NewAgentSession(sessions db.AgentSessionRepository, agents db.AgentRepository, workspaces db.WorkspaceRepository) *AgentSession {
	return &AgentSession{
		SessionRepository:   sessions,
		AgentRepository:     agents,
		WorkspaceRepository: workspaces,
	}
}

// CreateAgentSession creates a session identified by query id and/or labels.
func (s *AgentSession) CreateAgentSession(workspaceId, agentId db.Id, externalId string, labels map[string]string, req *CreateAgentSessionRequest) (*AgentSessionResponse, error) {
	workspace, err := s.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := s.AgentRepository.GetById(agentId)
	if err != nil {
		return nil, err
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return nil, ErrAgentNotFound
	}

	if lo.IsEmpty(externalId) && len(labels) == 0 {
		return nil, ErrInvalidSessionReference
	}

	var labelsJSON *string
	if len(labels) > 0 {
		labelsRaw, err := json.Marshal(labels)
		if err != nil {
			return nil, ErrInvalidSessionLabels
		}
		labelsJSON = new(string(labelsRaw))

		existing, err := s.SessionRepository.GetByAgentIdAndExactLabels(agent.Id, string(labelsRaw))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedCreateAgentSession, err)
		}
		if existing != nil {
			return nil, ErrAgentSessionAlreadyExists
		}
	}

	if lo.IsEmpty(externalId) {
		token, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedCreateAgentSession, err)
		}
		externalId = token.String()
	} else {
		existing, err := s.SessionRepository.GetByAgentIdAndExternalId(agent.Id, externalId)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedCreateAgentSession, err)
		}
		if existing != nil {
			return nil, ErrAgentSessionAlreadyExists
		}
	}

	session := &db.AgentSession{
		ExternalId:  externalId,
		WorkspaceId: workspaceId,
		AgentId:     agent.Id,
		Labels:      labelsJSON,
		Status:      "active",
	}
	if req != nil {
		session.Title = req.Title
	}

	err = s.SessionRepository.Create(session)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateAgentSession, err)
	}

	var labelMap map[string]string
	if session.Labels != nil && lo.IsNotEmpty(*session.Labels) {
		_ = json.Unmarshal([]byte(*session.Labels), &labelMap)
	}

	response := &AgentSessionResponse{
		Id:          session.Id,
		ExternalId:  session.ExternalId,
		WorkspaceId: session.WorkspaceId,
		AgentId:     session.AgentId,
		Title:       session.Title,
		Status:      session.Status,
		Labels:      labelMap,
		Meta:        util.JSONRawFromString(session.Meta),
		LastActivityAt: lo.TernaryF(session.LastActivityAt != nil, func() string {
			return session.LastActivityAt.UTC().Format(time.RFC3339)
		}, func() string { return "" }),
		CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: session.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return response, nil
}

// GetAgentSession returns one session by internal id.
func (s *AgentSession) GetAgentSession(workspaceId, agentId, sessionId db.Id) (*AgentSessionResponse, error) {
	session, err := s.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		switch {
		case errors.Is(err, ErrWorkspaceNotFound),
			errors.Is(err, ErrAgentNotFound),
			errors.Is(err, ErrAgentSessionNotFound):
			return nil, err
		default:
			return nil, fmt.Errorf("%w: %v", ErrFailedGetAgentSession, err)
		}
	}

	var labelMap map[string]string
	if session.Labels != nil && lo.IsNotEmpty(*session.Labels) {
		_ = json.Unmarshal([]byte(*session.Labels), &labelMap)
	}

	response := &AgentSessionResponse{
		Id:          session.Id,
		ExternalId:  session.ExternalId,
		WorkspaceId: session.WorkspaceId,
		AgentId:     session.AgentId,
		Title:       session.Title,
		Status:      session.Status,
		Labels:      labelMap,
		Meta:        util.JSONRawFromString(session.Meta),
		LastActivityAt: lo.TernaryF(session.LastActivityAt != nil, func() string {
			return session.LastActivityAt.UTC().Format(time.RFC3339)
		}, func() string { return "" }),
		CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: session.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return response, nil
}

// LoadAgentSession loads a session after verifying it belongs to the workspace and agent.
func (s *AgentSession) LoadAgentSession(workspaceId, agentId, sessionId db.Id) (*db.AgentSession, error) {
	workspace, err := s.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := s.AgentRepository.GetById(agentId)
	if err != nil {
		return nil, err
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return nil, ErrAgentNotFound
	}

	session, err := s.SessionRepository.GetById(sessionId)
	if err != nil {
		return nil, err
	}
	if session == nil || session.AgentId != agent.Id || session.WorkspaceId != workspaceId {
		return nil, ErrAgentSessionNotFound
	}

	return session, nil
}

// ListAgentSessions returns paginated sessions, optionally filtered by id or labels.
func (s *AgentSession) ListAgentSessions(workspaceId, agentId db.Id, externalId string, labels map[string]string, limit, offset int) (*ListAgentSessionsResponse, error) {
	var (
		sessions []*db.AgentSession
		total    int64
	)

	workspace, err := s.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := s.AgentRepository.GetById(agentId)
	if err != nil {
		return nil, err
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return nil, ErrAgentNotFound
	}

	if lo.IsNotEmpty(externalId) {
		session, err := s.SessionRepository.GetByAgentIdAndExternalId(agent.Id, externalId)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedListAgentSessions, err)
		}
		if session != nil {
			sessions = []*db.AgentSession{session}
			total = 1
		}
	} else if len(labels) > 0 {
		labelsRaw, err := json.Marshal(labels)
		if err != nil {
			return nil, ErrInvalidSessionLabels
		}
		labelsJSON := string(labelsRaw)

		total, err = s.SessionRepository.CountByAgentIdAndLabels(agent.Id, labelsJSON)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedListAgentSessions, err)
		}

		sessions, err = s.SessionRepository.ListByAgentIdAndLabels(agent.Id, labelsJSON, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedListAgentSessions, err)
		}
	} else {
		total, err = s.SessionRepository.CountByAgentId(agent.Id)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedListAgentSessions, err)
		}

		sessions, err = s.SessionRepository.ListByAgentId(agent.Id, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedListAgentSessions, err)
		}
	}

	list := make([]*AgentSessionResponse, 0, len(sessions))
	for _, item := range sessions {
		var labelMap map[string]string
		if item.Labels != nil && lo.IsNotEmpty(*item.Labels) {
			_ = json.Unmarshal([]byte(*item.Labels), &labelMap)
		}

		response := &AgentSessionResponse{
			Id:          item.Id,
			ExternalId:  item.ExternalId,
			WorkspaceId: item.WorkspaceId,
			AgentId:     item.AgentId,
			Title:       item.Title,
			Status:      item.Status,
			Labels:      labelMap,
			Meta:        util.JSONRawFromString(item.Meta),
			LastActivityAt: lo.TernaryF(item.LastActivityAt != nil, func() string {
				return item.LastActivityAt.UTC().Format(time.RFC3339)
			}, func() string { return "" }),
			CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: item.UpdatedAt.UTC().Format(time.RFC3339),
		}

		list = append(list, response)
	}

	return &ListAgentSessionsResponse{Sessions: list, Total: total}, nil
}

// GetAgentSessionByLabels returns one session by external id or exact label match.
func (s *AgentSession) GetAgentSessionByLabels(workspaceId, agentId db.Id, externalId string, labels map[string]string) (*AgentSessionResponse, error) {
	workspace, err := s.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := s.AgentRepository.GetById(agentId)
	if err != nil {
		return nil, err
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return nil, ErrAgentNotFound
	}

	var session *db.AgentSession
	if lo.IsNotEmpty(externalId) {
		session, err = s.SessionRepository.GetByAgentIdAndExternalId(agent.Id, externalId)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedGetAgentSession, err)
		}
	} else if len(labels) > 0 {
		labelsRaw, err := json.Marshal(labels)
		if err != nil {
			return nil, ErrInvalidSessionLabels
		}

		session, err = s.SessionRepository.GetByAgentIdAndExactLabels(agent.Id, string(labelsRaw))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedGetAgentSession, err)
		}
	} else {
		return nil, ErrInvalidSessionReference
	}

	if session == nil {
		return nil, ErrAgentSessionNotFound
	}

	var labelMap map[string]string
	if session.Labels != nil && lo.IsNotEmpty(*session.Labels) {
		_ = json.Unmarshal([]byte(*session.Labels), &labelMap)
	}

	response := &AgentSessionResponse{
		Id:          session.Id,
		ExternalId:  session.ExternalId,
		WorkspaceId: session.WorkspaceId,
		AgentId:     session.AgentId,
		Title:       session.Title,
		Status:      session.Status,
		Labels:      labelMap,
		Meta:        util.JSONRawFromString(session.Meta),
		LastActivityAt: lo.TernaryF(session.LastActivityAt != nil, func() string {
			return session.LastActivityAt.UTC().Format(time.RFC3339)
		}, func() string { return "" }),
		CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: session.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return response, nil
}

// UpdateAgentSessionByLabels updates a session identified by external id or query labels.
func (s *AgentSession) UpdateAgentSessionByLabels(workspaceId, agentId db.Id, externalId string, labels map[string]string, req *UpdateAgentSessionRequest) (*AgentSessionResponse, error) {
	workspace, err := s.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := s.AgentRepository.GetById(agentId)
	if err != nil {
		return nil, err
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return nil, ErrAgentNotFound
	}

	var session *db.AgentSession
	if lo.IsNotEmpty(externalId) {
		session, err = s.SessionRepository.GetByAgentIdAndExternalId(agent.Id, externalId)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedUpdateAgentSession, err)
		}
	} else if len(labels) > 0 {
		labelsRaw, err := json.Marshal(labels)
		if err != nil {
			return nil, ErrInvalidSessionLabels
		}

		session, err = s.SessionRepository.GetByAgentIdAndExactLabels(agent.Id, string(labelsRaw))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedUpdateAgentSession, err)
		}
	} else {
		return nil, ErrInvalidSessionReference
	}

	if session == nil {
		return nil, ErrAgentSessionNotFound
	}

	if req.Title != nil {
		session.Title = req.Title
	}
	if req.Labels != nil {
		if len(req.Labels) == 0 {
			return nil, ErrInvalidSessionLabels
		}

		nextLabelsRaw, err := json.Marshal(req.Labels)
		if err != nil {
			return nil, ErrInvalidSessionLabels
		}
		nextLabelsJSON := string(nextLabelsRaw)

		if session.Labels == nil || *session.Labels != nextLabelsJSON {
			existing, err := s.SessionRepository.GetByAgentIdAndExactLabels(agent.Id, nextLabelsJSON)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrFailedUpdateAgentSession, err)
			}
			if existing != nil && existing.Id != session.Id {
				return nil, ErrAgentSessionAlreadyExists
			}
		}

		session.Labels = &nextLabelsJSON
	}

	err = s.SessionRepository.Update(session)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateAgentSession, err)
	}

	var labelMap map[string]string
	if session.Labels != nil && lo.IsNotEmpty(*session.Labels) {
		_ = json.Unmarshal([]byte(*session.Labels), &labelMap)
	}

	response := &AgentSessionResponse{
		Id:          session.Id,
		ExternalId:  session.ExternalId,
		WorkspaceId: session.WorkspaceId,
		AgentId:     session.AgentId,
		Title:       session.Title,
		Status:      session.Status,
		Labels:      labelMap,
		Meta:        util.JSONRawFromString(session.Meta),
		LastActivityAt: lo.TernaryF(session.LastActivityAt != nil, func() string {
			return session.LastActivityAt.UTC().Format(time.RFC3339)
		}, func() string { return "" }),
		CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: session.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return response, nil
}

// DeleteAgentSessionByLabels deletes a session identified by external id or query labels.
func (s *AgentSession) DeleteAgentSessionByLabels(workspaceId, agentId db.Id, externalId string, labels map[string]string) error {
	workspace, err := s.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	agent, err := s.AgentRepository.GetById(agentId)
	if err != nil {
		return err
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return ErrAgentNotFound
	}

	var session *db.AgentSession
	if lo.IsNotEmpty(externalId) {
		session, err = s.SessionRepository.GetByAgentIdAndExternalId(agent.Id, externalId)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedDeleteAgentSession, err)
		}
	} else if len(labels) > 0 {
		labelsRaw, err := json.Marshal(labels)
		if err != nil {
			return ErrInvalidSessionLabels
		}

		session, err = s.SessionRepository.GetByAgentIdAndExactLabels(agent.Id, string(labelsRaw))
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedDeleteAgentSession, err)
		}
	} else {
		return ErrInvalidSessionReference
	}

	if session == nil {
		return ErrAgentSessionNotFound
	}

	err = s.SessionRepository.Delete(session.Id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteAgentSession, err)
	}

	return s.EnqueueDeleteTask(workspaceId, agent.Id, session.Id)
}

// EnqueueDeleteTask enqueues an async task to delete a session resource.
func (s *AgentSession) EnqueueDeleteTask(workspaceId, agentId, sessionId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.delete",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"agentId":   agentId.String(),
			"sessionId": sessionId.String(),
		},
	})

	return err
}
