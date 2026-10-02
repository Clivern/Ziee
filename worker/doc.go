// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/broker"
)

// HandleDocumentIndex indexes a document.
func (h *handlers) HandleDocumentIndex(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	err = h.Knowledge.Index(ctx, db.Id(payload["documentId"]))
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	return h.Tasks.Complete(taskId, "")
}

// HandleDocumentDelete deletes a document.
func (h *handlers) HandleDocumentDelete(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	err = h.Knowledge.Delete(ctx, db.Id(payload["documentId"]), payload["internalId"])
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	return h.Tasks.Complete(taskId, "")
}
