// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// SessionMemory is a single row in the session_memories table.
type SessionMemory struct {
	Id         Id
	InternalId Id
	SessionId  Id
	Kind       string
	Content    string
	Meta       *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// SessionMemoryRepository is the interface for session memory CRUD.
type SessionMemoryRepository interface {
	Create(memory *SessionMemory) error
	GetById(id Id) (*SessionMemory, error)
	GetByInternalId(internalId Id) (*SessionMemory, error)
	Update(memory *SessionMemory) error
	Delete(id Id) error
	DeleteBySessionId(sessionId Id) error
	CountBySessionId(sessionId Id) (int64, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
	ListBySessionId(sessionId Id, limit, offset int) ([]*SessionMemory, error)
}

type SessionMemoryRepositoryPostgres struct {
	db *sql.DB
}

// NewSessionMemoryRepository returns a session memory repository.
func NewSessionMemoryRepository(db *sql.DB) SessionMemoryRepository {
	return &SessionMemoryRepositoryPostgres{db: db}
}

// Create inserts a session memory row.
func (r *SessionMemoryRepositoryPostgres) Create(memory *SessionMemory) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	memory.Id = id

	internalId, err := NewId()
	if err != nil {
		return err
	}
	memory.InternalId = internalId

	return r.db.QueryRow(
		`INSERT INTO session_memories (
			id, internal_id, session_id, kind, content, meta
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`,
		memory.Id.String(),
		memory.InternalId.String(),
		memory.SessionId.String(),
		memory.Kind,
		memory.Content,
		memory.Meta,
	).Scan(
		&memory.CreatedAt,
		&memory.UpdatedAt,
	)
}

// GetById returns a session memory by id.
func (r *SessionMemoryRepositoryPostgres) GetById(id Id) (*SessionMemory, error) {
	item := &SessionMemory{}
	err := r.db.QueryRow(
		`SELECT
			id, internal_id, session_id, kind, content, meta, created_at, updated_at
		FROM session_memories
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.InternalId,
		&item.SessionId,
		&item.Kind,
		&item.Content,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByInternalId returns a session memory by internal id.
func (r *SessionMemoryRepositoryPostgres) GetByInternalId(internalId Id) (*SessionMemory, error) {
	item := &SessionMemory{}
	err := r.db.QueryRow(
		`SELECT
			id, internal_id, session_id, kind, content, meta, created_at, updated_at
		FROM session_memories
		WHERE internal_id = $1`,
		internalId.String(),
	).Scan(
		&item.Id,
		&item.InternalId,
		&item.SessionId,
		&item.Kind,
		&item.Content,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates an existing session memory row.
func (r *SessionMemoryRepositoryPostgres) Update(memory *SessionMemory) error {
	return r.db.QueryRow(
		`UPDATE session_memories
		SET
			kind = $1,
			content = $2,
			meta = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING updated_at`,
		memory.Kind,
		memory.Content,
		memory.Meta,
		time.Now().UTC(),
		memory.Id.String(),
	).Scan(&memory.UpdatedAt)
}

// Delete deletes a session memory row.
func (r *SessionMemoryRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(
		`DELETE FROM session_memories
		WHERE id = $1`,
		id.String(),
	)
	return err
}

// DeleteBySessionId deletes session memory rows by session id.
func (r *SessionMemoryRepositoryPostgres) DeleteBySessionId(sessionId Id) error {
	_, err := r.db.Exec(
		`DELETE FROM session_memories
		WHERE session_id = $1`,
		sessionId.String(),
	)
	return err
}

// CountBySessionId counts session memory rows by session id.
func (r *SessionMemoryRepositoryPostgres) CountBySessionId(sessionId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM session_memories
		WHERE session_id = $1`,
		sessionId.String(),
	).Scan(&count)
	return count, err
}

// CountByWorkspaceId counts session memory rows across a workspace.
func (r *SessionMemoryRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM session_memories sm
		INNER JOIN sessions s ON s.id = sm.session_id
		WHERE s.workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// ListBySessionId lists session memory rows by session id.
func (r *SessionMemoryRepositoryPostgres) ListBySessionId(sessionId Id, limit, offset int) ([]*SessionMemory, error) {
	rows, err := r.db.Query(
		`SELECT
			id, internal_id, session_id, kind, content, meta, created_at, updated_at
		FROM session_memories
		WHERE session_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		sessionId.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*SessionMemory
	for rows.Next() {
		item := &SessionMemory{}
		err := rows.Scan(
			&item.Id,
			&item.InternalId,
			&item.SessionId,
			&item.Kind,
			&item.Content,
			&item.Meta,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}
