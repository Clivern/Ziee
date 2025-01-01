// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package knowledge

import (
	"context"
	"fmt"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/migration"
	"github.com/actx0/ziee/pkg/ai"
	"github.com/actx0/ziee/pkg/qdrant"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

// Search returns knowledge chunks relevant to the query.
func (s *Service) Search(ctx context.Context, req SearchRequest) ([]SearchHit, ai.Usage, error) {
	vector, usage, err := s.embed.EmbedQuery(ctx, req.Query, ai.EmbedOptions{
		Dimensions: int64(viper.GetInt("app.ai.embed.qdb.vector_size")),
	})
	if err != nil {
		return nil, usage, fmt.Errorf("%w: embed query: %v", ErrSearchFailed, err)
	}

	filters := map[string]string{
		migration.PayloadWorkspaceID: req.WorkspaceId.String(),
	}
	if tag, ok := req.Labels[migration.PayloadTag]; ok && lo.IsNotEmpty(tag) {
		filters[migration.PayloadTag] = tag
	}

	results, err := s.vectors.Search(ctx, migration.WorkspaceKnowledgeCollection, qdrant.Query{
		Vector:  qdrant.Float32Vector(vector),
		Filters: filters,
		Limit:   req.Limit,
	})
	if err != nil {
		return nil, usage, fmt.Errorf("%w: vector search: %v", ErrSearchFailed, err)
	}

	hits := make([]SearchHit, 0, len(results))
	for _, result := range results {
		internalIdStr := result.Payload[migration.PayloadDocumentInternalID]
		if lo.IsEmpty(internalIdStr) {
			continue
		}

		document, err := s.documents.GetByInternalId(db.Id(internalIdStr))
		if err != nil {
			log.Warn().
				Err(err).
				Str("internal_id", internalIdStr).
				Msg("Failed to resolve document for search hit")
			continue
		}
		if document == nil {
			continue
		}

		labels := map[string]string{}
		if tag := result.Payload[migration.PayloadTag]; lo.IsNotEmpty(tag) {
			labels[migration.PayloadTag] = tag
		}

		hits = append(hits, SearchHit{
			DocumentId: document.Id,
			ChunkId:    result.Id,
			Score:      result.Score,
			Text:       result.Payload[migration.PayloadText],
			Labels:     labels,
		})
	}

	log.Info().
		Str("workspace_id", req.WorkspaceId.String()).
		Int("hits", len(hits)).
		Int64("prompt_tokens", usage.PromptTokens).
		Int64("total_tokens", usage.TotalTokens).
		Int64("total_cost", usage.Cost).
		Msg("Document search completed")

	return hits, usage, nil
}
