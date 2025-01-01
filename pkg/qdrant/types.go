// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package qdrant

const (
	IndexTypeString IndexType = "string"
	IndexTypeUUID   IndexType = "uuid"
)

// Point is a vector record stored in Qdrant.
type Point struct {
	Id      string
	Vector  []float32
	Payload map[string]string
}

// Result is a scored vector match from Qdrant.
type Result struct {
	Id      string
	Score   float32
	Payload map[string]string
}

// IndexType is the payload field type for a collection index.
type IndexType string

// Index is a payload field index on a collection.
type Index struct {
	Field string
	Type  IndexType
}

// Query is a vector search request.
type Query struct {
	Vector  []float32
	Filters map[string]string
	Limit   uint64
}
