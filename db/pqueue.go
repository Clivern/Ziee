// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const (
	PQueuePriorityHigh   = "high"
	PQueuePriorityMedium = "medium"
	PQueuePriorityLow    = "low"

	PQueueStatusDraft    = "draft"
	PQueueStatusOpened   = "opened"
	PQueueStatusClosed   = "closed"
	PQueueStatusReopened = "reopened"
	PQueueStatusQueued   = "queued"
	PQueueStatusChecking = "checking"
	PQueueStatusSuccess  = "success"
	PQueueStatusFailed   = "failed"
	PQueueStatusMerged   = "merged"
	PQueueStatusDequeued = "dequeued"
)

// PQueue is a pull request waiting in a repository merge queue.
type PQueue struct {
	Id        Id
	RepoId    Id
	RemoteId  int64
	Priority  string
	Rank      int
	Status    string
	Checksum  string
	Meta      string
	OpenedAt  *time.Time
	MergedAt  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PQueueRepository is the interface for merge-queue PR persistence.
type PQueueRepository interface {
	Create(item *PQueue) error
	GetById(id Id) (*PQueue, error)
	GetByRepoIdAndRemoteId(repoId Id, remoteId int64) (*PQueue, error)
	Update(item *PQueue) error
	ListUnmergedByRepoId(repoId Id) ([]*PQueue, error)
	ListQueuedByRepoId(repoId Id, limit int) ([]*PQueue, error)
}

type PQueueRepositoryPostgres struct {
	db *sql.DB
}

// NewPQueueRepository returns the repository for the pqueue table.
func NewPQueueRepository(db *sql.DB) PQueueRepository {
	return &PQueueRepositoryPostgres{db: db}
}

// Create inserts a merge-queue PR row.
func (r *PQueueRepositoryPostgres) Create(item *PQueue) error {
	id, err := NewId()
	if err != nil {
		return err
	}

	item.Id = id

	if item.Priority == "" {
		item.Priority = PQueuePriorityMedium
	}
	if item.Meta == "" {
		item.Meta = "{}"
	}

	return r.db.QueryRow(
		`INSERT INTO pqueue
		(id, repo_id, remote_id, priority, rank, status, checksum, meta, opened_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at`,
		item.Id.String(),
		item.RepoId.String(),
		item.RemoteId,
		item.Priority,
		item.Rank,
		item.Status,
		item.Checksum,
		item.Meta,
		item.OpenedAt,
		item.MergedAt,
	).Scan(&item.CreatedAt, &item.UpdatedAt)
}

// GetById returns a merge-queue PR by id.
func (r *PQueueRepositoryPostgres) GetById(id Id) (*PQueue, error) {
	item := &PQueue{}
	err := r.db.QueryRow(
		`SELECT id, repo_id, remote_id, priority, rank, status, checksum, meta, opened_at, merged_at, created_at, updated_at
		FROM pqueue
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.RepoId,
		&item.RemoteId,
		&item.Priority,
		&item.Rank,
		&item.Status,
		&item.Checksum,
		&item.Meta,
		&item.OpenedAt,
		&item.MergedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// GetByRepoIdAndRemoteId returns a merge-queue PR by repo and remote id.
func (r *PQueueRepositoryPostgres) GetByRepoIdAndRemoteId(repoId Id, remoteId int64) (*PQueue, error) {
	item := &PQueue{}
	err := r.db.QueryRow(
		`SELECT id, repo_id, remote_id, priority, rank, status, checksum, meta, opened_at, merged_at, created_at, updated_at
		FROM pqueue
		WHERE repo_id = $1 AND remote_id = $2`,
		repoId.String(),
		remoteId,
	).Scan(
		&item.Id,
		&item.RepoId,
		&item.RemoteId,
		&item.Priority,
		&item.Rank,
		&item.Status,
		&item.Checksum,
		&item.Meta,
		&item.OpenedAt,
		&item.MergedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// Update updates a merge-queue PR row.
func (r *PQueueRepositoryPostgres) Update(item *PQueue) error {
	_, err := r.db.Exec(
		`UPDATE pqueue
		SET
			repo_id = $1,
			remote_id = $2,
			priority = $3,
			rank = $4,
			status = $5,
			checksum = $6,
			meta = $7,
			opened_at = $8,
			merged_at = $9,
			updated_at = $10
		WHERE id = $11`,
		item.RepoId.String(),
		item.RemoteId,
		item.Priority,
		item.Rank,
		item.Status,
		item.Checksum,
		item.Meta,
		item.OpenedAt,
		item.MergedAt,
		time.Now().UTC(),
		item.Id.String(),
	)

	return err
}

// ListUnmergedByRepoId lists non-merged PRs for a repo, high priority first then by rank.
func (r *PQueueRepositoryPostgres) ListUnmergedByRepoId(repoId Id) ([]*PQueue, error) {
	rows, err := r.db.Query(
		`SELECT id, repo_id, remote_id, priority, rank, status, checksum, meta, opened_at, merged_at, created_at, updated_at
		FROM pqueue
		WHERE repo_id = $1 AND status != $2
		ORDER BY
			CASE priority
				WHEN 'high' THEN 1
				WHEN 'medium' THEN 2
				WHEN 'low' THEN 3
			END,
			rank`,
		repoId.String(),
		PQueueStatusMerged,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list []*PQueue
	for rows.Next() {
		item := &PQueue{}
		if err := rows.Scan(
			&item.Id,
			&item.RepoId,
			&item.RemoteId,
			&item.Priority,
			&item.Rank,
			&item.Status,
			&item.Checksum,
			&item.Meta,
			&item.OpenedAt,
			&item.MergedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	return list, rows.Err()
}

// ListQueuedByRepoId returns the top queued PRs for a repo, high priority first then by rank.
func (r *PQueueRepositoryPostgres) ListQueuedByRepoId(repoId Id, limit int) ([]*PQueue, error) {
	rows, err := r.db.Query(
		`SELECT id, repo_id, remote_id, priority, rank, status, checksum, meta, opened_at, merged_at, created_at, updated_at
		FROM pqueue
		WHERE repo_id = $1 AND status = $2
		ORDER BY
			CASE priority
				WHEN 'high' THEN 1
				WHEN 'medium' THEN 2
				WHEN 'low' THEN 3
			END,
			rank
		LIMIT $3`,
		repoId.String(),
		PQueueStatusQueued,
		limit,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list []*PQueue
	for rows.Next() {
		item := &PQueue{}
		if err := rows.Scan(
			&item.Id,
			&item.RepoId,
			&item.RemoteId,
			&item.Priority,
			&item.Rank,
			&item.Status,
			&item.Checksum,
			&item.Meta,
			&item.OpenedAt,
			&item.MergedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	return list, rows.Err()
}
