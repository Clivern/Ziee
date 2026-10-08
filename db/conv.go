// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const (
	ConversationTargetIssue       = "issue"
	ConversationTargetPullRequest = "pull_request"

	ConversationRoleUser      = "user"
	ConversationRoleAssistant = "assistant"
)

// RepositoryConversation is one message in the assistant history of an issue or pull request.
type RepositoryConversation struct {
	Id           Id
	RepositoryId Id
	TargetType   string
	TargetNumber int
	Role         string
	Author       *string
	CommentId    *int64
	Body         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RepositoryConversationRepository is the interface for repository conversation persistence.
type RepositoryConversationRepository interface {
	Create(item *RepositoryConversation) error
	GetById(id Id) (*RepositoryConversation, error)
	ListByTarget(repositoryId Id, targetType string, targetNumber int) ([]*RepositoryConversation, error)
	DeleteByTarget(repositoryId Id, targetType string, targetNumber int) error
}

type RepositoryConversationRepositoryPostgres struct {
	db *sql.DB
}

// NewRepositoryConversationRepository returns the repository for repository_conversations.
func NewRepositoryConversationRepository(db *sql.DB) RepositoryConversationRepository {
	return &RepositoryConversationRepositoryPostgres{db: db}
}

// Create inserts a repository conversation message.
func (r *RepositoryConversationRepositoryPostgres) Create(item *RepositoryConversation) error {
	id, err := NewId()
	if err != nil {
		return err
	}

	item.Id = id

	return r.db.QueryRow(
		`INSERT INTO repository_conversations (
			id, repository_id, target_type, target_number, role, author, comment_id, body
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`,
		item.Id.String(),
		item.RepositoryId.String(),
		item.TargetType,
		item.TargetNumber,
		item.Role,
		item.Author,
		item.CommentId,
		item.Body,
	).Scan(&item.CreatedAt, &item.UpdatedAt)
}

// GetById returns a repository conversation message by id.
func (r *RepositoryConversationRepositoryPostgres) GetById(id Id) (*RepositoryConversation, error) {
	item := &RepositoryConversation{}
	err := r.db.QueryRow(
		`SELECT
			id, repository_id, target_type, target_number, role,
			author, comment_id, body, created_at, updated_at
		FROM repository_conversations
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.RepositoryId,
		&item.TargetType,
		&item.TargetNumber,
		&item.Role,
		&item.Author,
		&item.CommentId,
		&item.Body,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// ListByTarget lists the message history of one issue or pull request, oldest first.
func (r *RepositoryConversationRepositoryPostgres) ListByTarget(repositoryId Id, targetType string, targetNumber int) ([]*RepositoryConversation, error) {
	rows, err := r.db.Query(
		`SELECT
			id, repository_id, target_type, target_number, role,
			author, comment_id, body, created_at, updated_at
		FROM repository_conversations
		WHERE repository_id = $1 AND target_type = $2 AND target_number = $3
		ORDER BY created_at ASC`,
		repositoryId.String(),
		targetType,
		targetNumber,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*RepositoryConversation
	for rows.Next() {
		item := &RepositoryConversation{}
		err = rows.Scan(
			&item.Id,
			&item.RepositoryId,
			&item.TargetType,
			&item.TargetNumber,
			&item.Role,
			&item.Author,
			&item.CommentId,
			&item.Body,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

// DeleteByTarget deletes the message history of one issue or pull request.
func (r *RepositoryConversationRepositoryPostgres) DeleteByTarget(repositoryId Id, targetType string, targetNumber int) error {
	_, err := r.db.Exec(
		`DELETE FROM repository_conversations
		WHERE repository_id = $1 AND target_type = $2 AND target_number = $3`,
		repositoryId.String(),
		targetType,
		targetNumber,
	)

	return err
}

// RepositoryConversationMeta is a single row in the repository_conversations_meta table.
type RepositoryConversationMeta struct {
	Id             Id
	ConversationId Id
	Key            string
	Value          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// RepositoryConversationMetaRepository is the interface for conversation metadata CRUD.
type RepositoryConversationMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*RepositoryConversationMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByConversationId(id Id) ([]*RepositoryConversationMeta, error)
	Upsert(id Id, key, value string) error
}

type RepositoryConversationMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewRepositoryConversationMetaRepository returns the repository for conversation metadata.
func NewRepositoryConversationMetaRepository(db *sql.DB) RepositoryConversationMetaRepository {
	return &RepositoryConversationMetaRepositoryPostgres{db: db}
}

// Create inserts a conversation metadata row.
func (r *RepositoryConversationMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`INSERT INTO repository_conversations_meta (id, conversation_id, key, value)
		VALUES ($1, $2, $3, to_jsonb($4::text))`,
		metaId.String(), id.String(), key, value,
	)

	return err
}

// Get returns conversation metadata by key.
func (r *RepositoryConversationMetaRepositoryPostgres) Get(id Id, key string) (*RepositoryConversationMeta, error) {
	meta := &RepositoryConversationMeta{}
	err := r.db.QueryRow(
		`SELECT id, conversation_id, key, value #>> '{}', created_at, updated_at
		FROM repository_conversations_meta
		WHERE conversation_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(&meta.Id, &meta.ConversationId, &meta.Key, &meta.Value, &meta.CreatedAt, &meta.UpdatedAt)
	if isNotFound(err) {
		return nil, nil
	}

	return meta, err
}

// Update updates an existing conversation metadata row.
func (r *RepositoryConversationMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE repository_conversations_meta
		SET value = to_jsonb($1::text), updated_at = $2
		WHERE conversation_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)

	return err
}

// Delete deletes a conversation metadata row.
func (r *RepositoryConversationMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM repository_conversations_meta
		WHERE conversation_id = $1 AND key = $2`,
		id.String(), key,
	)

	return err
}

// ListByConversationId lists conversation metadata rows by conversation id.
func (r *RepositoryConversationMetaRepositoryPostgres) ListByConversationId(id Id) ([]*RepositoryConversationMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, conversation_id, key, value #>> '{}', created_at, updated_at
		FROM repository_conversations_meta
		WHERE conversation_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list []*RepositoryConversationMeta
	for rows.Next() {
		meta := &RepositoryConversationMeta{}
		err := rows.Scan(&meta.Id, &meta.ConversationId, &meta.Key, &meta.Value, &meta.CreatedAt, &meta.UpdatedAt)
		if err != nil {
			return nil, err
		}

		list = append(list, meta)
	}

	return list, rows.Err()
}

// Upsert creates or updates conversation metadata.
func (r *RepositoryConversationMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}

	return r.Update(id, key, value)
}
