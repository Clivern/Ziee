// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package qdrant

import (
	"fmt"

	qdrantsdk "github.com/qdrant/go-client/qdrant"
)

// ParseFieldType parses a Qdrant field type from a string.
func ParseFieldType(fieldType IndexType) *qdrantsdk.FieldType {
	switch fieldType {
	case IndexTypeUUID:
		return qdrantsdk.FieldType_FieldTypeUuid.Enum()
	default:
		return qdrantsdk.FieldType_FieldTypeKeyword.Enum()
	}
}

// PointIdString formats a Qdrant point id as a string.
func PointIdString(id *qdrantsdk.PointId) string {
	switch value := id.PointIdOptions.(type) {
	case *qdrantsdk.PointId_Uuid:
		return value.Uuid
	case *qdrantsdk.PointId_Num:
		return fmt.Sprintf("%d", value.Num)
	default:
		return ""
	}
}

// Float32Vector converts embedding values to Qdrant's float32 vector format.
func Float32Vector(values []float64) []float32 {
	out := make([]float32, len(values))
	for i, value := range values {
		out[i] = float32(value)
	}
	return out
}
