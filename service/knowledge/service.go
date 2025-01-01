// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.
//
// Knowledge service: document indexing and retrieval.
//
//	Index  → storage.GetDocument → chunk → ai.EmbedClient → qdrant.Upsert
//	Search → ai.EmbedClient.EmbedQuery → qdrant.Search
//	Delete → qdrant.Delete

package knowledge

import (
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/ai"
	"github.com/clivern/actx0/pkg/qdrant"
	"github.com/clivern/actx0/pkg/storage"
)

// Dependencies are the collaborators required by the knowledge service.
type Dependencies struct {
	Embed         *ai.EmbedClient
	Vectors       *qdrant.Client
	Documents     db.WorkspaceDocumentRepository
	Store         storage.Store
	Usage         db.UsageRepository
	Subscriptions db.SubscriptionRepository
}

// Service handles agent knowledge indexing and retrieval.
type Service struct {
	embed         *ai.EmbedClient
	vectors       *qdrant.Client
	documents     db.WorkspaceDocumentRepository
	store         storage.Store
	usage         db.UsageRepository
	subscriptions db.SubscriptionRepository
}

// New returns a knowledge service.
func New(deps Dependencies) *Service {
	return &Service{
		embed:         deps.Embed,
		vectors:       deps.Vectors,
		documents:     deps.Documents,
		store:         deps.Store,
		usage:         deps.Usage,
		subscriptions: deps.Subscriptions,
	}
}

// MarkAsFailed marks a workspace document as failed
func (s *Service) MarkAsFailed(document *db.WorkspaceDocument) {
	document.Status = db.WorkspaceDocumentStatusFailed
	now := time.Now().UTC()
	document.ProcessedAt = &now

	s.documents.Update(document)
}
