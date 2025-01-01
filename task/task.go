// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package task

import (
	"context"
	"fmt"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/module"
)

// Knowledge indexes and deletes workspace documents.
type Knowledge interface {
	Index(ctx context.Context, documentId db.Id) error
	Delete(ctx context.Context, documentId db.Id, internalId string) error
}

// Dependencies are the services required by async task handlers.
type Dependencies struct {
	Knowledge Knowledge
}

// Handlers registers and runs async task handlers.
type Handlers struct {
	knowledge Knowledge
}

// New returns task handlers backed by services only.
func New(deps Dependencies) *Handlers {
	return &Handlers{
		knowledge: deps.Knowledge,
	}
}

// Register attaches all task handlers to the async manager.
func Register(am *module.Async, deps Dependencies) {
	New(deps).Register(am)
}

// Register attaches handlers to the async module.
func (h *Handlers) Register(am *module.Async) {
	am.RegisterHandler("task.doc.index", h.HandleDocumentIndex)
	am.RegisterHandler("task.doc.delete", h.HandleDocumentDelete)
}

// Run executes a handler synchronously without HTTP or the worker pool.
func (h *Handlers) Run(ctx context.Context, taskType string, task *db.AsyncTask) (string, error) {
	switch taskType {
	case "task.doc.index":
		return h.HandleDocumentIndex(ctx, task)
	case "task.doc.delete":
		return h.HandleDocumentDelete(ctx, task)
	default:
		return "", fmt.Errorf("unknown task type: %s", taskType)
	}
}
