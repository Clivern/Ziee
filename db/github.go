// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const RepositoryMetaConfigPath = "_github_config_path"

// Repository is a GitHub repository that installed the Ziee GitHub App.
type Repository struct {
	Id             Id
	WorkspaceId    Id
	InstallationId int64
	GitHubId       int64
	NodeId         string
	Owner          string
	Name           string
	FullName       string
	Private        bool
	Meta           *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// RepositoryRepository is the interface for workspace GitHub repo CRUD.
type RepositoryRepository interface {
	Create(repo *Repository) error
	Upsert(repo *Repository) error
	GetById(id Id) (*Repository, error)
	GetByGitHubId(githubId int64) (*Repository, error)
	Update(repo *Repository) error
	Delete(id Id) error
	DeleteByGitHubId(githubId int64) error
	DeleteByInstallationId(installationId int64) error
	ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*Repository, error)
	ListByInstallationId(installationId int64) ([]*Repository, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
}

type RepositoryRepositoryPostgres struct {
	db *sql.DB
}

// NewRepositoryRepository returns the repository for workspace GitHub repos.
func NewRepositoryRepository(db *sql.DB) RepositoryRepository {
	return &RepositoryRepositoryPostgres{db: db}
}

// Create inserts a workspace GitHub repo row.
func (r *RepositoryRepositoryPostgres) Create(repo *Repository) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	repo.Id = id

	return r.db.QueryRow(
		`INSERT INTO repositories
		(id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at`,
		repo.Id.String(),
		repo.WorkspaceId.String(),
		repo.InstallationId,
		repo.GitHubId,
		repo.NodeId,
		repo.Owner,
		repo.Name,
		repo.FullName,
		repo.Private,
		repo.Meta,
	).Scan(&repo.CreatedAt, &repo.UpdatedAt)
}

// Upsert inserts or updates a workspace GitHub repo by GitHub repository id.
func (r *RepositoryRepositoryPostgres) Upsert(repo *Repository) error {
	id, err := NewId()
	if err != nil {
		return err
	}

	return r.db.QueryRow(
		`INSERT INTO repositories
		(id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (github_id) DO UPDATE SET
			workspace_id = EXCLUDED.workspace_id,
			installation_id = EXCLUDED.installation_id,
			node_id = EXCLUDED.node_id,
			owner = EXCLUDED.owner,
			name = EXCLUDED.name,
			full_name = EXCLUDED.full_name,
			private = EXCLUDED.private,
			meta = EXCLUDED.meta,
			updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
		RETURNING id, created_at, updated_at`,
		id.String(),
		repo.WorkspaceId.String(),
		repo.InstallationId,
		repo.GitHubId,
		repo.NodeId,
		repo.Owner,
		repo.Name,
		repo.FullName,
		repo.Private,
		repo.Meta,
	).Scan(&repo.Id, &repo.CreatedAt, &repo.UpdatedAt)
}

// GetById returns a workspace GitHub repo by id.
func (r *RepositoryRepositoryPostgres) GetById(id Id) (*Repository, error) {
	item := &Repository{}
	err := r.db.QueryRow(
		`SELECT id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta, created_at, updated_at
		FROM repositories
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.InstallationId,
		&item.GitHubId,
		&item.NodeId,
		&item.Owner,
		&item.Name,
		&item.FullName,
		&item.Private,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByGitHubId returns a workspace GitHub repo by GitHub repository id.
func (r *RepositoryRepositoryPostgres) GetByGitHubId(githubId int64) (*Repository, error) {
	item := &Repository{}
	err := r.db.QueryRow(
		`SELECT id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta, created_at, updated_at
		FROM repositories
		WHERE github_id = $1`,
		githubId,
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.InstallationId,
		&item.GitHubId,
		&item.NodeId,
		&item.Owner,
		&item.Name,
		&item.FullName,
		&item.Private,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates a workspace GitHub repo.
func (r *RepositoryRepositoryPostgres) Update(repo *Repository) error {
	_, err := r.db.Exec(
		`UPDATE repositories
		SET
			workspace_id = $1,
			installation_id = $2,
			github_id = $3,
			node_id = $4,
			owner = $5,
			name = $6,
			full_name = $7,
			private = $8,
			meta = $9,
			updated_at = $10
		WHERE id = $11`,
		repo.WorkspaceId.String(),
		repo.InstallationId,
		repo.GitHubId,
		repo.NodeId,
		repo.Owner,
		repo.Name,
		repo.FullName,
		repo.Private,
		repo.Meta,
		time.Now().UTC(),
		repo.Id.String(),
	)
	return err
}

// Delete deletes a workspace GitHub repo row.
func (r *RepositoryRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(`DELETE FROM repositories WHERE id = $1`, id.String())
	return err
}

// DeleteByGitHubId deletes a workspace GitHub repo by GitHub repository id.
func (r *RepositoryRepositoryPostgres) DeleteByGitHubId(githubId int64) error {
	_, err := r.db.Exec(`DELETE FROM repositories WHERE github_id = $1`, githubId)
	return err
}

// DeleteByInstallationId deletes all repos for a GitHub App installation.
func (r *RepositoryRepositoryPostgres) DeleteByInstallationId(installationId int64) error {
	_, err := r.db.Exec(`DELETE FROM repositories WHERE installation_id = $1`, installationId)
	return err
}

// ListByWorkspaceId lists GitHub repos by workspace id.
func (r *RepositoryRepositoryPostgres) ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*Repository, error) {
	rows, err := r.db.Query(
		`SELECT id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta, created_at, updated_at
		FROM repositories
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

	var list []*Repository
	for rows.Next() {
		item := &Repository{}
		if err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.InstallationId,
			&item.GitHubId,
			&item.NodeId,
			&item.Owner,
			&item.Name,
			&item.FullName,
			&item.Private,
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

// ListByInstallationId lists GitHub repos by GitHub App installation id.
func (r *RepositoryRepositoryPostgres) ListByInstallationId(installationId int64) ([]*Repository, error) {
	rows, err := r.db.Query(
		`SELECT id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta, created_at, updated_at
		FROM repositories
		WHERE installation_id = $1
		ORDER BY created_at DESC`,
		installationId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Repository
	for rows.Next() {
		item := &Repository{}
		if err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.InstallationId,
			&item.GitHubId,
			&item.NodeId,
			&item.Owner,
			&item.Name,
			&item.FullName,
			&item.Private,
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

// CountByWorkspaceId counts GitHub repos by workspace id.
func (r *RepositoryRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM repositories
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// RepositoryMeta is a single row in the repositories_meta table.
type RepositoryMeta struct {
	Id           Id
	RepositoryId Id
	Key          string
	Value        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RepositoryMetaRepository is the interface for workspace GitHub repo metadata CRUD.
type RepositoryMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*RepositoryMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByRepositoryId(id Id) ([]*RepositoryMeta, error)
	Upsert(id Id, key, value string) error
}

type RepositoryMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewRepositoryMetaRepository returns the repository for workspace GitHub repo metadata.
func NewRepositoryMetaRepository(db *sql.DB) RepositoryMetaRepository {
	return &RepositoryMetaRepositoryPostgres{db: db}
}

// Create inserts a workspace GitHub repo metadata row.
func (r *RepositoryMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`INSERT INTO repositories_meta (id, repository_id, key, value)
		VALUES ($1, $2, $3, to_jsonb($4::text))`,
		metaId.String(), id.String(), key, value,
	)
	return err
}

// Get returns workspace GitHub repo metadata by key.
func (r *RepositoryMetaRepositoryPostgres) Get(id Id, key string) (*RepositoryMeta, error) {
	meta := &RepositoryMeta{}
	err := r.db.QueryRow(
		`SELECT id, repository_id, key, value #>> '{}', created_at, updated_at
		FROM repositories_meta
		WHERE repository_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(&meta.Id, &meta.RepositoryId, &meta.Key, &meta.Value, &meta.CreatedAt, &meta.UpdatedAt)
	if isNotFound(err) {
		return nil, nil
	}
	return meta, err
}

// Update updates an existing workspace GitHub repo metadata row.
func (r *RepositoryMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE repositories_meta
		SET value = to_jsonb($1::text), updated_at = $2
		WHERE repository_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)
	return err
}

// Delete deletes a workspace GitHub repo metadata row.
func (r *RepositoryMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM repositories_meta
		WHERE repository_id = $1 AND key = $2`,
		id.String(), key,
	)
	return err
}

// ListByRepositoryId lists workspace GitHub repo metadata rows by repo id.
func (r *RepositoryMetaRepositoryPostgres) ListByRepositoryId(id Id) ([]*RepositoryMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, repository_id, key, value #>> '{}', created_at, updated_at
		FROM repositories_meta
		WHERE repository_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*RepositoryMeta
	for rows.Next() {
		meta := &RepositoryMeta{}
		err := rows.Scan(&meta.Id, &meta.RepositoryId, &meta.Key, &meta.Value, &meta.CreatedAt, &meta.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, meta)
	}
	return list, rows.Err()
}

// Upsert creates or updates workspace GitHub repo metadata.
func (r *RepositoryMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}
	return r.Update(id, key, value)
}
