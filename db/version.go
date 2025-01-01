// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"

	"github.com/clivern/actx0/pkg/util"
)

// PromptVersion is a single row in the prompt_versions table.
type PromptVersion struct {
	Id            Id
	PromptId      Id
	Version       int
	Type          string
	Content       string
	Config        *string
	Labels        *string
	CommitMessage *string
	CommitHash    string
	Meta          *string
	Status        string
	Production    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// PromptVersionRepository is the interface for prompt version CRUD.
type PromptVersionRepository interface {
	Create(version *PromptVersion) error
	GetById(id Id) (*PromptVersion, error)
	GetByPromptIdAndVersion(promptId Id, version int) (*PromptVersion, error)
	GetLatestByPromptId(promptId Id) (*PromptVersion, error)
	GetProductionByPromptId(promptId Id) (*PromptVersion, error)
	Update(version *PromptVersion) error
	Delete(id Id) error
	ListByPromptId(promptId Id, limit, offset int) ([]*PromptVersion, error)
	CountByPromptId(promptId Id) (int64, error)
	ClearLabelFromPromptId(promptId Id, label string) error
	ClearProductionFromPromptId(promptId Id, exceptId *Id) error
}

type PromptVersionRepositoryPostgres struct {
	db *sql.DB
}

// NewPromptVersionRepository returns a prompt version repository.
func NewPromptVersionRepository(db *sql.DB) PromptVersionRepository {
	return &PromptVersionRepositoryPostgres{db: db}
}

// Create inserts a prompt version row for an existing prompt.
func (r *PromptVersionRepositoryPostgres) Create(version *PromptVersion) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	version.Id = id

	if version.Version == 0 {
		var maxVersion int
		err := r.db.QueryRow(
			`SELECT COALESCE(MAX(version), 0) FROM prompt_versions WHERE prompt_id = $1`,
			version.PromptId.String(),
		).Scan(&maxVersion)
		if err != nil {
			return err
		}
		version.Version = maxVersion + 1
	}

	commitHash, err := util.RandomHash()
	if err != nil {
		return err
	}
	version.CommitHash = commitHash

	if version.Production {
		if err := r.ClearProductionFromPromptId(version.PromptId, nil); err != nil {
			return err
		}
	}

	err = r.db.QueryRow(
		`INSERT INTO prompt_versions (
			id, prompt_id, version, type, content, config, labels,
			commit_message, commit_hash, meta, status, production
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at`,
		version.Id.String(),
		version.PromptId.String(),
		version.Version,
		version.Type,
		version.Content,
		version.Config,
		version.Labels,
		version.CommitMessage,
		version.CommitHash,
		version.Meta,
		version.Status,
		version.Production,
	).Scan(&version.CreatedAt, &version.UpdatedAt)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`UPDATE prompts SET updated_at = $1 WHERE id = $2`,
		time.Now().UTC(),
		version.PromptId.String(),
	)
	return err
}

// GetById returns a prompt version by id.
func (r *PromptVersionRepositoryPostgres) GetById(id Id) (*PromptVersion, error) {
	item := &PromptVersion{}
	err := r.db.QueryRow(
		`SELECT
			id, prompt_id, version, type, content, config, labels,
			commit_message, commit_hash, meta, status, production, created_at, updated_at
		FROM prompt_versions
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.PromptId,
		&item.Version,
		&item.Type,
		&item.Content,
		&item.Config,
		&item.Labels,
		&item.CommitMessage,
		&item.CommitHash,
		&item.Meta,
		&item.Status,
		&item.Production,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByPromptIdAndVersion returns a prompt version by prompt id and version number.
func (r *PromptVersionRepositoryPostgres) GetByPromptIdAndVersion(promptId Id, version int) (*PromptVersion, error) {
	item := &PromptVersion{}
	err := r.db.QueryRow(
		`SELECT
			id, prompt_id, version, type, content, config, labels,
			commit_message, commit_hash, meta, status, production, created_at, updated_at
		FROM prompt_versions
		WHERE prompt_id = $1 AND version = $2`,
		promptId.String(),
		version,
	).Scan(
		&item.Id,
		&item.PromptId,
		&item.Version,
		&item.Type,
		&item.Content,
		&item.Config,
		&item.Labels,
		&item.CommitMessage,
		&item.CommitHash,
		&item.Meta,
		&item.Status,
		&item.Production,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetLatestByPromptId returns the newest version for a prompt.
func (r *PromptVersionRepositoryPostgres) GetLatestByPromptId(promptId Id) (*PromptVersion, error) {
	item := &PromptVersion{}
	err := r.db.QueryRow(
		`SELECT
			id, prompt_id, version, type, content, config, labels,
			commit_message, commit_hash, meta, status, production, created_at, updated_at
		FROM prompt_versions
		WHERE prompt_id = $1
		ORDER BY version DESC
		LIMIT 1`,
		promptId.String(),
	).Scan(
		&item.Id,
		&item.PromptId,
		&item.Version,
		&item.Type,
		&item.Content,
		&item.Config,
		&item.Labels,
		&item.CommitMessage,
		&item.CommitHash,
		&item.Meta,
		&item.Status,
		&item.Production,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetProductionByPromptId returns the production version for a prompt.
func (r *PromptVersionRepositoryPostgres) GetProductionByPromptId(promptId Id) (*PromptVersion, error) {
	item := &PromptVersion{}
	err := r.db.QueryRow(
		`SELECT
			id, prompt_id, version, type, content, config, labels,
			commit_message, commit_hash, meta, status, production, created_at, updated_at
		FROM prompt_versions
		WHERE prompt_id = $1 AND production = TRUE
		LIMIT 1`,
		promptId.String(),
	).Scan(
		&item.Id,
		&item.PromptId,
		&item.Version,
		&item.Type,
		&item.Content,
		&item.Config,
		&item.Labels,
		&item.CommitMessage,
		&item.CommitHash,
		&item.Meta,
		&item.Status,
		&item.Production,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates an existing prompt version row.
func (r *PromptVersionRepositoryPostgres) Update(version *PromptVersion) error {
	_, err := r.db.Exec(
		`UPDATE prompt_versions
		SET
			type = $1,
			content = $2,
			config = $3,
			labels = $4,
			commit_message = $5,
			meta = $6,
			status = $7,
			production = $8,
			updated_at = $9
		WHERE id = $10`,
		version.Type,
		version.Content,
		version.Config,
		version.Labels,
		version.CommitMessage,
		version.Meta,
		version.Status,
		version.Production,
		time.Now().UTC(),
		version.Id.String(),
	)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`UPDATE prompts SET updated_at = $1 WHERE id = $2`,
		time.Now().UTC(),
		version.PromptId.String(),
	)
	return err
}

// Delete deletes a single prompt version; removes the parent prompt when none remain.
func (r *PromptVersionRepositoryPostgres) Delete(id Id) error {
	var promptId Id
	err := r.db.QueryRow(
		`SELECT prompt_id FROM prompt_versions WHERE id = $1`,
		id.String(),
	).Scan(&promptId)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`DELETE FROM prompt_versions WHERE id = $1`, id.String())
	if err != nil {
		return err
	}

	var remaining int
	err = r.db.QueryRow(
		`SELECT COUNT(*) FROM prompt_versions WHERE prompt_id = $1`,
		promptId.String(),
	).Scan(&remaining)
	if err != nil {
		return err
	}
	if remaining == 0 {
		_, err = r.db.Exec(`DELETE FROM prompts WHERE id = $1`, promptId.String())
	}
	return err
}

// ClearLabelFromPromptId removes a label from all versions of a prompt.
func (r *PromptVersionRepositoryPostgres) ClearLabelFromPromptId(promptId Id, label string) error {
	_, err := r.db.Exec(
		`UPDATE prompt_versions
		SET
			labels = (
				SELECT CASE
					WHEN COUNT(*) = 0 THEN NULL
					ELSE jsonb_agg(to_jsonb(x))
				END
				FROM jsonb_array_elements_text(COALESCE(labels, '[]'::jsonb)) AS t(x)
				WHERE x <> $1
			),
			updated_at = $2
		WHERE prompt_id = $3
			AND labels @> jsonb_build_array($1)`,
		label,
		time.Now().UTC(),
		promptId.String(),
	)
	return err
}

// ClearProductionFromPromptId clears the production flag on sibling versions.
func (r *PromptVersionRepositoryPostgres) ClearProductionFromPromptId(promptId Id, exceptId *Id) error {
	query := `UPDATE prompt_versions SET production = FALSE, updated_at = $1 WHERE prompt_id = $2 AND production = TRUE`
	args := []interface{}{time.Now().UTC(), promptId.String()}
	if exceptId != nil {
		query += ` AND id != $3`
		args = append(args, exceptId.String())
	}
	_, err := r.db.Exec(query, args...)
	return err
}

// ListByPromptId lists versions for a prompt.
func (r *PromptVersionRepositoryPostgres) ListByPromptId(promptId Id, limit, offset int) ([]*PromptVersion, error) {
	rows, err := r.db.Query(
		`SELECT
			id, prompt_id, version, type, content, config, labels,
			commit_message, commit_hash, meta, status, production, created_at, updated_at
		FROM prompt_versions
		WHERE prompt_id = $1
		ORDER BY version DESC
		LIMIT $2 OFFSET $3`,
		promptId.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*PromptVersion
	for rows.Next() {
		item := &PromptVersion{}
		err := rows.Scan(
			&item.Id,
			&item.PromptId,
			&item.Version,
			&item.Type,
			&item.Content,
			&item.Config,
			&item.Labels,
			&item.CommitMessage,
			&item.CommitHash,
			&item.Meta,
			&item.Status,
			&item.Production,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// CountByPromptId counts versions for a prompt.
func (r *PromptVersionRepositoryPostgres) CountByPromptId(promptId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM prompt_versions WHERE prompt_id = $1`,
		promptId.String(),
	).Scan(&count)
	return count, err
}
