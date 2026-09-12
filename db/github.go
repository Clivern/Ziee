// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const WorkspaceGitHubRepoMetaConfigPath = "_github_config_path"

// WorkspaceGitHubRepo is a GitHub repository that installed the Ziee GitHub App.
type WorkspaceGitHubRepo struct {
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

// WorkspaceGitHubRepoRepository is the interface for workspace GitHub repo CRUD.
type WorkspaceGitHubRepoRepository interface {
	Create(repo *WorkspaceGitHubRepo) error
	Upsert(repo *WorkspaceGitHubRepo) error
	GetById(id Id) (*WorkspaceGitHubRepo, error)
	GetByGitHubId(githubId int64) (*WorkspaceGitHubRepo, error)
	Update(repo *WorkspaceGitHubRepo) error
	Delete(id Id) error
	DeleteByGitHubId(githubId int64) error
	DeleteByInstallationId(installationId int64) error
	ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*WorkspaceGitHubRepo, error)
	ListByInstallationId(installationId int64) ([]*WorkspaceGitHubRepo, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
}

type WorkspaceGitHubRepoRepositoryPostgres struct {
	db *sql.DB
}

// NewWorkspaceGitHubRepoRepository returns the repository for workspace GitHub repos.
func NewWorkspaceGitHubRepoRepository(db *sql.DB) WorkspaceGitHubRepoRepository {
	return &WorkspaceGitHubRepoRepositoryPostgres{db: db}
}

// Create inserts a workspace GitHub repo row.
func (r *WorkspaceGitHubRepoRepositoryPostgres) Create(repo *WorkspaceGitHubRepo) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	repo.Id = id

	return r.db.QueryRow(
		`INSERT INTO workspace_github_repos
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
func (r *WorkspaceGitHubRepoRepositoryPostgres) Upsert(repo *WorkspaceGitHubRepo) error {
	id, err := NewId()
	if err != nil {
		return err
	}

	return r.db.QueryRow(
		`INSERT INTO workspace_github_repos
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
func (r *WorkspaceGitHubRepoRepositoryPostgres) GetById(id Id) (*WorkspaceGitHubRepo, error) {
	item := &WorkspaceGitHubRepo{}
	err := r.db.QueryRow(
		`SELECT id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta, created_at, updated_at
		FROM workspace_github_repos
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
func (r *WorkspaceGitHubRepoRepositoryPostgres) GetByGitHubId(githubId int64) (*WorkspaceGitHubRepo, error) {
	item := &WorkspaceGitHubRepo{}
	err := r.db.QueryRow(
		`SELECT id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta, created_at, updated_at
		FROM workspace_github_repos
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
func (r *WorkspaceGitHubRepoRepositoryPostgres) Update(repo *WorkspaceGitHubRepo) error {
	_, err := r.db.Exec(
		`UPDATE workspace_github_repos
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
func (r *WorkspaceGitHubRepoRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(`DELETE FROM workspace_github_repos WHERE id = $1`, id.String())
	return err
}

// DeleteByGitHubId deletes a workspace GitHub repo by GitHub repository id.
func (r *WorkspaceGitHubRepoRepositoryPostgres) DeleteByGitHubId(githubId int64) error {
	_, err := r.db.Exec(`DELETE FROM workspace_github_repos WHERE github_id = $1`, githubId)
	return err
}

// DeleteByInstallationId deletes all repos for a GitHub App installation.
func (r *WorkspaceGitHubRepoRepositoryPostgres) DeleteByInstallationId(installationId int64) error {
	_, err := r.db.Exec(`DELETE FROM workspace_github_repos WHERE installation_id = $1`, installationId)
	return err
}

// ListByWorkspaceId lists GitHub repos by workspace id.
func (r *WorkspaceGitHubRepoRepositoryPostgres) ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*WorkspaceGitHubRepo, error) {
	rows, err := r.db.Query(
		`SELECT id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta, created_at, updated_at
		FROM workspace_github_repos
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

	var list []*WorkspaceGitHubRepo
	for rows.Next() {
		item := &WorkspaceGitHubRepo{}
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
func (r *WorkspaceGitHubRepoRepositoryPostgres) ListByInstallationId(installationId int64) ([]*WorkspaceGitHubRepo, error) {
	rows, err := r.db.Query(
		`SELECT id, workspace_id, installation_id, github_id, node_id, owner, name, full_name, private, meta, created_at, updated_at
		FROM workspace_github_repos
		WHERE installation_id = $1
		ORDER BY created_at DESC`,
		installationId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*WorkspaceGitHubRepo
	for rows.Next() {
		item := &WorkspaceGitHubRepo{}
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
func (r *WorkspaceGitHubRepoRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM workspace_github_repos
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// WorkspaceGitHubRepoMeta is a single row in the workspace_github_repos_meta table.
type WorkspaceGitHubRepoMeta struct {
	Id                    Id
	WorkspaceGitHubRepoId Id
	Key                   string
	Value                 string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// WorkspaceGitHubRepoMetaRepository is the interface for workspace GitHub repo metadata CRUD.
type WorkspaceGitHubRepoMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*WorkspaceGitHubRepoMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByRepoId(id Id) ([]*WorkspaceGitHubRepoMeta, error)
	Upsert(id Id, key, value string) error
}

type WorkspaceGitHubRepoMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewWorkspaceGitHubRepoMetaRepository returns the repository for workspace GitHub repo metadata.
func NewWorkspaceGitHubRepoMetaRepository(db *sql.DB) WorkspaceGitHubRepoMetaRepository {
	return &WorkspaceGitHubRepoMetaRepositoryPostgres{db: db}
}

// Create inserts a workspace GitHub repo metadata row.
func (r *WorkspaceGitHubRepoMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`INSERT INTO workspace_github_repos_meta (id, workspace_github_repo_id, key, value)
		VALUES ($1, $2, $3, to_jsonb($4::text))`,
		metaId.String(), id.String(), key, value,
	)
	return err
}

// Get returns workspace GitHub repo metadata by key.
func (r *WorkspaceGitHubRepoMetaRepositoryPostgres) Get(id Id, key string) (*WorkspaceGitHubRepoMeta, error) {
	meta := &WorkspaceGitHubRepoMeta{}
	err := r.db.QueryRow(
		`SELECT id, workspace_github_repo_id, key, value #>> '{}', created_at, updated_at
		FROM workspace_github_repos_meta
		WHERE workspace_github_repo_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(&meta.Id, &meta.WorkspaceGitHubRepoId, &meta.Key, &meta.Value, &meta.CreatedAt, &meta.UpdatedAt)
	if isNotFound(err) {
		return nil, nil
	}
	return meta, err
}

// Update updates an existing workspace GitHub repo metadata row.
func (r *WorkspaceGitHubRepoMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE workspace_github_repos_meta
		SET value = to_jsonb($1::text), updated_at = $2
		WHERE workspace_github_repo_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)
	return err
}

// Delete deletes a workspace GitHub repo metadata row.
func (r *WorkspaceGitHubRepoMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM workspace_github_repos_meta
		WHERE workspace_github_repo_id = $1 AND key = $2`,
		id.String(), key,
	)
	return err
}

// ListByRepoId lists workspace GitHub repo metadata rows by repo id.
func (r *WorkspaceGitHubRepoMetaRepositoryPostgres) ListByRepoId(id Id) ([]*WorkspaceGitHubRepoMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, workspace_github_repo_id, key, value #>> '{}', created_at, updated_at
		FROM workspace_github_repos_meta
		WHERE workspace_github_repo_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*WorkspaceGitHubRepoMeta
	for rows.Next() {
		meta := &WorkspaceGitHubRepoMeta{}
		err := rows.Scan(&meta.Id, &meta.WorkspaceGitHubRepoId, &meta.Key, &meta.Value, &meta.CreatedAt, &meta.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, meta)
	}
	return list, rows.Err()
}

// Upsert creates or updates workspace GitHub repo metadata.
func (r *WorkspaceGitHubRepoMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}
	return r.Update(id, key, value)
}
