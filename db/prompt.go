// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Prompt is a prompt row with its version count.
type Prompt struct {
	Id           Id
	Name         string
	Handle       string
	Description  string
	VersionCount int64
}

// PromptRepository is the interface for prompt parent-row CRUD.
type PromptRepository interface {
	Create(workspaceId Id, name, handle, description string) (Id, error)
	ExistsByHandle(workspaceId Id, handle string) (bool, error)
	GetById(workspaceId Id, id Id) (*Prompt, error)
	GetByHandle(workspaceId Id, handle string) (*Prompt, error)
	Delete(workspaceId Id, promptId Id) error
	ListByWorkspace(workspaceId Id, limit, offset int) ([]*Prompt, error)
	Count() (int64, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
}

type PromptRepositoryPostgres struct {
	db *sql.DB
}

// NewPromptRepository returns a prompt repository.
func NewPromptRepository(db *sql.DB) PromptRepository {
	return &PromptRepositoryPostgres{db: db}
}

// Create inserts a new prompt row.
func (r *PromptRepositoryPostgres) Create(workspaceId Id, name, handle, description string) (Id, error) {
	id, err := NewId()
	if err != nil {
		return "", err
	}
	_, err = r.db.Exec(
		`INSERT INTO prompts (id, workspace_id, name, description, handle) VALUES ($1, $2, $3, $4, $5)`,
		id.String(),
		workspaceId.String(),
		name,
		description,
		handle,
	)
	return id, err
}

// ExistsByHandle reports whether a handle is already used in the workspace.
func (r *PromptRepositoryPostgres) ExistsByHandle(workspaceId Id, handle string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM prompts WHERE workspace_id = $1 AND handle = $2)`,
		workspaceId.String(),
		handle,
	).Scan(&exists)
	return exists, err
}

// GetById returns a prompt with its version count.
func (r *PromptRepositoryPostgres) GetById(workspaceId, id Id) (*Prompt, error) {
	item := &Prompt{}
	err := r.db.QueryRow(`
		SELECT
			p.id,
			p.name,
			p.handle,
			p.description,
			(SELECT COUNT(*) FROM prompt_versions v WHERE v.prompt_id = p.id)
		FROM prompts p
		WHERE p.workspace_id = $1 AND p.id = $2`,
		workspaceId.String(),
		id.String(),
	).Scan(
		&item.Id,
		&item.Name,
		&item.Handle,
		&item.Description,
		&item.VersionCount,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByHandle returns a prompt by workspace and handle.
func (r *PromptRepositoryPostgres) GetByHandle(workspaceId Id, handle string) (*Prompt, error) {
	item := &Prompt{}
	err := r.db.QueryRow(`
		SELECT
			p.id,
			p.name,
			p.handle,
			p.description,
			(SELECT COUNT(*) FROM prompt_versions v WHERE v.prompt_id = p.id)
		FROM prompts p
		WHERE p.workspace_id = $1 AND p.handle = $2`,
		workspaceId.String(),
		handle,
	).Scan(
		&item.Id,
		&item.Name,
		&item.Handle,
		&item.Description,
		&item.VersionCount,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Delete deletes a prompt and all of its versions.
func (r *PromptRepositoryPostgres) Delete(workspaceId Id, promptId Id) error {
	result, err := r.db.Exec(
		`DELETE FROM prompts WHERE id = $1 AND workspace_id = $2`,
		promptId.String(),
		workspaceId.String(),
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ListByWorkspace lists prompts for a workspace with version counts.
func (r *PromptRepositoryPostgres) ListByWorkspace(workspaceId Id, limit, offset int) ([]*Prompt, error) {
	rows, err := r.db.Query(`
		SELECT
			p.id,
			p.name,
			p.handle,
			p.description,
			(SELECT COUNT(*) FROM prompt_versions v WHERE v.prompt_id = p.id)
		FROM prompts p
		WHERE p.workspace_id = $1
		ORDER BY p.handle ASC
		LIMIT $2 OFFSET $3`, workspaceId.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Prompt
	for rows.Next() {
		item := &Prompt{}
		if err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Handle,
			&item.Description,
			&item.VersionCount,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// CountByWorkspaceId counts prompt rows by workspace id.
func (r *PromptRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM prompts
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// Count returns the total number of prompt rows.
func (r *PromptRepositoryPostgres) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM prompts`).Scan(&count)
	return count, err
}

// PromptMeta is a single row in the prompts_meta table.
type PromptMeta struct {
	Id        Id
	PromptId  Id
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PromptMetaRepository is the interface for prompt metadata CRUD.
type PromptMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*PromptMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByPromptId(id Id) ([]*PromptMeta, error)
	Upsert(id Id, key, value string) error
}

type PromptMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewPromptMetaRepository returns the repository for prompt metadata.
func NewPromptMetaRepository(db *sql.DB) PromptMetaRepository {
	return &PromptMetaRepositoryPostgres{db: db}
}

// Create inserts a prompt metadata row.
func (r *PromptMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO prompts_meta (id, prompt_id, key, value)
		VALUES ($1, $2, $3, $4)`,
		metaId.String(), id.String(), key, value,
	)
	return err
}

// Get returns prompt metadata by key.
func (r *PromptMetaRepositoryPostgres) Get(id Id, key string) (*PromptMeta, error) {
	meta := &PromptMeta{}
	err := r.db.QueryRow(
		`SELECT id, prompt_id, key, value, created_at, updated_at
		FROM prompts_meta
		WHERE prompt_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(
		&meta.Id,
		&meta.PromptId,
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

// Update updates an existing prompt metadata row.
func (r *PromptMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE prompts_meta
		SET value = $1, updated_at = $2
		WHERE prompt_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)
	return err
}

// Delete deletes a prompt metadata row.
func (r *PromptMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM prompts_meta
		WHERE prompt_id = $1 AND key = $2`,
		id.String(), key,
	)
	return err
}

// ListByPromptId lists prompt metadata rows by prompt id.
func (r *PromptMetaRepositoryPostgres) ListByPromptId(id Id) ([]*PromptMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, prompt_id, key, value, created_at, updated_at
		FROM prompts_meta
		WHERE prompt_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*PromptMeta
	for rows.Next() {
		meta := &PromptMeta{}
		err := rows.Scan(
			&meta.Id,
			&meta.PromptId,
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

// Upsert creates or updates prompt metadata.
func (r *PromptMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}
	return r.Update(id, key, value)
}
