// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package event

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	dat, err := os.ReadFile("../fixtures/status.json")
	require.NoError(t, err)

	var status Status
	ok, err := status.LoadFromJSON(dat)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "Clivern", status.Commit.Commit.Author.Name)
}
