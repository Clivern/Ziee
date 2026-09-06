// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"testing"

	"github.com/clivern/ziee/policy"
	"github.com/clivern/ziee/policy/action"
	v1 "github.com/clivern/ziee/policy/spec/v1"

	"github.com/stretchr/testify/assert"
)

func TestUnitEvaluateIssueOpened(t *testing.T) {
	file := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			Rules: []v1.Rule{
				{
					Name:   "hotfix",
					When:   v1.Clauses{{Title: "[Hh]otfix"}},
					Labels: v1.Labels{Add: []string{"hotfix"}},
					Assign: []string{"clivern"},
				},
				{
					Name:   "bot",
					When:   v1.Clauses{{AuthorIn: []string{"dependabot"}}},
					Labels: v1.Labels{Add: []string{"bot"}},
				},
			},
		},
	}

	plan := EvaluateIssueOpened(file, Event{
		Issue: Issue{
			Title:  "Hotfix: cache",
			Author: "maya",
		},
	}, nil)

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"hotfix"}},
		{Kind: policy.Assign, Users: []string{"clivern"}},
	}, plan.Actions)
}
