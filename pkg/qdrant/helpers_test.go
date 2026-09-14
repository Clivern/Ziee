// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package qdrant

import (
	"testing"

	qdrantsdk "github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/assert"
)

func TestUnitHelpers(t *testing.T) {
	t.Run("ParseFieldType", func(t *testing.T) {
		assert.Equal(t, qdrantsdk.FieldType_FieldTypeUuid.Enum(), ParseFieldType(IndexTypeUUID))
		assert.Equal(t, qdrantsdk.FieldType_FieldTypeKeyword.Enum(), ParseFieldType(IndexTypeString))
	})

	t.Run("PointIdString", func(t *testing.T) {
		uuid := &qdrantsdk.PointId{PointIdOptions: &qdrantsdk.PointId_Uuid{Uuid: "abc"}}
		assert.Equal(t, "abc", PointIdString(uuid))

		num := &qdrantsdk.PointId{PointIdOptions: &qdrantsdk.PointId_Num{Num: 42}}
		assert.Equal(t, "42", PointIdString(num))

		assert.Equal(t, "", PointIdString(&qdrantsdk.PointId{}))
	})

	t.Run("Float32Vector", func(t *testing.T) {
		assert.Equal(t, []float32{1.5, 2.25}, Float32Vector([]float64{1.5, 2.25}))
		assert.Empty(t, Float32Vector(nil))
	})
}
