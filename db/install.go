// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const (
	GitHubInstallationStatusPending   = "pending"
	GitHubInstallationStatusFailed    = "failed"
	GitHubInstallationStatusProcessed = "processed"
)

// GitHubInstallation is a GitHub App installation waiting to be attached to a workspace.
type GitHubInstallation struct {
	Id                  Id
	GitHubId            int64
	GitHubUserId        string
	AccountId           int64
	AccountLogin        string
	AccountType         string
	WorkspaceId         Id
	Status              string
	RepositorySelection string
	HTMLURL             string
	Meta                *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// GitHubInstallationRepository is the interface for GitHub App installation CRUD.
type GitHubInstallationRepository interface {
	Upsert(installation *GitHubInstallation) error
	GetById(id Id) (*GitHubInstallation, error)
	GetByGitHubId(githubId int64) (*GitHubInstallation, error)
	ListPendingByGitHubUserId(githubUserId string) ([]*GitHubInstallation, error)
	Attach(id, workspaceId Id) error
	UpdateStatus(id Id, status string) error
	DeleteByGitHubId(githubId int64) error
}

type GitHubInstallationRepositoryPostgres struct {
	db *sql.DB
}

// NewGitHubInstallationRepository returns the repository for GitHub App installations.
func NewGitHubInstallationRepository(db *sql.DB) GitHubInstallationRepository {
	return &GitHubInstallationRepositoryPostgres{db: db}
}

// Upsert inserts or updates a GitHub App installation by GitHub installation id.
func (r *GitHubInstallationRepositoryPostgres) Upsert(installation *GitHubInstallation) error {
	id, err := NewId()
	if err != nil {
		return err
	}

	return r.db.QueryRow(
		`INSERT INTO github_installations
		(id, github_id, github_user_id, account_id, account_login, account_type, repository_selection, html_url, meta)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (github_id) DO UPDATE SET
			github_user_id = EXCLUDED.github_user_id,
			account_id = EXCLUDED.account_id,
			account_login = EXCLUDED.account_login,
			account_type = EXCLUDED.account_type,
			repository_selection = EXCLUDED.repository_selection,
			html_url = EXCLUDED.html_url,
			meta = EXCLUDED.meta,
			updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
		RETURNING id, workspace_id, status, created_at, updated_at`,
		id.String(),
		installation.GitHubId,
		installation.GitHubUserId,
		installation.AccountId,
		installation.AccountLogin,
		installation.AccountType,
		installation.RepositorySelection,
		installation.HTMLURL,
		installation.Meta,
	).Scan(
		&installation.Id,
		&installation.WorkspaceId,
		&installation.Status,
		&installation.CreatedAt,
		&installation.UpdatedAt,
	)
}

// GetById returns a GitHub App installation by id.
func (r *GitHubInstallationRepositoryPostgres) GetById(id Id) (*GitHubInstallation, error) {
	item := &GitHubInstallation{}
	err := r.db.QueryRow(
		`SELECT id, github_id, github_user_id, account_id, account_login, account_type, workspace_id, status, repository_selection, html_url, meta, created_at, updated_at
		FROM github_installations
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.GitHubId,
		&item.GitHubUserId,
		&item.AccountId,
		&item.AccountLogin,
		&item.AccountType,
		&item.WorkspaceId,
		&item.Status,
		&item.RepositorySelection,
		&item.HTMLURL,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// GetByGitHubId returns a GitHub App installation by GitHub installation id.
func (r *GitHubInstallationRepositoryPostgres) GetByGitHubId(githubId int64) (*GitHubInstallation, error) {
	item := &GitHubInstallation{}
	err := r.db.QueryRow(
		`SELECT id, github_id, github_user_id, account_id, account_login, account_type, workspace_id, status, repository_selection, html_url, meta, created_at, updated_at
		FROM github_installations
		WHERE github_id = $1`,
		githubId,
	).Scan(
		&item.Id,
		&item.GitHubId,
		&item.GitHubUserId,
		&item.AccountId,
		&item.AccountLogin,
		&item.AccountType,
		&item.WorkspaceId,
		&item.Status,
		&item.RepositorySelection,
		&item.HTMLURL,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// ListPendingByGitHubUserId lists pending GitHub App installations for a GitHub user ident.
func (r *GitHubInstallationRepositoryPostgres) ListPendingByGitHubUserId(githubUserId string) ([]*GitHubInstallation, error) {
	rows, err := r.db.Query(
		`SELECT id, github_id, github_user_id, account_id, account_login, account_type, workspace_id, status, repository_selection, html_url, meta, created_at, updated_at
		FROM github_installations
		WHERE github_user_id = $1 AND status = $2
		ORDER BY created_at DESC`,
		githubUserId,
		GitHubInstallationStatusPending,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*GitHubInstallation{}
	for rows.Next() {
		item := &GitHubInstallation{}
		if err := rows.Scan(
			&item.Id,
			&item.GitHubId,
			&item.GitHubUserId,
			&item.AccountId,
			&item.AccountLogin,
			&item.AccountType,
			&item.WorkspaceId,
			&item.Status,
			&item.RepositorySelection,
			&item.HTMLURL,
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

// Attach sets the workspace for a GitHub App installation.
func (r *GitHubInstallationRepositoryPostgres) Attach(id, workspaceId Id) error {
	_, err := r.db.Exec(
		`UPDATE github_installations
		SET workspace_id = $1, status = $2, updated_at = $3
		WHERE id = $4`,
		workspaceId.String(),
		GitHubInstallationStatusProcessed,
		time.Now().UTC(),
		id.String(),
	)

	return err
}

// UpdateStatus sets the processing status of a GitHub App installation.
func (r *GitHubInstallationRepositoryPostgres) UpdateStatus(id Id, status string) error {
	_, err := r.db.Exec(
		`UPDATE github_installations
		SET status = $1, updated_at = $2
		WHERE id = $3`,
		status,
		time.Now().UTC(),
		id.String(),
	)

	return err
}

// DeleteByGitHubId deletes a GitHub App installation by GitHub installation id.
func (r *GitHubInstallationRepositoryPostgres) DeleteByGitHubId(githubId int64) error {
	_, err := r.db.Exec(`DELETE FROM github_installations WHERE github_id = $1`, githubId)

	return err
}
