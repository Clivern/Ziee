// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package v1

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitParseRepoFile(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	data, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", ".ziee.yml"))
	require.NoError(t, err)

	file, err := Parse(data)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0", file.Version)
	assert.False(t, file.MergeQueue.Enabled)
	assert.Equal(t, "serial", file.MergeQueue.Mode)
	assert.Equal(t, 5, file.MergeQueue.MaxParallelChecks)
	assert.Equal(t, "always", file.MergeQueue.ResetOnExternalMerge)
	assert.Equal(t, "outcomes", file.MergeQueue.Comments)
	assert.Equal(t, "state/queued", file.MergeQueue.Labels.Queued)
	assert.Equal(t, "state/checking", file.MergeQueue.Labels.Checking)
	assert.Equal(t, "state/dequeued", file.MergeQueue.Labels.Dequeued)

	assert.True(t, file.MergeQueue.PRTriage.AI.Enabled)
	assert.Equal(t, "platform_monorepo", file.MergeQueue.PRTriage.AI.KBTags[0].Tag)

	area := file.MergeQueue.PRTriage.Rules[0]
	assert.Equal(t, "area-api", area.Name)
	assert.Equal(t, []string{"api/**", "pkg/**"}, area.When[0].Files)
	assert.Equal(t, []string{"area/api"}, area.Labels.Add)
	assert.Equal(t, []string{"clivern"}, area.Assign)
	assert.Equal(t, []string{"maya"}, area.Reviewers)
	assert.Equal(t, []string{"core"}, area.ReviewTeams)

	sizeS := file.MergeQueue.PRTriage.Rules[2]
	assert.Equal(t, 20, *sizeS.When[0].MaxFilesChanged)

	sizeL := file.MergeQueue.PRTriage.Rules[3]
	assert.Equal(t, 21, *sizeL.When[0].MinFilesChanged)
	assert.Equal(t, []string{"size/s"}, sizeL.Labels.Remove)

	bug := file.MergeQueue.PRTriage.Rules[7]
	assert.Equal(t, []string{"dependabot", "ziee-bot"}, bug.When[0].AuthorNotIn)
	assert.Equal(t, "bug", bug.When[1].Intention)

	fromSre := file.MergeQueue.PRTriage.Rules[5]
	assert.Equal(t, []string{"sre"}, fromSre.When[0].AuthorInTeam)

	hotfixTitle := file.MergeQueue.PRTriage.Rules[6]
	assert.Equal(t, "[Hh]otfix", hotfixTitle.When[0].Title)

	queue := file.MergeQueue.Commands["queue"]
	assert.Equal(t, "write", queue.Allow[0].Permission)
	assert.Equal(t, []string{"sre", "core"}, queue.Allow[1].Teams)
	assert.Equal(t, []string{"clivern"}, queue.Allow[2].Users)

	assert.Equal(t, "hotfix", file.MergeQueue.PriorityRules[0].Name)
	assert.Equal(t, "hotfix", file.MergeQueue.PriorityRules[0].When[0].Label)
	assert.True(t, file.MergeQueue.PriorityRules[0].InterruptChecks)

	hotfixQ := file.MergeQueue.QueueRules[0]
	assert.Equal(t, 1, *hotfixQ.BatchSize.Value)
	assert.Equal(t, 1, hotfixQ.MergeWhen[0].Approvals.Min)
	assert.Equal(t, []string{"clivern"}, hotfixQ.MergeWhen[0].Approvals.From)
	assert.Equal(t, "ci", hotfixQ.MergeWhen[1].Check)
	assert.Equal(t, []string{"sre"}, hotfixQ.Allow[0].Teams)

	defaultQ := file.MergeQueue.QueueRules[1]
	assert.Equal(t, 1, defaultQ.BatchSize.Min)
	assert.Equal(t, 8, defaultQ.BatchSize.Max)
	assert.Equal(t, "5m", defaultQ.BatchMaxWait)
	assert.Equal(t, "ci-full", defaultQ.MergeWhen[0].Check)

	assert.False(t, file.PRReviews.Enabled)
	assert.True(t, file.IssueTriage.Enabled)
	assert.Equal(t, "outcomes", file.IssueTriage.Comments)
	assert.Equal(t, "platform_monorepo", file.IssueTriage.AI.KBTags[0].Tag)
	assert.Equal(t, []string{"clivern"}, file.IssueTriage.Commands["close"].Allow[1].Users)
	assert.Equal(t, "bug", file.IssueTriage.Rules[3].When[1].Intention)
}

func TestUnitParseWhenAllowBatchSize(t *testing.T) {
	file, err := Parse([]byte(`
version: 1.0.0
merge_queue:
  enabled: true
  pr_triage:
    rules:
      - name: size
        when:
          - max_files_changed: 3
          - min_files_changed: 1
          - author_not_in: [bot]
          - body: "repro"
        labels:
          add: [size/s]
  queue_rules:
    - name: default
      batch_size:
        min: 2
        max: 4
issue_triage:
  enabled: true
  commands:
    label:
      allow:
        - users: [maya]
`))
	require.NoError(t, err)

	when := file.MergeQueue.PRTriage.Rules[0].When
	assert.Equal(t, 3, *when[0].MaxFilesChanged)
	assert.Equal(t, 1, *when[1].MinFilesChanged)
	assert.Equal(t, []string{"bot"}, when[2].AuthorNotIn)
	assert.Equal(t, "repro", when[3].Body)
	assert.Equal(t, 2, file.MergeQueue.QueueRules[0].BatchSize.Min)
	assert.Equal(t, 4, file.MergeQueue.QueueRules[0].BatchSize.Max)
	assert.Equal(t, []string{"maya"}, file.IssueTriage.Commands["label"].Allow[0].Users)
}
