// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package knowledge

import (
	"context"
	"errors"
	"fmt"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/migration"
	"github.com/actx0/ziee/pkg/storage"

	"github.com/rs/zerolog/log"
)

// Delete removes a workspace document from storage and the database
func (s *Service) Delete(ctx context.Context, documentId db.Id, internalId string) error {
	err := s.store.DeleteDocument(ctx, internalId)
	if err != nil && !errors.Is(err, storage.ErrDocumentNotFound) {
		return fmt.Errorf("%w: delete storage: %v", ErrDeleteFailed, err)
	}

	err = s.documents.Delete(documentId)
	if err != nil {
		return fmt.Errorf("%w: delete document: %v", ErrDeleteFailed, err)
	}

	err = s.vectors.DeleteByFilter(ctx, migration.WorkspaceKnowledgeCollection, map[string]string{
		migration.PayloadDocumentInternalID: internalId,
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("document_id", documentId.String()).
			Str("internal_id", internalId).
			Msg("Failed to delete document vectors")
		return fmt.Errorf("%w: delete vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}
