// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// SessionMessage is a single row in the session_messages table.
type SessionMessage struct {
	Id         Id
	InternalId Id
	SessionId  Id
	Role       string
	Content    string
	Meta       *string
	CreatedAt  time.Time
}

// SessionMessageRepository is the interface for session message persistence.
type SessionMessageRepository interface {
	Create(message *SessionMessage) error
	GetById(id Id) (*SessionMessage, error)
	GetByInternalId(internalId Id) (*SessionMessage, error)
	Update(message *SessionMessage) error
	Delete(id Id) error
	DeleteBySessionId(sessionId Id) error
	CountBySessionId(sessionId Id) (int64, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
	ListBySessionId(sessionId Id, limit, offset int) ([]*SessionMessage, error)
}

type SessionMessageRepositoryPostgres struct {
	db *sql.DB
}

// NewSessionMessageRepository returns a session message repository.
func NewSessionMessageRepository(db *sql.DB) SessionMessageRepository {
	return &SessionMessageRepositoryPostgres{db: db}
}

// Create inserts a session message row.
func (r *SessionMessageRepositoryPostgres) Create(message *SessionMessage) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	message.Id = id

	internalId, err := NewId()
	if err != nil {
		return err
	}
	message.InternalId = internalId

	return r.db.QueryRow(
		`INSERT INTO session_messages (
			id, internal_id, session_id, role, content, meta
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`,
		message.Id.String(),
		message.InternalId.String(),
		message.SessionId.String(),
		message.Role,
		message.Content,
		message.Meta,
	).Scan(&message.CreatedAt)
}

// GetById returns a session message by id.
func (r *SessionMessageRepositoryPostgres) GetById(id Id) (*SessionMessage, error) {
	item := &SessionMessage{}
	err := r.db.QueryRow(
		`SELECT
			id, internal_id, session_id, role, content, meta, created_at
		FROM session_messages
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.InternalId,
		&item.SessionId,
		&item.Role,
		&item.Content,
		&item.Meta,
		&item.CreatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByInternalId returns a session message by internal id.
func (r *SessionMessageRepositoryPostgres) GetByInternalId(internalId Id) (*SessionMessage, error) {
	item := &SessionMessage{}
	err := r.db.QueryRow(
		`SELECT
			id, internal_id, session_id, role, content, meta, created_at
		FROM session_messages
		WHERE internal_id = $1`,
		internalId.String(),
	).Scan(
		&item.Id,
		&item.InternalId,
		&item.SessionId,
		&item.Role,
		&item.Content,
		&item.Meta,
		&item.CreatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates a session message row.
func (r *SessionMessageRepositoryPostgres) Update(message *SessionMessage) error {
	_, err := r.db.Exec(
		`UPDATE session_messages
		SET role = $2, content = $3, meta = $4
		WHERE id = $1`,
		message.Id.String(),
		message.Role,
		message.Content,
		message.Meta,
	)
	return err
}

// Delete deletes a session message row.
func (r *SessionMessageRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(
		`DELETE FROM session_messages
		WHERE id = $1`,
		id.String(),
	)
	return err
}

// DeleteBySessionId deletes session message rows by session id.
func (r *SessionMessageRepositoryPostgres) DeleteBySessionId(sessionId Id) error {
	_, err := r.db.Exec(
		`DELETE FROM session_messages
		WHERE session_id = $1`,
		sessionId.String(),
	)
	return err
}

// CountBySessionId counts session message rows by session id.
func (r *SessionMessageRepositoryPostgres) CountBySessionId(sessionId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM session_messages
		WHERE session_id = $1`,
		sessionId.String(),
	).Scan(&count)
	return count, err
}

// CountByWorkspaceId counts session message rows across a workspace.
func (r *SessionMessageRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM session_messages sm
		INNER JOIN sessions s ON s.id = sm.session_id
		WHERE s.workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// ListBySessionId lists session message rows by session id.
func (r *SessionMessageRepositoryPostgres) ListBySessionId(sessionId Id, limit, offset int) ([]*SessionMessage, error) {
	rows, err := r.db.Query(
		`SELECT
			id, internal_id, session_id, role, content, meta, created_at
		FROM session_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`,
		sessionId.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*SessionMessage
	for rows.Next() {
		item := &SessionMessage{}
		err := rows.Scan(
			&item.Id,
			&item.InternalId,
			&item.SessionId,
			&item.Role,
			&item.Content,
			&item.Meta,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}
