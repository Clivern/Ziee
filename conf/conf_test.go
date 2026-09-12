// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package conf

import (
	"testing"

	"github.com/clivern/ziee/pkg/github/policy/spec"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitDefaultZieeYML(t *testing.T) {
	file, err := spec.Parse([]byte(DefaultZieeYML))
	require.NoError(t, err)

	assert.Equal(t, "1.0.0", file.Version)
	assert.Equal(t, "state/queued", file.Labels[0].Name)
	assert.Equal(t, "maintainers", file.Teams[0].Name)
	assert.Equal(t, "docs/", file.Knowledge[0].Path)
	assert.False(t, file.MergeQueue.Enabled)
	assert.Equal(t, "serial", file.MergeQueue.Mode)
	assert.Equal(t, "state/queued", file.MergeQueue.Labels.Queued)
	assert.False(t, file.MergeQueue.PRTriage.AI.Enabled)
	assert.Equal(t, "size-s", file.MergeQueue.PRTriage.Rules[0].Name)
	assert.Equal(t, "write", file.MergeQueue.Commands["queue"].Allow[0].Permission)
	assert.Equal(t, "hotfix", file.MergeQueue.PriorityRules[0].Name)
	assert.Equal(t, "default", file.MergeQueue.QueueRules[0].Name)
	assert.False(t, file.IssueTriage.Enabled)
	assert.False(t, file.PRReviews.Enabled)
}
