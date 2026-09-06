// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/storage"
	"github.com/clivern/ziee/pkg/util"
	"github.com/clivern/ziee/service/knowledge"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

var (
	ErrDocumentNotFound      = errors.New("document not found")
	ErrFailedCreateDocument  = errors.New("failed create document")
	ErrFailedStoreDocument   = errors.New("failed store document")
	ErrFailedListDocuments   = errors.New("failed list documents")
	ErrFailedDeleteDocument  = errors.New("failed delete document")
	ErrFailedSearchDocuments = errors.New("failed search documents")
)

// DocumentMeasure is a numeric document field with a unit for API responses.
type DocumentMeasure struct {
	Value int64  `json:"value"`
	Unit  string `json:"unit"`
}

// DocumentResponse is a workspace document shaped for API responses.
type DocumentResponse struct {
	Id               db.Id           `json:"id"`
	WorkspaceId      db.Id           `json:"workspaceId"`
	Title            string          `json:"title"`
	Filename         string          `json:"filename"`
	ContentType      string          `json:"contentType"`
	Checksum         string          `json:"checksum"`
	Size             DocumentMeasure `json:"size"`
	CharCount        int64           `json:"charCount"`
	Labels           []string        `json:"labels,omitempty"`
	ChunkingStrategy string          `json:"chunkingStrategy,omitempty"`
	ChunkSize        int             `json:"chunkSize,omitempty"`
	ChunkOverlap     int             `json:"chunkOverlap,omitempty"`
	Status           string          `json:"status"`
	CreatedAt        string          `json:"createdAt"`
	UpdatedAt        string          `json:"updatedAt"`
}

// ListDocumentsResponse is returned when listing documents.
type ListDocumentsResponse struct {
	Documents []*DocumentResponse
	Total     int64
}

// SearchDocumentsRequest is the body for searching workspace documents.
type SearchDocumentsRequest struct {
	Query  string            `json:"query" validate:"required,min=1" label:"Query"`
	Labels map[string]string `json:"labels" validate:"omitempty" label:"Labels"`
	Limit  int               `json:"limit" validate:"omitempty,min=1,max=100" label:"Limit"`
}

// SearchDocumentHit is a scored document chunk for API responses.
type SearchDocumentHit struct {
	DocumentId db.Id             `json:"documentId"`
	ChunkId    string            `json:"chunkId"`
	Score      float32           `json:"score"`
	Text       string            `json:"text"`
	Labels     map[string]string `json:"labels,omitempty"`
}

// SearchDocumentsResponse is returned when searching documents.
type SearchDocumentsResponse struct {
	Results []SearchDocumentHit `json:"results"`
}

// Document is the module for workspace document CRUD.
type Document struct {
	DocumentRepository     db.WorkspaceDocumentRepository
	WorkspaceRepository    db.WorkspaceRepository
	Store                  storage.Store
	Knowledge              *knowledge.Service
	UsageRepository        db.UsageRepository
	SubscriptionRepository db.SubscriptionRepository
}

// NewDocument creates a document module with the given repositories and store.
func NewDocument(documents db.WorkspaceDocumentRepository, workspaces db.WorkspaceRepository, store storage.Store) *Document {
	return &Document{
		DocumentRepository:  documents,
		WorkspaceRepository: workspaces,
		Store:               store,
	}
}

// UploadDocument stores document metadata in the database and text in storage.
func (d *Document) UploadDocument(ctx context.Context, form *util.UploadForm, workspaceId db.Id) (*DocumentResponse, error) {
	workspace, err := d.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	var labels *string
	if len(form.Labels) > 0 {
		raw, err := json.Marshal(form.Labels)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedCreateDocument, err)
		}
		labels = new(string(raw))
	}

	raw, err := json.Marshal(util.DefaultChunkingOptions(form.CharCount, form.Filename))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateDocument, err)
	}
	chunking := string(raw)

	document := &db.WorkspaceDocument{
		WorkspaceId:    workspaceId,
		Title:          form.Title,
		Filename:       form.Filename,
		ContentType:    form.ContentType,
		Checksum:       form.Checksum,
		Size:           form.Size,
		CharCount:      form.CharCount,
		Labels:         labels,
		ChunkingConfig: &chunking,
		Status:         db.WorkspaceDocumentStatusProcessing,
	}

	err = d.DocumentRepository.Create(document)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateDocument, err)
	}

	_, err = d.Store.PutDocument(ctx, document.InternalId.String(), form.Content)
	if err != nil {
		derr := d.DocumentRepository.Delete(document.Id)
		if derr != nil {
			log.Error().
				Err(derr).
				Str("documentId", document.Id.String()).
				Msg("Failed to rollback document after storage error")
		}
		return nil, fmt.Errorf("%w: %v", ErrFailedStoreDocument, err)
	}

	err = d.EnqueueIndexTask(workspaceId, document.Id)
	if err != nil {
		log.Error().
			Err(err).
			Str("documentId", document.Id.String()).
			Msg("Failed to enqueue document index task")
	}

	chkOpts := util.DefaultChunkingOptions(document.CharCount, document.Filename)

	return &DocumentResponse{
		Id:          document.Id,
		WorkspaceId: document.WorkspaceId,
		Title:       document.Title,
		Filename:    document.Filename,
		ContentType: document.ContentType,
		Checksum:    document.Checksum,
		Size: DocumentMeasure{
			Value: document.Size,
			Unit:  "bytes",
		},
		CharCount:        document.CharCount,
		Labels:           form.Labels,
		ChunkingStrategy: chkOpts.Strategy,
		ChunkSize:        chkOpts.Size,
		ChunkOverlap:     chkOpts.Overlap,
		Status:           document.Status,
		CreatedAt:        document.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:        document.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// DeleteDocument enqueues removal of a workspace document from storage, Qdrant, and the database.
func (d *Document) DeleteDocument(_ context.Context, workspaceId, documentId db.Id) error {
	workspace, err := d.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	document, err := d.DocumentRepository.GetById(documentId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteDocument, err)
	}
	if document == nil || document.WorkspaceId != workspaceId {
		return ErrDocumentNotFound
	}

	err = d.DocumentRepository.Delete(documentId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteDocument, err)
	}

	return d.EnqueueDeleteTask(workspaceId, documentId, document.InternalId)
}

// ListDocuments returns paginated documents for a workspace.
func (d *Document) ListDocuments(workspaceId db.Id, limit, offset int) (*ListDocumentsResponse, error) {
	workspace, err := d.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	documents, err := d.DocumentRepository.ListByWorkspaceId(workspaceId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListDocuments, err)
	}

	total, err := d.DocumentRepository.CountByWorkspaceId(workspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListDocuments, err)
	}

	responses := make([]*DocumentResponse, 0, len(documents))
	for _, document := range documents {
		var labels []string
		if document.Labels != nil && lo.IsNotEmpty(*document.Labels) {
			_ = json.Unmarshal([]byte(*document.Labels), &labels)
		}

		chkOpts := util.DefaultChunkingOptions(document.CharCount, document.Filename)

		responses = append(responses, &DocumentResponse{
			Id:          document.Id,
			WorkspaceId: document.WorkspaceId,
			Title:       document.Title,
			Filename:    document.Filename,
			ContentType: document.ContentType,
			Checksum:    document.Checksum,
			Size: DocumentMeasure{
				Value: document.Size,
				Unit:  "bytes",
			},
			CharCount:        document.CharCount,
			Labels:           labels,
			ChunkingStrategy: chkOpts.Strategy,
			ChunkSize:        chkOpts.Size,
			ChunkOverlap:     chkOpts.Overlap,
			Status:           document.Status,
			CreatedAt:        document.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:        document.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &ListDocumentsResponse{
		Documents: responses,
		Total:     total,
	}, nil
}

// SearchDocuments finds relevant document chunks for a semantic query.
func (d *Document) SearchDocuments(ctx context.Context, workspaceId db.Id, query string, labels map[string]string, limit int) (*SearchDocumentsResponse, error) {
	workspace, err := d.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	hits, usage, err := d.Knowledge.Search(ctx, knowledge.SearchRequest{
		WorkspaceId: workspaceId,
		Query:       query,
		Limit:       uint64(limit),
		Labels:      labels,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedSearchDocuments, err)
	}

	err = NewUsage().IncrementAIUsage(
		d.UsageRepository,
		d.SubscriptionRepository,
		workspaceId,
		usage.TotalTokens,
		usage.Cost,
	)
	if err != nil {
		log.Warn().
			Err(err).
			Str("workspaceId", workspaceId.String()).
			Int64("total_tokens", usage.TotalTokens).
			Int64("total_cost", usage.Cost).
			Msg("Failed to record AI usage for document search")
	}

	results := make([]SearchDocumentHit, 0, len(hits))
	for _, hit := range hits {
		results = append(results, SearchDocumentHit{
			DocumentId: hit.DocumentId,
			ChunkId:    hit.ChunkId,
			Score:      hit.Score,
			Text:       hit.Text,
			Labels:     hit.Labels,
		})
	}

	return &SearchDocumentsResponse{
		Results: results,
	}, nil
}

// EnqueueIndexTask records an async task and publishes a NATS message to index a document.
func (d *Document) EnqueueIndexTask(workspaceId, documentId db.Id) error {
	return d.enqueueTask(db.AsyncTaskTypeDocIndex, conf.NATSSubjectDocIndex, map[string]string{
		"workspaceId": workspaceId.String(),
		"documentId":  documentId.String(),
	}, workspaceId)
}

// EnqueueDeleteTask records an async task and publishes a NATS message to delete a document.
func (d *Document) EnqueueDeleteTask(workspaceId, documentId, internalId db.Id) error {
	return d.enqueueTask(db.AsyncTaskTypeDocDelete, conf.NATSSubjectDocDelete, map[string]string{
		"workspaceId": workspaceId.String(),
		"documentId":  documentId.String(),
		"internalId":  internalId.String(),
	}, workspaceId)
}

func (d *Document) enqueueTask(taskType, subject string, payload map[string]string, workspaceId db.Id) error {
	taskId, err := db.NewId()
	if err != nil {
		return err
	}

	payload["taskId"] = taskId.String()
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body := string(raw)

	err = db.NewAsyncTaskRepository(db.GetDB()).Create(&db.AsyncTask{
		Id:          taskId,
		WorkspaceId: workspaceId,
		Type:        taskType,
		Status:      db.AsyncTaskStatusPending,
		Payload:     &body,
	})
	if err != nil {
		return err
	}

	bus := GetBus()
	err = bus.Publish(subject, raw)
	if err != nil {
		return err
	}

	return bus.Flush()
}
