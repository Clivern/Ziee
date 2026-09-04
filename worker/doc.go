// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/nats"
)

var ErrInvalidPayload = errors.New("invalid worker payload")

// Knowledge indexes and deletes workspace documents.
type Knowledge interface {
	Index(ctx context.Context, documentId db.Id) error
	Delete(ctx context.Context, documentId db.Id, internalId string) error
}

// Dependencies are the services required by worker handlers.
type Dependencies struct {
	Knowledge Knowledge
}

type handlers struct {
	knowledge Knowledge
}

// Register attaches all worker handlers.
func Register(deps Dependencies) {
	h := &handlers{knowledge: deps.Knowledge}

	On(conf.NATSSubjectDocIndex, h.HandleDocumentIndex)
	On(conf.NATSSubjectDocDelete, h.HandleDocumentDelete)
}

// HandleDocumentIndex indexes a document.
func (h *handlers) HandleDocumentIndex(ctx context.Context, msg *nats.Msg) error {
	var payload map[string]string

	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	return h.knowledge.Index(ctx, db.Id(payload["documentId"]))
}

// HandleDocumentDelete deletes a document.
func (h *handlers) HandleDocumentDelete(ctx context.Context, msg *nats.Msg) error {
	var payload map[string]string

	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	return h.knowledge.Delete(ctx, db.Id(payload["documentId"]), payload["internalId"])
}
