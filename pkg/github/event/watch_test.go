// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package event

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatch(t *testing.T) {
	dat, err := os.ReadFile("../fixtures/watch.json")
	require.NoError(t, err)

	var watch Watch
	ok, err := watch.LoadFromJSON(dat)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "Clivern", watch.Repository.Owner.Login)
}
