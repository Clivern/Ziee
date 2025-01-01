// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"

	"github.com/samber/lo"
)

const (
	AgentKindUnmanaged = "unmanaged"
)

// Agent is a single row in the agents table.
type Agent struct {
	Id          Id
	WorkspaceId Id
	Name        string
	Kind        string
	PromptId    *Id
	KbLabels    *string
	Handle      string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AgentRepository is the interface for agent CRUD.
type AgentRepository interface {
	Create(agent *Agent) error
	GetById(id Id) (*Agent, error)
	GetByHandle(handle string) (*Agent, error)
	Update(agent *Agent) error
	Delete(id Id) error
	ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*Agent, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
}

type AgentRepositoryPostgres struct {
	db *sql.DB
}

// AgentMeta is a single row in the agents_meta table.
type AgentMeta struct {
	Id        Id
	AgentId   Id
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AgentMetaRepository is the interface for agent meta CRUD.
type AgentMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*AgentMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByAgentId(id Id) ([]*AgentMeta, error)
	Upsert(id Id, key, value string) error
}

type AgentMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewAgentRepository returns an agent repository.
func NewAgentRepository(db *sql.DB) AgentRepository {
	return &AgentRepositoryPostgres{db: db}
}

// Create inserts an agent row.
func (r *AgentRepositoryPostgres) Create(agent *Agent) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	agent.Id = id

	// Set default kind if not provided
	// TODO: Remove this once we support managed agents
	if lo.IsEmpty(agent.Kind) {
		agent.Kind = AgentKindUnmanaged
	}

	err = r.db.QueryRow(
		`INSERT INTO agents (
			id, workspace_id, name, kind, prompt_id, kb_labels, handle, description
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING status, created_at, updated_at`,
		agent.Id.String(),
		agent.WorkspaceId.String(),
		agent.Name,
		agent.Kind,
		lo.FromPtr(agent.PromptId),
		agent.KbLabels,
		agent.Handle,
		agent.Description,
	).Scan(
		&agent.Status,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)
	return err
}

// GetById returns an agent by id.
func (r *AgentRepositoryPostgres) GetById(id Id) (*Agent, error) {
	item := &Agent{}
	err := r.db.QueryRow(
		`SELECT
			id, workspace_id, name, kind, prompt_id, kb_labels, handle, description, status, created_at, updated_at
		FROM agents
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.Name,
		&item.Kind,
		&item.PromptId,
		&item.KbLabels,
		&item.Handle,
		&item.Description,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByHandle returns an agent by handle.
func (r *AgentRepositoryPostgres) GetByHandle(handle string) (*Agent, error) {
	item := &Agent{}
	err := r.db.QueryRow(
		`SELECT
			id, workspace_id, name, kind, prompt_id, kb_labels, handle, description, status, created_at, updated_at
		FROM agents
		WHERE handle = $1`,
		handle,
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.Name,
		&item.Kind,
		&item.PromptId,
		&item.KbLabels,
		&item.Handle,
		&item.Description,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates an existing agent row.
func (r *AgentRepositoryPostgres) Update(agent *Agent) error {
	_, err := r.db.Exec(
		`UPDATE agents
		SET
			workspace_id = $1,
			name = $2,
			kind = $3,
			prompt_id = $4,
			kb_labels = $5,
			description = $6,
			updated_at = $7
		WHERE id = $8`,
		agent.WorkspaceId.String(),
		agent.Name,
		agent.Kind,
		lo.FromPtr(agent.PromptId),
		agent.KbLabels,
		agent.Description,
		time.Now().UTC(),
		agent.Id.String(),
	)
	return err
}

// Delete deletes an agent row.
func (r *AgentRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(
		`DELETE FROM agents
		WHERE id = $1`,
		id.String(),
	)
	return err
}

// ListByWorkspaceId lists agent rows by workspace id.
func (r *AgentRepositoryPostgres) ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*Agent, error) {
	rows, err := r.db.Query(
		`SELECT
			id, workspace_id, name, kind, prompt_id, kb_labels, handle, description, status, created_at, updated_at
		FROM agents
		WHERE workspace_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		workspaceId.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Agent
	for rows.Next() {
		item := &Agent{}
		if err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.Name,
			&item.Kind,
			&item.PromptId,
			&item.KbLabels,
			&item.Handle,
			&item.Description,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// CountByWorkspaceId counts agent rows by workspace id.
func (r *AgentRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM agents
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// NewAgentMetaRepository returns the repository for agent meta.
func NewAgentMetaRepository(db *sql.DB) AgentMetaRepository {
	return &AgentMetaRepositoryPostgres{db: db}
}

// Create inserts an agent metadata row.
func (r *AgentMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO agents_meta (
			id, agent_id, key, value
		)
		VALUES ($1, $2, $3, $4)`,
		metaId.String(), id.String(), key, value,
	)
	return err
}

// Get returns agent metadata by key.
func (r *AgentMetaRepositoryPostgres) Get(id Id, key string) (*AgentMeta, error) {
	meta := &AgentMeta{}
	err := r.db.QueryRow(
		`SELECT
			id, agent_id, key, value, created_at, updated_at
		FROM agents_meta
		WHERE agent_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(
		&meta.Id,
		&meta.AgentId,
		&meta.Key,
		&meta.Value,
		&meta.CreatedAt,
		&meta.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return meta, err
}

// Update updates an existing agent metadata row.
func (r *AgentMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE agents_meta
		SET
			value = $1,
			updated_at = $2
		WHERE agent_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)
	return err
}

// Delete deletes an agent metadata row.
func (r *AgentMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM agents_meta
		WHERE agent_id = $1 AND key = $2`,
		id.String(), key,
	)
	return err
}

// ListByAgentId lists agent metadata rows by agent id.
func (r *AgentMetaRepositoryPostgres) ListByAgentId(id Id) ([]*AgentMeta, error) {
	rows, err := r.db.Query(
		`SELECT
			id, agent_id, key, value, created_at, updated_at
		FROM agents_meta
		WHERE agent_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*AgentMeta
	for rows.Next() {
		meta := &AgentMeta{}
		err := rows.Scan(
			&meta.Id,
			&meta.AgentId,
			&meta.Key,
			&meta.Value,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, meta)
	}
	return list, rows.Err()
}

// Upsert creates or updates agent metadata.
func (r *AgentMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}
	return r.Update(id, key, value)
}
