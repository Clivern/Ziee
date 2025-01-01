// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// AgentSession is a single row in the sessions table for an agent conversation.
type AgentSession struct {
	Id             Id
	ExternalId     string
	WorkspaceId    Id
	AgentId        Id
	Title          *string
	Status         string
	Labels         *string
	Meta           *string
	LastActivityAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// AgentSessionRepository is the interface for agent session CRUD.
type AgentSessionRepository interface {
	Create(session *AgentSession) error
	GetById(id Id) (*AgentSession, error)
	GetByAgentIdAndExternalId(agentId Id, externalId string) (*AgentSession, error)
	GetByAgentIdAndExactLabels(agentId Id, labels string) (*AgentSession, error)
	Update(session *AgentSession) error
	TouchLastActivityAt(sessionId Id, at time.Time) error
	Delete(id Id) error
	ListByAgentId(agentId Id, limit, offset int) ([]*AgentSession, error)
	ListByAgentIdAndLabels(agentId Id, labels string, limit, offset int) ([]*AgentSession, error)
	CountByAgentId(agentId Id) (int64, error)
	CountByAgentIdAndLabels(agentId Id, labels string) (int64, error)
}

type AgentSessionRepositoryPostgres struct {
	db *sql.DB
}

// AgentSessionMeta is a single row in the sessions_meta table.
type AgentSessionMeta struct {
	Id        Id
	SessionId Id
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AgentSessionMetaRepository is the interface for agent session meta CRUD.
type AgentSessionMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*AgentSessionMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListBySessionId(id Id) ([]*AgentSessionMeta, error)
	Upsert(id Id, key, value string) error
}

type AgentSessionMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewAgentSessionRepository returns an agent session repository.
func NewAgentSessionRepository(db *sql.DB) AgentSessionRepository {
	return &AgentSessionRepositoryPostgres{db: db}
}

// Create inserts an agent session row.
func (r *AgentSessionRepositoryPostgres) Create(session *AgentSession) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	session.Id = id

	err = r.db.QueryRow(
		`INSERT INTO sessions (
			id, external_id, workspace_id, agent_id, title, status, labels, meta
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`,
		session.Id.String(),
		session.ExternalId,
		session.WorkspaceId.String(),
		session.AgentId.String(),
		session.Title,
		session.Status,
		session.Labels,
		session.Meta,
	).Scan(
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	return err
}

// GetById returns an agent session by id.
func (r *AgentSessionRepositoryPostgres) GetById(id Id) (*AgentSession, error) {
	item := &AgentSession{}
	err := r.db.QueryRow(
		`SELECT
			id, external_id, workspace_id, agent_id, title, status, labels, meta, last_activity_at, created_at, updated_at
		FROM sessions
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.ExternalId,
		&item.WorkspaceId,
		&item.AgentId,
		&item.Title,
		&item.Status,
		&item.Labels,
		&item.Meta,
		&item.LastActivityAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByAgentIdAndExternalId returns an agent session by agent id and caller-provided external id.
func (r *AgentSessionRepositoryPostgres) GetByAgentIdAndExternalId(agentId Id, externalId string) (*AgentSession, error) {
	item := &AgentSession{}
	err := r.db.QueryRow(
		`SELECT
			id, external_id, workspace_id, agent_id, title, status, labels, meta, last_activity_at, created_at, updated_at
		FROM sessions
		WHERE agent_id = $1 AND external_id = $2`,
		agentId.String(),
		externalId,
	).Scan(
		&item.Id,
		&item.ExternalId,
		&item.WorkspaceId,
		&item.AgentId,
		&item.Title,
		&item.Status,
		&item.Labels,
		&item.Meta,
		&item.LastActivityAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByAgentIdAndExactLabels returns an agent session by agent id and exact labels match.
func (r *AgentSessionRepositoryPostgres) GetByAgentIdAndExactLabels(agentId Id, labels string) (*AgentSession, error) {
	item := &AgentSession{}
	err := r.db.QueryRow(
		`SELECT
			id, external_id, workspace_id, agent_id, title, status, labels, meta, last_activity_at, created_at, updated_at
		FROM sessions
		WHERE agent_id = $1 AND labels = $2::jsonb`,
		agentId.String(),
		labels,
	).Scan(
		&item.Id,
		&item.ExternalId,
		&item.WorkspaceId,
		&item.AgentId,
		&item.Title,
		&item.Status,
		&item.Labels,
		&item.Meta,
		&item.LastActivityAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates an existing agent session row.
func (r *AgentSessionRepositoryPostgres) Update(session *AgentSession) error {
	_, err := r.db.Exec(
		`UPDATE sessions
		SET
			workspace_id = $1,
			agent_id = $2,
			title = $3,
			status = $4,
			labels = $5,
			meta = $6,
			updated_at = $7
		WHERE id = $8`,
		session.WorkspaceId.String(),
		session.AgentId.String(),
		session.Title,
		session.Status,
		session.Labels,
		session.Meta,
		time.Now().UTC(),
		session.Id.String(),
	)
	return err
}

// TouchLastActivityAt sets last activity and updated timestamps for a session.
func (r *AgentSessionRepositoryPostgres) TouchLastActivityAt(sessionId Id, at time.Time) error {
	at = at.UTC()
	_, err := r.db.Exec(
		`UPDATE sessions
		SET last_activity_at = $2, updated_at = $2
		WHERE id = $1`,
		sessionId.String(),
		at,
	)
	return err
}

// Delete deletes an agent session row.
func (r *AgentSessionRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(
		`DELETE FROM sessions
		WHERE id = $1`,
		id.String(),
	)
	return err
}

// ListByAgentId lists agent session rows by agent id.
func (r *AgentSessionRepositoryPostgres) ListByAgentId(agentId Id, limit, offset int) ([]*AgentSession, error) {
	rows, err := r.db.Query(
		`SELECT
			id, external_id, workspace_id, agent_id, title, status, labels, meta, last_activity_at, created_at, updated_at
		FROM sessions
		WHERE agent_id = $1
		ORDER BY COALESCE(last_activity_at, created_at) DESC
		LIMIT $2 OFFSET $3`,
		agentId.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*AgentSession
	for rows.Next() {
		item := &AgentSession{}
		if err := rows.Scan(
			&item.Id,
			&item.ExternalId,
			&item.WorkspaceId,
			&item.AgentId,
			&item.Title,
			&item.Status,
			&item.Labels,
			&item.Meta,
			&item.LastActivityAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// ListByAgentIdAndLabels lists agent session rows by agent id and labels.
func (r *AgentSessionRepositoryPostgres) ListByAgentIdAndLabels(agentId Id, labels string, limit, offset int) ([]*AgentSession, error) {
	rows, err := r.db.Query(
		`SELECT
			id, external_id, workspace_id, agent_id, title, status, labels, meta, last_activity_at, created_at, updated_at
		FROM sessions
		WHERE agent_id = $1 AND labels @> $2::jsonb
		ORDER BY COALESCE(last_activity_at, created_at) DESC
		LIMIT $3 OFFSET $4`,
		agentId.String(),
		labels,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*AgentSession
	for rows.Next() {
		item := &AgentSession{}
		if err := rows.Scan(
			&item.Id,
			&item.ExternalId,
			&item.WorkspaceId,
			&item.AgentId,
			&item.Title,
			&item.Status,
			&item.Labels,
			&item.Meta,
			&item.LastActivityAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// CountByAgentId counts agent session rows by agent id.
func (r *AgentSessionRepositoryPostgres) CountByAgentId(agentId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM sessions
		WHERE agent_id = $1`,
		agentId.String(),
	).Scan(&count)
	return count, err
}

// CountByAgentIdAndLabels counts agent session rows by agent id and labels containment.
func (r *AgentSessionRepositoryPostgres) CountByAgentIdAndLabels(agentId Id, labels string) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM sessions
		WHERE agent_id = $1 AND labels @> $2::jsonb`,
		agentId.String(),
		labels,
	).Scan(&count)
	return count, err
}

// NewAgentSessionMetaRepository returns the repository for agent session meta.
func NewAgentSessionMetaRepository(db *sql.DB) AgentSessionMetaRepository {
	return &AgentSessionMetaRepositoryPostgres{db: db}
}

// Create inserts an agent session metadata row.
func (r *AgentSessionMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO sessions_meta (
			id, session_id, key, value
		)
		VALUES ($1, $2, $3, $4)`,
		metaId.String(), id.String(), key, value,
	)
	return err
}

// Get returns agent session metadata by key.
func (r *AgentSessionMetaRepositoryPostgres) Get(id Id, key string) (*AgentSessionMeta, error) {
	meta := &AgentSessionMeta{}
	err := r.db.QueryRow(
		`SELECT
			id, session_id, key, value, created_at, updated_at
		FROM sessions_meta
		WHERE session_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(
		&meta.Id,
		&meta.SessionId,
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

// Update updates an existing agent session metadata row.
func (r *AgentSessionMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE sessions_meta
		SET
			value = $1,
			updated_at = $2
		WHERE session_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)
	return err
}

// Delete deletes an agent session metadata row.
func (r *AgentSessionMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM sessions_meta
		WHERE session_id = $1 AND key = $2`,
		id.String(), key,
	)
	return err
}

// ListBySessionId lists agent session metadata rows by session id.
func (r *AgentSessionMetaRepositoryPostgres) ListBySessionId(id Id) ([]*AgentSessionMeta, error) {
	rows, err := r.db.Query(
		`SELECT
			id, session_id, key, value, created_at, updated_at
		FROM sessions_meta
		WHERE session_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*AgentSessionMeta
	for rows.Next() {
		meta := &AgentSessionMeta{}
		err := rows.Scan(
			&meta.Id,
			&meta.SessionId,
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

// Upsert creates or updates agent session metadata.
func (r *AgentSessionMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}
	return r.Update(id, key, value)
}
