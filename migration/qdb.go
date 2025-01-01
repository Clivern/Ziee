// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package migration

import (
	"context"

	"github.com/clivern/actx0/pkg/qdrant"
)

const (
	WorkspaceKnowledgeCollection   = "workspace_knowledge"
	AgentSessionMessagesCollection = "agent_session_messages"
	AgentSessionMemoriesCollection = "agent_session_memories"

	PayloadWorkspaceID        = "_ws_id"
	PayloadDocumentInternalID = "_doc_int_id"
	PayloadText               = "_text"
	PayloadChunkIndex         = "_chk_idx"
	PayloadAgentID            = "_agent_id"
	PayloadSessionID          = "_session_id"
	PayloadMemoryID           = "_mem_id"
	PayloadMessageID          = "_msg_id"
	PayloadTag                = "tag"
)

// EnsureCollections creates the default Qdrant collections when they do not exist.
func EnsureCollections(ctx context.Context, client *qdrant.Client) error {
	collections := []struct {
		name    string
		indexes []qdrant.Index
	}{
		{
			name: WorkspaceKnowledgeCollection,
			indexes: []qdrant.Index{
				{Field: PayloadWorkspaceID, Type: qdrant.IndexTypeUUID},
				{Field: PayloadDocumentInternalID, Type: qdrant.IndexTypeUUID},
				{Field: PayloadTag, Type: qdrant.IndexTypeString},
			},
		},
		{
			name: AgentSessionMessagesCollection,
			indexes: []qdrant.Index{
				{Field: PayloadWorkspaceID, Type: qdrant.IndexTypeUUID},
				{Field: PayloadAgentID, Type: qdrant.IndexTypeUUID},
				{Field: PayloadSessionID, Type: qdrant.IndexTypeUUID},
				{Field: PayloadMessageID, Type: qdrant.IndexTypeUUID},
			},
		},
		{
			name: AgentSessionMemoriesCollection,
			indexes: []qdrant.Index{
				{Field: PayloadWorkspaceID, Type: qdrant.IndexTypeUUID},
				{Field: PayloadAgentID, Type: qdrant.IndexTypeUUID},
				{Field: PayloadSessionID, Type: qdrant.IndexTypeUUID},
				{Field: PayloadMemoryID, Type: qdrant.IndexTypeUUID},
			},
		},
	}

	for _, collection := range collections {
		err := client.EnsureCollection(
			ctx,
			collection.name,
			collection.indexes...,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
