// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const (
	SandboxStatusActive  = "active"
	SandboxStatusStopped = "stopped"
	SandboxStatusRemoved = "removed"
)

// Sandbox is a repository sandbox with config, usage, and expiry.
type Sandbox struct {
	Id             Id
	RepositoryId   Id
	Config         string
	Usage          string
	Port           int
	Status         string
	Token          string
	RemoteId       string
	ExpiresAt      time.Time
	LastActivityAt time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SandboxRepository is the interface for repository sandbox persistence.
type SandboxRepository interface {
	Create(item *Sandbox) error
	GetById(id Id) (*Sandbox, error)
	GetByToken(token string) (*Sandbox, error)
	GetByRemoteId(remoteId string) (*Sandbox, error)
	ListByRepositoryId(repositoryId Id) ([]*Sandbox, error)
	Touch(id Id) error
	Delete(id Id) error
	DeleteExpired() (int64, error)
}

type SandboxRepositoryPostgres struct {
	db *sql.DB
}

// NewSandboxRepository returns the repository for the sandbox table.
func NewSandboxRepository(db *sql.DB) SandboxRepository {
	return &SandboxRepositoryPostgres{db: db}
}

// Create inserts a sandbox row.
func (r *SandboxRepositoryPostgres) Create(item *Sandbox) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	item.Id = id

	if item.Config == "" {
		item.Config = "{}"
	}
	if item.Usage == "" {
		item.Usage = "{}"
	}
	if item.Status == "" {
		item.Status = SandboxStatusActive
	}
	if item.LastActivityAt.IsZero() {
		item.LastActivityAt = time.Now().UTC()
	}

	return r.db.QueryRow(
		`INSERT INTO sandbox (id, repository_id, config, usage, port, status, token, remote_id, expires_at, last_activity_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at`,
		item.Id.String(),
		item.RepositoryId.String(),
		item.Config,
		item.Usage,
		item.Port,
		item.Status,
		item.Token,
		item.RemoteId,
		item.ExpiresAt,
		item.LastActivityAt,
	).Scan(&item.CreatedAt, &item.UpdatedAt)
}

// GetById returns a non-expired sandbox by id.
func (r *SandboxRepositoryPostgres) GetById(id Id) (*Sandbox, error) {
	item := &Sandbox{}
	err := r.db.QueryRow(
		`SELECT id, repository_id, config, usage, port, status, token, remote_id, expires_at, last_activity_at, created_at, updated_at
		FROM sandbox
		WHERE id = $1 AND expires_at > $2`,
		id.String(),
		time.Now().UTC(),
	).Scan(
		&item.Id,
		&item.RepositoryId,
		&item.Config,
		&item.Usage,
		&item.Port,
		&item.Status,
		&item.Token,
		&item.RemoteId,
		&item.ExpiresAt,
		&item.LastActivityAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// GetByToken returns a non-expired sandbox by token.
func (r *SandboxRepositoryPostgres) GetByToken(token string) (*Sandbox, error) {
	item := &Sandbox{}
	err := r.db.QueryRow(
		`SELECT id, repository_id, config, usage, port, status, token, remote_id, expires_at, last_activity_at, created_at, updated_at
		FROM sandbox
		WHERE token = $1 AND expires_at > $2`,
		token,
		time.Now().UTC(),
	).Scan(
		&item.Id,
		&item.RepositoryId,
		&item.Config,
		&item.Usage,
		&item.Port,
		&item.Status,
		&item.Token,
		&item.RemoteId,
		&item.ExpiresAt,
		&item.LastActivityAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// GetByRemoteId returns a non-expired sandbox by remote id.
func (r *SandboxRepositoryPostgres) GetByRemoteId(remoteId string) (*Sandbox, error) {
	item := &Sandbox{}
	err := r.db.QueryRow(
		`SELECT id, repository_id, config, usage, port, status, token, remote_id, expires_at, last_activity_at, created_at, updated_at
		FROM sandbox
		WHERE remote_id = $1 AND expires_at > $2`,
		remoteId,
		time.Now().UTC(),
	).Scan(
		&item.Id,
		&item.RepositoryId,
		&item.Config,
		&item.Usage,
		&item.Port,
		&item.Status,
		&item.Token,
		&item.RemoteId,
		&item.ExpiresAt,
		&item.LastActivityAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// ListByRepositoryId lists non-expired sandbox rows for a repository.
func (r *SandboxRepositoryPostgres) ListByRepositoryId(repositoryId Id) ([]*Sandbox, error) {
	rows, err := r.db.Query(
		`SELECT id, repository_id, config, usage, port, status, token, remote_id, expires_at, last_activity_at, created_at, updated_at
		FROM sandbox
		WHERE repository_id = $1 AND expires_at > $2
		ORDER BY created_at DESC`,
		repositoryId.String(),
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*Sandbox, 0)
	for rows.Next() {
		item := &Sandbox{}
		err := rows.Scan(
			&item.Id,
			&item.RepositoryId,
			&item.Config,
			&item.Usage,
			&item.Port,
			&item.Status,
			&item.Token,
			&item.RemoteId,
			&item.ExpiresAt,
			&item.LastActivityAt,
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

// Touch updates last activity for a sandbox.
func (r *SandboxRepositoryPostgres) Touch(id Id) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE sandbox
		SET last_activity_at = $1, updated_at = $1
		WHERE id = $2`,
		now,
		id.String(),
	)

	return err
}

// Delete removes a sandbox by id.
func (r *SandboxRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(`DELETE FROM sandbox WHERE id = $1`, id.String())

	return err
}

// DeleteExpired removes sandbox rows that have passed their expiry.
func (r *SandboxRepositoryPostgres) DeleteExpired() (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM sandbox WHERE expires_at <= $1`,
		time.Now().UTC(),
	)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
