// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// AsyncTask is a single row in the async_tasks table.
type AsyncTask struct {
	Id          Id
	WorkspaceId Id
	Type        string
	Status      string
	Payload     *string
	Result      *string
	Error       *string
	Attempts    int
	Priority    int
	RunAt       *time.Time
	LockedAt    *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AsyncTaskRepository is the interface for async task CRUD.
type AsyncTaskRepository interface {
	Create(task *AsyncTask) error
	GetById(id Id) (*AsyncTask, error)
	Update(task *AsyncTask) error
	Delete(id Id) error
	ListByStatus(status string, limit, offset int) ([]*AsyncTask, error)
	ListDue(status string, now time.Time, limit int) ([]*AsyncTask, error)
	ClaimDue(now time.Time, limit int) ([]*AsyncTask, error)
	ReleaseClaim(id Id) error
	CountByStatus(status string) (int64, error)
	DeleteCompletedOlderThan(before time.Time) (int64, error)
}

type AsyncTaskRepositoryPostgres struct {
	db *sql.DB
}

// AsyncTaskMeta is a single row in the async_tasks_meta table.
type AsyncTaskMeta struct {
	Id          Id
	AsyncTaskId Id
	Key         string
	Value       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AsyncTaskMetaRepository is the interface for async task metadata CRUD.
type AsyncTaskMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*AsyncTaskMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByAsyncTaskId(id Id) ([]*AsyncTaskMeta, error)
	Upsert(id Id, key, value string) error
}

type AsyncTaskMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewAsyncTaskRepository returns an async task repository.
func NewAsyncTaskRepository(db *sql.DB) AsyncTaskRepository {
	return &AsyncTaskRepositoryPostgres{db: db}
}

// Create inserts an async task row.
func (r *AsyncTaskRepositoryPostgres) Create(task *AsyncTask) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	task.Id = id

	err = r.db.QueryRow(
		`INSERT INTO async_tasks (
			id, workspace_id, type, status, payload, result, error, attempts, priority,
			run_at, locked_at, completed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at`,
		task.Id.String(),
		task.WorkspaceId,
		task.Type,
		task.Status,
		task.Payload,
		task.Result,
		task.Error,
		task.Attempts,
		task.Priority,
		task.RunAt,
		task.LockedAt,
		task.CompletedAt,
	).Scan(&task.CreatedAt, &task.UpdatedAt)
	return err
}

// GetById returns an async task by id.
func (r *AsyncTaskRepositoryPostgres) GetById(id Id) (*AsyncTask, error) {
	item := &AsyncTask{}
	err := r.db.QueryRow(
		`SELECT
			id, workspace_id, type, status, payload, result, error, attempts, priority,
			run_at, locked_at, completed_at, created_at, updated_at
		FROM async_tasks
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.Type,
		&item.Status,
		&item.Payload,
		&item.Result,
		&item.Error,
		&item.Attempts,
		&item.Priority,
		&item.RunAt,
		&item.LockedAt,
		&item.CompletedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates an existing async task row.
func (r *AsyncTaskRepositoryPostgres) Update(task *AsyncTask) error {
	_, err := r.db.Exec(
		`UPDATE async_tasks
		SET
			workspace_id = $1,
			type = $2,
			status = $3,
			payload = $4,
			result = $5,
			error = $6,
			attempts = $7,
			priority = $8,
			run_at = $9,
			locked_at = $10,
			completed_at = $11,
			updated_at = $12
		WHERE id = $13`,
		task.WorkspaceId,
		task.Type,
		task.Status,
		task.Payload,
		task.Result,
		task.Error,
		task.Attempts,
		task.Priority,
		task.RunAt,
		task.LockedAt,
		task.CompletedAt,
		time.Now().UTC(),
		task.Id.String(),
	)
	return err
}

// Delete deletes an async task row.
func (r *AsyncTaskRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(`DELETE FROM async_tasks WHERE id = $1`, id.String())
	return err
}

// ListByStatus lists async task rows by status.
func (r *AsyncTaskRepositoryPostgres) ListByStatus(status string, limit, offset int) ([]*AsyncTask, error) {
	rows, err := r.db.Query(
		`SELECT
			id, workspace_id, type, status, payload, result, error, attempts, priority,
			run_at, locked_at, completed_at, created_at, updated_at
		FROM async_tasks
		WHERE status = $1
		ORDER BY priority DESC, created_at ASC
		LIMIT $2 OFFSET $3`,
		status,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*AsyncTask
	for rows.Next() {
		item := &AsyncTask{}
		if err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.Type,
			&item.Status,
			&item.Payload,
			&item.Result,
			&item.Error,
			&item.Attempts,
			&item.Priority,
			&item.RunAt,
			&item.LockedAt,
			&item.CompletedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// ListDue lists pending tasks that are due to run.
func (r *AsyncTaskRepositoryPostgres) ListDue(status string, now time.Time, limit int) ([]*AsyncTask, error) {
	rows, err := r.db.Query(
		`SELECT
			id, workspace_id, type, status, payload, result, error, attempts, priority,
			run_at, locked_at, completed_at, created_at, updated_at
		FROM async_tasks
		WHERE status = $1
		ORDER BY priority DESC, created_at ASC
		LIMIT $2`,
		status,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*AsyncTask
	for rows.Next() {
		item := &AsyncTask{}
		if err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.Type,
			&item.Status,
			&item.Payload,
			&item.Result,
			&item.Error,
			&item.Attempts,
			&item.Priority,
			&item.RunAt,
			&item.LockedAt,
			&item.CompletedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// ClaimDue atomically claims pending tasks due for execution.
func (r *AsyncTaskRepositoryPostgres) ClaimDue(now time.Time, limit int) ([]*AsyncTask, error) {
	now = now.UTC()
	rows, err := r.db.Query(
		`UPDATE async_tasks
		SET status = 'running', run_at = $1, locked_at = $1, updated_at = $1
		WHERE id IN (
			SELECT id FROM async_tasks
			WHERE status = 'pending'
			ORDER BY priority DESC, created_at ASC
			LIMIT $2
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, workspace_id, type, status, payload, result, error, attempts, priority,
			run_at, locked_at, completed_at, created_at, updated_at`,
		now,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*AsyncTask
	for rows.Next() {
		item := &AsyncTask{}
		if err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.Type,
			&item.Status,
			&item.Payload,
			&item.Result,
			&item.Error,
			&item.Attempts,
			&item.Priority,
			&item.RunAt,
			&item.LockedAt,
			&item.CompletedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// ReleaseClaim returns a claimed task to pending status.
func (r *AsyncTaskRepositoryPostgres) ReleaseClaim(id Id) error {
	_, err := r.db.Exec(
		`UPDATE async_tasks
		SET status = 'pending', run_at = NULL, locked_at = NULL, updated_at = $1
		WHERE id = $2 AND status = 'running'`,
		time.Now().UTC(),
		id.String(),
	)
	return err
}

// CountByStatus counts async task rows by status.
func (r *AsyncTaskRepositoryPostgres) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM async_tasks
		WHERE status = $1`,
		status,
	).Scan(&count)
	return count, err
}

// DeleteCompletedOlderThan deletes completed tasks older than the cutoff.
func (r *AsyncTaskRepositoryPostgres) DeleteCompletedOlderThan(before time.Time) (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM async_tasks
		WHERE status = $1 AND completed_at < $2`,
		"completed",
		before.UTC(),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// NewAsyncTaskMetaRepository returns the repository for async task metadata.
func NewAsyncTaskMetaRepository(db *sql.DB) AsyncTaskMetaRepository {
	return &AsyncTaskMetaRepositoryPostgres{db: db}
}

// Create inserts an async task metadata row.
func (r *AsyncTaskMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO async_tasks_meta (id, async_task_id, key, value)
		VALUES ($1, $2, $3, $4)`,
		metaId.String(), id.String(), key, value,
	)
	return err
}

// Get returns async task metadata by key.
func (r *AsyncTaskMetaRepositoryPostgres) Get(id Id, key string) (*AsyncTaskMeta, error) {
	meta := &AsyncTaskMeta{}
	err := r.db.QueryRow(
		`SELECT id, async_task_id, key, value, created_at, updated_at
		FROM async_tasks_meta
		WHERE async_task_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(
		&meta.Id,
		&meta.AsyncTaskId,
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

// Update updates an existing async task metadata row.
func (r *AsyncTaskMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE async_tasks_meta
		SET value = $1, updated_at = $2
		WHERE async_task_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)
	return err
}

// Delete deletes an async task metadata row.
func (r *AsyncTaskMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM async_tasks_meta
		WHERE async_task_id = $1 AND key = $2`,
		id.String(), key,
	)
	return err
}

// ListByAsyncTaskId lists async task metadata rows by async task id.
func (r *AsyncTaskMetaRepositoryPostgres) ListByAsyncTaskId(id Id) ([]*AsyncTaskMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, async_task_id, key, value, created_at, updated_at
		FROM async_tasks_meta
		WHERE async_task_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*AsyncTaskMeta
	for rows.Next() {
		meta := &AsyncTaskMeta{}
		err := rows.Scan(
			&meta.Id,
			&meta.AsyncTaskId,
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

// Upsert creates or updates async task metadata.
func (r *AsyncTaskMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}
	return r.Update(id, key, value)
}
