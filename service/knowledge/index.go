// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package knowledge

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/migration"
	"github.com/clivern/actx0/pkg/ai"
	"github.com/clivern/actx0/pkg/chunk"
	"github.com/clivern/actx0/pkg/qdrant"
	"github.com/clivern/actx0/pkg/storage"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Index chunks, embeds, and stores a workspace document.
func (s *Service) Index(ctx context.Context, documentId db.Id) error {
	var promptTokens, totalTokens, totalCost int64

	document, err := s.documents.GetById(documentId)
	if err != nil {
		return err
	}
	if document == nil {
		return ErrDocumentNotFound
	}

	text, err := s.store.GetDocument(ctx, document.InternalId.String())
	if err != nil {
		if errors.Is(err, storage.ErrDocumentNotFound) {
			s.MarkAsFailed(document)
			return ErrDocumentNotFound
		}
		s.MarkAsFailed(document)
		return fmt.Errorf("%w: load document: %v", ErrIndexFailed, err)
	}

	chunker, err := chunk.New(GetChunkConfig(document))
	if err != nil {
		s.MarkAsFailed(document)
		return fmt.Errorf("%w: chunker: %v", ErrIndexFailed, err)
	}

	parts, err := chunker.Split(text)
	if err != nil {
		s.MarkAsFailed(document)
		return fmt.Errorf("%w: split text: %v", ErrIndexFailed, err)
	}

	if len(parts) == 0 {
		s.MarkAsFailed(document)
		return fmt.Errorf("%w: no parts found", ErrIndexFailed)
	}

	points := make([]qdrant.Point, 0, len(parts))
	for i, part := range parts {
		embeddings, cu, err := s.embed.EmbedDocuments(ctx, []string{part}, ai.EmbedOptions{
			Dimensions: int64(viper.GetInt("app.ai.embed.qdb.vector_size")),
		})
		if err != nil {
			s.MarkAsFailed(document)
			return fmt.Errorf("%w: embed chunk %d: %v", ErrIndexFailed, i, err)
		}
		if len(embeddings) != 1 {
			s.MarkAsFailed(document)
			return fmt.Errorf("%w: embed chunk %d: expected 1 embedding, got %d", ErrIndexFailed, i, len(embeddings))
		}

		chunkId, err := db.NewId()
		if err != nil {
			s.MarkAsFailed(document)
			return fmt.Errorf("%w: chunk id: %v", ErrIndexFailed, err)
		}

		points = append(points, qdrant.Point{
			Id:      chunkId.String(),
			Vector:  qdrant.Float32Vector(embeddings[0]),
			Payload: GetDocumentPayload(document, part, i),
		})

		promptTokens += cu.PromptTokens
		totalTokens += cu.TotalTokens
		totalCost += cu.Cost
	}

	sub, err := s.subscriptions.GetByWorkspaceId(document.WorkspaceId)
	if err != nil {
		log.Warn().
			Err(err).
			Str("workspace_id", document.WorkspaceId.String()).
			Int64("total_tokens", totalTokens).
			Int64("total_cost", totalCost).
			Msg("Failed to record AI usage")
	}

	if sub != nil && sub.CurrentPeriodStart != nil && sub.CurrentPeriodEnd != nil {
		// Increment AI tokens usage
		s.usage.IncrementByPeriod(
			document.WorkspaceId,
			db.UsageTypeAITokens,
			sub.CurrentPeriodStart.UTC(),
			sub.CurrentPeriodEnd.UTC(),
			totalTokens,
			db.UsageUnitTokens,
		)

		// Increment AI cost usage
		s.usage.IncrementByPeriod(
			document.WorkspaceId,
			db.UsageTypeAICost,
			sub.CurrentPeriodStart.UTC(),
			sub.CurrentPeriodEnd.UTC(),
			totalCost,
			db.UsageUnitNanoUSD,
		)
	}

	log.Info().
		Str("document_id", documentId.String()).
		Int("chunks", len(parts)).
		Int64("prompt_tokens", promptTokens).
		Int64("total_tokens", totalTokens).
		Int64("total_cost", totalCost).
		Msg("Document embedded")

	err = s.vectors.DeleteByFilter(ctx, migration.WorkspaceKnowledgeCollection, map[string]string{
		migration.PayloadDocumentInternalID: document.InternalId.String(),
	})
	if err != nil {
		s.MarkAsFailed(document)
		return fmt.Errorf("%w: delete existing vectors: %v", ErrIndexFailed, err)
	}

	err = s.vectors.Upsert(
		ctx,
		migration.WorkspaceKnowledgeCollection,
		points,
	)

	if err != nil {
		s.MarkAsFailed(document)
		return fmt.Errorf("%w: upsert vectors: %v", ErrIndexFailed, err)
	}

	document.Status = db.WorkspaceDocumentStatusIndexed
	document.ProcessedAt = new(time.Now().UTC())

	return s.documents.Update(document)
}
