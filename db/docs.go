// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"

	"github.com/samber/lo"
)

const (
	DocumentStatusProcessing = "processing"
	DocumentStatusIndexed    = "indexed"
	DocumentStatusFailed     = "failed"
)

// Document is a single row in the documents table.
type Document struct {
	Id             Id
	InternalId     Id
	WorkspaceId    Id
	Title          string
	Filename       string
	ContentType    string
	Checksum       string
	Size           int64
	CharCount      int64
	Labels         *string
	ProcessedAt    *time.Time
	ChunkingConfig *string
	Meta           *string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// DocumentRepository is the interface for workspace document CRUD.
type DocumentRepository interface {
	Create(document *Document) error
	GetById(id Id) (*Document, error)
	GetByInternalId(internalId Id) (*Document, error)
	Update(document *Document) error
	Delete(id Id) error
	ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*Document, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
	SumSizeByWorkspaceId(workspaceId Id) (int64, error)
}

type DocumentRepositoryPostgres struct {
	db *sql.DB
}

// DocumentMeta is a single row in the documents_meta table.
type DocumentMeta struct {
	Id         Id
	DocumentId Id
	Key        string
	Value      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// DocumentMetaRepository is the interface for workspace document meta CRUD.
type DocumentMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*DocumentMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByDocumentId(id Id) ([]*DocumentMeta, error)
	Upsert(id Id, key, value string) error
}

type DocumentMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewDocumentRepository returns a workspace document repository.
func NewDocumentRepository(db *sql.DB) DocumentRepository {
	return &DocumentRepositoryPostgres{db: db}
}

// Create inserts a workspace document row.
func (r *DocumentRepositoryPostgres) Create(document *Document) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	document.Id = id

	internalId, err := NewId()
	if err != nil {
		return err
	}
	document.InternalId = internalId

	if lo.IsEmpty(document.Status) {
		document.Status = DocumentStatusProcessing
	}

	return r.db.QueryRow(
		`INSERT INTO documents (
			id, internal_id, workspace_id, title, filename, content_type, checksum, size, char_count,
			labels, processed_at, chunking_config, meta, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING created_at, updated_at`,
		document.Id.String(),
		document.InternalId.String(),
		document.WorkspaceId.String(),
		document.Title,
		document.Filename,
		document.ContentType,
		document.Checksum,
		document.Size,
		document.CharCount,
		document.Labels,
		document.ProcessedAt,
		document.ChunkingConfig,
		document.Meta,
		document.Status,
	).Scan(&document.CreatedAt, &document.UpdatedAt)
}

// GetById returns a workspace document by id.
func (r *DocumentRepositoryPostgres) GetById(id Id) (*Document, error) {
	item := &Document{}
	err := r.db.QueryRow(
		`SELECT
			id, internal_id, workspace_id, title, filename, content_type, checksum, size, char_count,
			labels, processed_at, chunking_config, meta, status, created_at, updated_at
		FROM documents
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.InternalId,
		&item.WorkspaceId,
		&item.Title,
		&item.Filename,
		&item.ContentType,
		&item.Checksum,
		&item.Size,
		&item.CharCount,
		&item.Labels,
		&item.ProcessedAt,
		&item.ChunkingConfig,
		&item.Meta,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByInternalId returns a workspace document by internal id.
func (r *DocumentRepositoryPostgres) GetByInternalId(internalId Id) (*Document, error) {
	item := &Document{}
	err := r.db.QueryRow(
		`SELECT
			id, internal_id, workspace_id, title, filename, content_type, checksum, size, char_count,
			labels, processed_at, chunking_config, meta, status, created_at, updated_at
		FROM documents
		WHERE internal_id = $1`,
		internalId.String(),
	).Scan(
		&item.Id,
		&item.InternalId,
		&item.WorkspaceId,
		&item.Title,
		&item.Filename,
		&item.ContentType,
		&item.Checksum,
		&item.Size,
		&item.CharCount,
		&item.Labels,
		&item.ProcessedAt,
		&item.ChunkingConfig,
		&item.Meta,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates an existing workspace document row.
func (r *DocumentRepositoryPostgres) Update(document *Document) error {
	_, err := r.db.Exec(
		`UPDATE documents
		SET
			workspace_id = $1,
			title = $2,
			filename = $3,
			content_type = $4,
			checksum = $5,
			size = $6,
			char_count = $7,
			labels = $8,
			processed_at = $9,
			chunking_config = $10,
			meta = $11,
			status = $12,
			updated_at = $13
		WHERE id = $14`,
		document.WorkspaceId.String(),
		document.Title,
		document.Filename,
		document.ContentType,
		document.Checksum,
		document.Size,
		document.CharCount,
		document.Labels,
		document.ProcessedAt,
		document.ChunkingConfig,
		document.Meta,
		document.Status,
		time.Now().UTC(),
		document.Id.String(),
	)
	return err
}

// Delete deletes a workspace document row.
func (r *DocumentRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(
		`DELETE FROM documents
		WHERE id = $1`,
		id.String(),
	)
	return err
}

// ListByWorkspaceId lists workspace document rows by workspace id.
func (r *DocumentRepositoryPostgres) ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*Document, error) {
	rows, err := r.db.Query(
		`SELECT
			id, internal_id, workspace_id, title, filename, content_type, checksum, size, char_count,
			labels, processed_at, chunking_config, meta, status, created_at, updated_at
		FROM documents
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

	var list []*Document
	for rows.Next() {
		item := &Document{}
		if err := rows.Scan(
			&item.Id,
			&item.InternalId,
			&item.WorkspaceId,
			&item.Title,
			&item.Filename,
			&item.ContentType,
			&item.Checksum,
			&item.Size,
			&item.CharCount,
			&item.Labels,
			&item.ProcessedAt,
			&item.ChunkingConfig,
			&item.Meta,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// CountByWorkspaceId counts workspace document rows by workspace id.
func (r *DocumentRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM documents
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// SumSizeByWorkspaceId sums workspace document size by workspace id.
func (r *DocumentRepositoryPostgres) SumSizeByWorkspaceId(workspaceId Id) (int64, error) {
	var total int64
	err := r.db.QueryRow(
		`SELECT COALESCE(SUM(size), 0)
		FROM documents
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&total)
	return total, err
}

// NewDocumentMetaRepository returns the repository for workspace document meta.
func NewDocumentMetaRepository(db *sql.DB) DocumentMetaRepository {
	return &DocumentMetaRepositoryPostgres{db: db}
}

// Create inserts a workspace document metadata row.
func (r *DocumentMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO documents_meta (
			id, document_id, key, value
		)
		VALUES ($1, $2, $3, $4)`,
		metaId.String(), id.String(), key, value,
	)
	return err
}

// Get returns workspace document metadata by key.
func (r *DocumentMetaRepositoryPostgres) Get(id Id, key string) (*DocumentMeta, error) {
	meta := &DocumentMeta{}
	err := r.db.QueryRow(
		`SELECT
			id, document_id, key, value, created_at, updated_at
		FROM documents_meta
		WHERE document_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(
		&meta.Id,
		&meta.DocumentId,
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

// Update updates an existing workspace document metadata row.
func (r *DocumentMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE documents_meta
		SET
			value = $1,
			updated_at = $2
		WHERE document_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)
	return err
}

// Delete deletes a workspace document metadata row.
func (r *DocumentMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM documents_meta
		WHERE document_id = $1 AND key = $2`,
		id.String(), key,
	)
	return err
}

// ListByDocumentId lists workspace document metadata rows by workspace document id.
func (r *DocumentMetaRepositoryPostgres) ListByDocumentId(id Id) ([]*DocumentMeta, error) {
	rows, err := r.db.Query(
		`SELECT
			id, document_id, key, value, created_at, updated_at
		FROM documents_meta
		WHERE document_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*DocumentMeta
	for rows.Next() {
		meta := &DocumentMeta{}
		err := rows.Scan(
			&meta.Id,
			&meta.DocumentId,
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

// Upsert creates or updates workspace document metadata.
func (r *DocumentMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}
	return r.Update(id, key, value)
}
