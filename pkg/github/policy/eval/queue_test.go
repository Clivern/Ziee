// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"testing"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitEvaluateMergeQueueCommentQueue(t *testing.T) {
	viper.Set("app.oauth.github.bot_name", "zieeai")

	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			Enabled:  true,
			Comments: policy.CommentsOutcomes,
			Labels: v1.QueueLabels{
				Queued:   "state/queued",
				Checking: "state/checking",
				Dequeued: "state/dequeued",
			},
			Commands: v1.Commands{
				"queue": {Allow: v1.Allow{{Permission: "write"}}},
			},
			PriorityRules: []v1.PriorityRule{
				{Name: "hotfix", When: v1.Clauses{{Label: "hotfix"}}, Priority: db.PQueuePriorityHigh},
				{Name: "normal", When: v1.Clauses{}, Priority: db.PQueuePriorityMedium},
			},
			QueueRules: []v1.QueueRule{
				{
					Name:      "hotfix",
					Allow:     v1.Allow{{Teams: []string{"sre"}}},
					QueueWhen: v1.Clauses{{Label: "hotfix"}},
				},
				{
					Name: "default",
					QueueWhen: v1.Clauses{
						{Approvals: &v1.Approvals{Min: 1}},
						{Check: "ci"},
						{Label: "state/queued"},
					},
				},
			},
		},
	}

	plan := EvaluatePRComment(conf, Event{
		Comment: "@zieeai queue",
		Issue:   Issue{Number: 12, Author: "maya", Labels: []string{"bug"}},
		Actor:   Actor{Login: "clivern"},
		Account: Account{Login: "clivern", Type: "User"},
	}, &stubClient{permission: "write"})

	assert.Equal(t, []action.Action{
		{Kind: policy.Queue, Priority: db.PQueuePriorityMedium, Rule: "default"},
		{Kind: policy.AddLabels, Labels: []string{"state/queued"}},
		{Kind: policy.RemoveLabels, Labels: []string{"state/dequeued", "state/checking"}},
		{Kind: policy.Comment, Body: "Queued this pull request as requested by @clivern."},
	}, plan.Actions)
}

func TestUnitEvaluateMergeQueueCommentQueueHotfix(t *testing.T) {
	viper.Set("app.oauth.github.bot_name", "zieeai")

	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			Enabled:  true,
			Comments: policy.CommentsOutcomes,
			Labels: v1.QueueLabels{
				Queued:   "state/queued",
				Checking: "state/checking",
				Dequeued: "state/dequeued",
			},
			Commands: v1.Commands{
				"queue": {Allow: v1.Allow{{Permission: "write"}}},
			},
			PriorityRules: []v1.PriorityRule{
				{Name: "hotfix", When: v1.Clauses{{Label: "hotfix"}}, Priority: db.PQueuePriorityHigh},
				{Name: "normal", When: v1.Clauses{}, Priority: db.PQueuePriorityMedium},
			},
			QueueRules: []v1.QueueRule{
				{
					Name:      "hotfix",
					Allow:     v1.Allow{{Teams: []string{"sre"}}},
					QueueWhen: v1.Clauses{{Label: "hotfix"}},
				},
				{Name: "default"},
			},
		},
	}

	plan := EvaluatePRComment(conf, Event{
		Comment: "@zieeai queue hotfix",
		Issue:   Issue{Number: 12, Author: "maya", Labels: []string{"hotfix"}},
		Actor:   Actor{Login: "clivern"},
	}, &stubClient{permission: "write", teams: []string{"sre"}})

	assert.Equal(t, []action.Action{
		{Kind: policy.Queue, Priority: db.PQueuePriorityHigh, Rule: "hotfix"},
		{Kind: policy.AddLabels, Labels: []string{"state/queued"}},
		{Kind: policy.RemoveLabels, Labels: []string{"state/dequeued", "state/checking"}},
		{Kind: policy.Comment, Body: "Queued this pull request as requested by @clivern."},
	}, plan.Actions)
}

func TestUnitEvaluateMergeQueueCommentDequeue(t *testing.T) {
	viper.Set("app.oauth.github.bot_name", "zieeai")

	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			Enabled:  true,
			Comments: policy.CommentsOutcomes,
			Labels: v1.QueueLabels{
				Queued:   "state/queued",
				Checking: "state/checking",
				Dequeued: "state/dequeued",
			},
			Commands: v1.Commands{
				"dequeue": {Allow: v1.Allow{{Permission: "write"}}},
			},
		},
	}

	plan := EvaluatePRComment(conf, Event{
		Comment: "@zieeai dequeue",
		Issue:   Issue{Number: 12, Author: "maya", Labels: []string{"state/queued"}},
		Actor:   Actor{Login: "clivern"},
	}, &stubClient{permission: "write"})

	assert.Equal(t, []action.Action{
		{Kind: policy.Dequeue},
		{Kind: policy.AddLabels, Labels: []string{"state/dequeued"}},
		{Kind: policy.RemoveLabels, Labels: []string{"state/queued", "state/checking"}},
		{Kind: policy.Comment, Body: "Dequeued this pull request as requested by @clivern."},
	}, plan.Actions)
}

func TestUnitEvaluateMergeQueueCommentRequeue(t *testing.T) {
	viper.Set("app.oauth.github.bot_name", "zieeai")

	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			Enabled:  true,
			Comments: policy.CommentsOutcomes,
			Labels: v1.QueueLabels{
				Queued:   "state/queued",
				Checking: "state/checking",
				Dequeued: "state/dequeued",
			},
			Commands: v1.Commands{
				"requeue": {Allow: v1.Allow{{Permission: "write"}}},
			},
			PriorityRules: []v1.PriorityRule{
				{Name: "normal", When: v1.Clauses{}, Priority: db.PQueuePriorityMedium},
			},
			QueueRules: []v1.QueueRule{
				{Name: "default"},
			},
		},
	}

	plan := EvaluatePRComment(conf, Event{
		Comment: "@zieeai requeue",
		Issue:   Issue{Number: 12, Author: "maya", Labels: []string{"state/queued"}},
		Actor:   Actor{Login: "clivern"},
	}, &stubClient{permission: "write"})

	assert.Equal(t, []action.Action{
		{Kind: policy.Dequeue},
		{Kind: policy.AddLabels, Labels: []string{"state/dequeued"}},
		{Kind: policy.RemoveLabels, Labels: []string{"state/queued", "state/checking"}},
		{Kind: policy.Queue, Priority: db.PQueuePriorityMedium, Rule: "default"},
		{Kind: policy.AddLabels, Labels: []string{"state/queued"}},
		{Kind: policy.RemoveLabels, Labels: []string{"state/dequeued", "state/checking"}},
		{Kind: policy.Comment, Body: "Requeued this pull request as requested by @clivern."},
	}, plan.Actions)
}

func TestUnitEvaluateMergeQueueDisabled(t *testing.T) {
	viper.Set("app.oauth.github.bot_name", "zieeai")

	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			Enabled: false,
			Commands: v1.Commands{
				"queue": {Allow: v1.Allow{{Permission: "write"}}},
			},
			QueueRules: []v1.QueueRule{{Name: "default"}},
		},
	}

	plan := EvaluatePRComment(conf, Event{
		Comment: "@zieeai queue",
		Issue:   Issue{Number: 12},
		Actor:   Actor{Login: "clivern"},
	}, &stubClient{permission: "write"})

	assert.Empty(t, plan.Actions)
}

func TestUnitQueueWhenForEntry(t *testing.T) {
	when := v1.Clauses{
		{Approvals: &v1.Approvals{Min: 1}},
		{Label: "state/queued"},
		{Check: "ci"},
	}

	got := QueueWhenForEntry(when, "state/queued")

	assert.Equal(t, v1.Clauses{
		{Approvals: &v1.Approvals{Min: 1}},
		{Check: "ci"},
	}, got)
}
