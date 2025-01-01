// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package task

import (
	"context"

	"github.com/actx0/ziee/db"
)

// HandleDocumentIndex handles the document index task.
func (h *Handlers) HandleDocumentIndex(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.knowledge.Index(ctx, db.Id(payload["documentId"]))
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "indexed",
	})
}

// HandleDocumentDelete removes a document from storage, Qdrant, and the database.
func (h *Handlers) HandleDocumentDelete(ctx context.Context, task *db.AsyncTask) (string, error) {
	payload, err := DecodePayload(task.Payload)
	if err != nil {
		return "", err
	}

	err = h.knowledge.Delete(ctx, db.Id(payload["documentId"]), payload["internalId"])
	if err != nil {
		return "", err
	}

	return EncodeResult(map[string]string{
		"status": "deleted",
	})
}
