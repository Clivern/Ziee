// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package event

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueComment(t *testing.T) {
	dat, err := os.ReadFile("../fixtures/issue_comment.json")
	require.NoError(t, err)

	var issueComment IssueComment
	ok, err := issueComment.LoadFromJSON(dat)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "Clivern", issueComment.Issue.User.Login)
}
