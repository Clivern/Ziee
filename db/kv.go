// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// KV is a key/value row with an optional expiry.
type KV struct {
	Id        Id
	Key       string
	Value     string
	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// KVRepository is a generic expiry-aware key/value store.
type KVRepository interface {
	Upsert(item *KV) error
	Get(key string) (*KV, error)
	Delete(key string) error
	DeleteExpired() (int64, error)
}

type KVRepositoryPostgres struct {
	db *sql.DB
}

// NewKVRepository returns the repository for the kv table.
func NewKVRepository(db *sql.DB) KVRepository {
	return &KVRepositoryPostgres{db: db}
}

// Upsert inserts or replaces a key/value row.
func (r *KVRepositoryPostgres) Upsert(item *KV) error {
	id, err := NewId()
	if err != nil {
		return err
	}

	return r.db.QueryRow(
		`INSERT INTO kv (id, key, value, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			expires_at = EXCLUDED.expires_at,
			updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
		RETURNING id, created_at, updated_at`,
		id.String(),
		item.Key,
		item.Value,
		item.ExpiresAt,
	).Scan(
		&item.Id,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
}

// Get returns a non-expired value for key. A null expires_at never expires.
func (r *KVRepositoryPostgres) Get(key string) (*KV, error) {
	item := &KV{}
	err := r.db.QueryRow(
		`SELECT id, key, value, expires_at, created_at, updated_at
		FROM kv
		WHERE key = $1 AND (expires_at IS NULL OR expires_at > $2)`,
		key,
		time.Now().UTC(),
	).Scan(
		&item.Id,
		&item.Key,
		&item.Value,
		&item.ExpiresAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}

	return item, err
}

// Delete removes a key/value row.
func (r *KVRepositoryPostgres) Delete(key string) error {
	_, err := r.db.Exec(`DELETE FROM kv WHERE key = $1`, key)
	return err
}

// DeleteExpired removes rows that have passed their expiry.
func (r *KVRepositoryPostgres) DeleteExpired() (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM kv WHERE expires_at IS NOT NULL AND expires_at <= $1`,
		time.Now().UTC(),
	)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
