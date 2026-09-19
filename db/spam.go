// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// RepositorySpamUser is a GitHub user blocked as spam on a repository.
type RepositorySpamUser struct {
	Id           Id
	RepositoryId Id
	GitHubId     *int64
	Username     *string
	Email        *string
	Meta         *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RepositorySpamUserRepository is the interface for repository spam user CRUD.
type RepositorySpamUserRepository interface {
	Create(user *RepositorySpamUser) error
	Upsert(user *RepositorySpamUser) error
	GetById(id Id) (*RepositorySpamUser, error)
	GetByGitHubId(repositoryId Id, githubId int64) (*RepositorySpamUser, error)
	Delete(id Id) error
	ListByRepositoryId(repositoryId Id) ([]*RepositorySpamUser, error)
}

type RepositorySpamUserRepositoryPostgres struct {
	db *sql.DB
}

// NewRepositorySpamUserRepository returns the repository for repository spam users.
func NewRepositorySpamUserRepository(db *sql.DB) RepositorySpamUserRepository {
	return &RepositorySpamUserRepositoryPostgres{db: db}
}

// Create inserts a repository spam user row.
func (r *RepositorySpamUserRepositoryPostgres) Create(user *RepositorySpamUser) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	user.Id = id

	return r.db.QueryRow(
		`INSERT INTO repository_spam_users (id, repository_id, github_id, username, email, meta)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`,
		user.Id.String(),
		user.RepositoryId.String(),
		user.GitHubId,
		user.Username,
		user.Email,
		user.Meta,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// Upsert inserts or updates a repository spam user by repository and GitHub id.
func (r *RepositorySpamUserRepositoryPostgres) Upsert(user *RepositorySpamUser) error {
	id, err := NewId()
	if err != nil {
		return err
	}

	return r.db.QueryRow(
		`INSERT INTO repository_spam_users (id, repository_id, github_id, username, email, meta)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (repository_id, github_id) DO UPDATE SET
			username = EXCLUDED.username,
			email = EXCLUDED.email,
			meta = EXCLUDED.meta,
			updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
		RETURNING id, created_at, updated_at`,
		id.String(),
		user.RepositoryId.String(),
		user.GitHubId,
		user.Username,
		user.Email,
		user.Meta,
	).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
}

// GetById returns a repository spam user by id.
func (r *RepositorySpamUserRepositoryPostgres) GetById(id Id) (*RepositorySpamUser, error) {
	item := &RepositorySpamUser{}
	err := r.db.QueryRow(
		`SELECT id, repository_id, github_id, username, email, meta, created_at, updated_at
		FROM repository_spam_users
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.RepositoryId,
		&item.GitHubId,
		&item.Username,
		&item.Email,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByGitHubId returns a repository spam user by GitHub id.
func (r *RepositorySpamUserRepositoryPostgres) GetByGitHubId(repositoryId Id, githubId int64) (*RepositorySpamUser, error) {
	item := &RepositorySpamUser{}
	err := r.db.QueryRow(
		`SELECT id, repository_id, github_id, username, email, meta, created_at, updated_at
		FROM repository_spam_users
		WHERE repository_id = $1 AND github_id = $2`,
		repositoryId.String(),
		githubId,
	).Scan(
		&item.Id,
		&item.RepositoryId,
		&item.GitHubId,
		&item.Username,
		&item.Email,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Delete deletes a repository spam user row.
func (r *RepositorySpamUserRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(`DELETE FROM repository_spam_users WHERE id = $1`, id.String())
	return err
}

// ListByRepositoryId lists spam users by repository id.
func (r *RepositorySpamUserRepositoryPostgres) ListByRepositoryId(repositoryId Id) ([]*RepositorySpamUser, error) {
	rows, err := r.db.Query(
		`SELECT id, repository_id, github_id, username, email, meta, created_at, updated_at
		FROM repository_spam_users
		WHERE repository_id = $1
		ORDER BY created_at DESC`,
		repositoryId.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*RepositorySpamUser
	for rows.Next() {
		item := &RepositorySpamUser{}
		if err := rows.Scan(
			&item.Id,
			&item.RepositoryId,
			&item.GitHubId,
			&item.Username,
			&item.Email,
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
