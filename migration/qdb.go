// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package migration

import (
	"context"

	"github.com/actx0/ziee/pkg/qdrant"
)

const (
	WorkspaceKnowledgeCollection = "workspace_knowledge"

	PayloadWorkspaceID        = "_ws_id"
	PayloadDocumentInternalID = "_doc_int_id"
	PayloadText               = "_text"
	PayloadChunkIndex         = "_chk_idx"
	PayloadTag                = "tag"
)

// EnsureCollections creates the default Qdrant collections when they do not exist.
func EnsureCollections(ctx context.Context, client *qdrant.Client) error {
	return client.EnsureCollection(
		ctx,
		WorkspaceKnowledgeCollection,
		qdrant.Index{Field: PayloadWorkspaceID, Type: qdrant.IndexTypeUUID},
		qdrant.Index{Field: PayloadDocumentInternalID, Type: qdrant.IndexTypeUUID},
		qdrant.Index{Field: PayloadTag, Type: qdrant.IndexTypeString},
	)
}
