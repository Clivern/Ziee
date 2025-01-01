// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package event

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssues(t *testing.T) {
	dat, err := os.ReadFile("../fixtures/issues.json")
	require.NoError(t, err)

	var issues Issues
	ok, err := issues.LoadFromJSON(dat)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "Clivern", issues.Issue.User.Login)
}
