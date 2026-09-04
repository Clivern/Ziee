// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// WorkspaceAccessKey is the DB row for a workspace-scoped access key.
type WorkspaceAccessKey struct {
	Id          Id
	WorkspaceId Id
	Name        string
	Key         string
	ExpiresAt   *time.Time
	Meta        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WorkspaceAccessKeyRepository is the interface for workspace access key CRUD.
type WorkspaceAccessKeyRepository interface {
	Create(accessKey *WorkspaceAccessKey) error
	GetById(id Id) (*WorkspaceAccessKey, error)
	GetByKey(key string) (*WorkspaceAccessKey, error)
	ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*WorkspaceAccessKey, error)
	Delete(id Id) error
	DeleteExpired() (int64, error)
	Count() (int64, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
}

type WorkspaceAccessKeyRepositoryPostgres struct {
	db *sql.DB
}

// NewWorkspaceAccessKeyRepository returns the repository for workspace access keys.
func NewWorkspaceAccessKeyRepository(db *sql.DB) WorkspaceAccessKeyRepository {
	return &WorkspaceAccessKeyRepositoryPostgres{db: db}
}

// Create inserts a workspace access key row.
func (r *WorkspaceAccessKeyRepositoryPostgres) Create(accessKey *WorkspaceAccessKey) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	accessKey.Id = id

	return r.db.QueryRow(
		`INSERT INTO workspace_access_keys
		(id, workspace_id, name, token, expires_at, meta)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`,
		accessKey.Id.String(),
		accessKey.WorkspaceId.String(),
		accessKey.Name,
		accessKey.Key,
		accessKey.ExpiresAt,
		accessKey.Meta,
	).Scan(&accessKey.CreatedAt, &accessKey.UpdatedAt)
}

// GetById returns a workspace access key by id.
func (r *WorkspaceAccessKeyRepositoryPostgres) GetById(id Id) (*WorkspaceAccessKey, error) {
	item := &WorkspaceAccessKey{}
	err := r.db.QueryRow(
		`SELECT id, workspace_id, name, token, expires_at, meta, created_at, updated_at
		FROM workspace_access_keys
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.Name,
		&item.Key,
		&item.ExpiresAt,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByKey returns a workspace access key by key.
func (r *WorkspaceAccessKeyRepositoryPostgres) GetByKey(key string) (*WorkspaceAccessKey, error) {
	item := &WorkspaceAccessKey{}
	err := r.db.QueryRow(
		`SELECT id, workspace_id, name, token, expires_at, meta, created_at, updated_at
		FROM workspace_access_keys
		WHERE token = $1 AND (expires_at IS NULL OR expires_at > $2)`,
		key,
		time.Now().UTC(),
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.Name,
		&item.Key,
		&item.ExpiresAt,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// ListByWorkspaceId lists workspace access key rows by workspace id.
func (r *WorkspaceAccessKeyRepositoryPostgres) ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*WorkspaceAccessKey, error) {
	rows, err := r.db.Query(
		`SELECT id, workspace_id, name, token, expires_at, meta, created_at, updated_at
		FROM workspace_access_keys
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

	var list []*WorkspaceAccessKey
	for rows.Next() {
		item := &WorkspaceAccessKey{}
		if err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.Name,
			&item.Key,
			&item.ExpiresAt,
			&item.Meta,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// Delete deletes a workspace access key row.
func (r *WorkspaceAccessKeyRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(`DELETE FROM workspace_access_keys WHERE id = $1`, id.String())
	return err
}

// DeleteExpired deletes expired workspace access keys.
func (r *WorkspaceAccessKeyRepositoryPostgres) DeleteExpired() (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM workspace_access_keys
		WHERE expires_at IS NOT NULL AND expires_at < $1`,
		time.Now().UTC(),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Count returns the total number of workspace access key rows.
func (r *WorkspaceAccessKeyRepositoryPostgres) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM workspace_access_keys
		WHERE expires_at IS NULL OR expires_at > $1`,
		time.Now().UTC(),
	).Scan(&count)
	return count, err
}

// CountByWorkspaceId counts workspace access key rows by workspace id.
func (r *WorkspaceAccessKeyRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM workspace_access_keys
		WHERE workspace_id = $1 AND (expires_at IS NULL OR expires_at > $2)`,
		workspaceId.String(),
		time.Now().UTC(),
	).Scan(&count)
	return count, err
}
