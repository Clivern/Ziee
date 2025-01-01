// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package knowledge

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/migration"
	"github.com/actx0/ziee/pkg/chunk"
	"github.com/actx0/ziee/pkg/util"

	"github.com/samber/lo"
)

const (
	StatusProcessing = "processing"
	StatusIndexed    = "indexed"
	StatusFailed     = "failed"
)

var (
	ErrNotImplemented   = errors.New("not implemented")
	ErrDocumentNotFound = errors.New("knowledge document not found")
	ErrIndexFailed      = errors.New("knowledge index failed")
	ErrDeleteFailed     = errors.New("knowledge delete failed")
	ErrSearchFailed     = errors.New("knowledge search failed")
)

// SearchRequest finds relevant knowledge chunks for a query.
type SearchRequest struct {
	WorkspaceId db.Id
	Query       string
	Limit       uint64
	Labels      map[string]string
}

// SearchHit is a scored knowledge chunk.
type SearchHit struct {
	DocumentId db.Id
	ChunkId    string
	Score      float32
	Text       string
	Labels     map[string]string
}

// Chunk is a text segment prepared for indexing.
type Chunk struct {
	Id   string
	Text string
}

// GetChunkConfig returns the server-side chunking configuration for a document.
func GetChunkConfig(document *db.WorkspaceDocument) chunk.Config {
	opts := util.DefaultChunkingOptions(document.CharCount, document.Filename)

	return chunk.Config{
		Splitter: chunk.Splitter(opts.Strategy),
		Size:     opts.Size,
		Overlap:  opts.Overlap,
	}
}

// GetDocumentPayload returns the payload for a workspace document chunk.
func GetDocumentPayload(document *db.WorkspaceDocument, text string, chunkIndex int) map[string]string {
	payload := map[string]string{
		migration.PayloadText:               text,
		migration.PayloadDocumentInternalID: document.InternalId.String(),
		migration.PayloadWorkspaceID:        document.WorkspaceId.String(),
		migration.PayloadChunkIndex:         strconv.Itoa(chunkIndex),
	}

	if tag := documentTag(document.Labels); lo.IsNotEmpty(tag) {
		payload[migration.PayloadTag] = tag
	}

	return payload
}

// documentTag builds the Qdrant tag for a workspace document.
func documentTag(labels *string) string {
	if labels == nil || lo.IsEmpty(*labels) {
		return ""
	}

	var items []string
	if err := json.Unmarshal([]byte(*labels), &items); err != nil {
		return ""
	}

	for _, item := range items {
		key, value, ok := strings.Cut(item, "=")
		if ok && key == migration.PayloadTag {
			return value
		}
	}

	return ""
}
