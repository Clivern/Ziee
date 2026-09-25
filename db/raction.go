// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const (
	RepositoryActionStatusCompleted = "completed"
	RepositoryActionStatusFailed    = "failed"
	RepositoryActionStatusNoOp      = "no_op"
)

// RepositoryAction is one evaluated plan for a repository issue or pull request.
type RepositoryAction struct {
	Id           Id
	RepositoryId Id
	AsyncTaskId  *Id
	DeliveryId   string
	EventKind    string
	IssueNumber  int
	Verb         *string
	Plan         string
	Status       string
	Error        *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RepositoryActionRepository is the interface for repository action persistence.
type RepositoryActionRepository interface {
	Create(item *RepositoryAction) error
	GetById(id Id) (*RepositoryAction, error)
	ListByRepositoryIssue(repositoryId Id, issueNumber int) ([]*RepositoryAction, error)
}

type RepositoryActionRepositoryPostgres struct {
	db *sql.DB
}

// NewRepositoryActionRepository returns the repository for repository_actions.
func NewRepositoryActionRepository(db *sql.DB) RepositoryActionRepository {
	return &RepositoryActionRepositoryPostgres{db: db}
}

// Create inserts a repository action row.
func (r *RepositoryActionRepositoryPostgres) Create(item *RepositoryAction) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	item.Id = id

	if item.Plan == "" {
		item.Plan = "[]"
	}

	return r.db.QueryRow(
		`INSERT INTO repository_actions (
			id, repository_id, async_task_id, delivery_id, event_kind,
			issue_number, verb, plan, status, error
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at`,
		item.Id.String(),
		item.RepositoryId.String(),
		item.AsyncTaskId,
		item.DeliveryId,
		item.EventKind,
		item.IssueNumber,
		item.Verb,
		item.Plan,
		item.Status,
		item.Error,
	).Scan(&item.CreatedAt, &item.UpdatedAt)
}

// GetById returns a repository action by id.
func (r *RepositoryActionRepositoryPostgres) GetById(id Id) (*RepositoryAction, error) {
	item := &RepositoryAction{}
	err := r.db.QueryRow(
		`SELECT
			id, repository_id, async_task_id, delivery_id, event_kind,
			issue_number, verb, plan, status, error, created_at, updated_at
		FROM repository_actions
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.RepositoryId,
		&item.AsyncTaskId,
		&item.DeliveryId,
		&item.EventKind,
		&item.IssueNumber,
		&item.Verb,
		&item.Plan,
		&item.Status,
		&item.Error,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// ListByRepositoryIssue lists actions for one issue or pull request, newest first.
func (r *RepositoryActionRepositoryPostgres) ListByRepositoryIssue(repositoryId Id, issueNumber int) ([]*RepositoryAction, error) {
	rows, err := r.db.Query(
		`SELECT
			id, repository_id, async_task_id, delivery_id, event_kind,
			issue_number, verb, plan, status, error, created_at, updated_at
		FROM repository_actions
		WHERE repository_id = $1 AND issue_number = $2
		ORDER BY created_at DESC`,
		repositoryId.String(),
		issueNumber,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*RepositoryAction
	for rows.Next() {
		item := &RepositoryAction{}
		err = rows.Scan(
			&item.Id,
			&item.RepositoryId,
			&item.AsyncTaskId,
			&item.DeliveryId,
			&item.EventKind,
			&item.IssueNumber,
			&item.Verb,
			&item.Plan,
			&item.Status,
			&item.Error,
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
