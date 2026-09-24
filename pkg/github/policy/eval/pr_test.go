// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"testing"

	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/stretchr/testify/assert"
)

func TestUnitEvaluatePROpenedSpam(t *testing.T) {
	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			PRTriage: v1.PRTriage{
				AI: v1.AI{Enabled: true},
				Rules: []v1.Rule{
					{
						Name: "spam",
						When: v1.Clauses{{
							Intention: v1.Intention{Name: "spam"},
						}},
						Labels:      v1.Labels{Add: []string{"spam"}},
						Close:       true,
						BlockAuthor: true,
					},
				},
			},
		},
	}

	client := &stubClient{intention: v1.Intention{Name: "spam"}}
	plan := EvaluatePROpened(conf, Event{
		Issue: Issue{Author: "spammer", Title: "buy crypto"},
	}, client)

	assert.Equal(t, "spammer", client.blockedLogin)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Close},
		{Kind: policy.BlockAuthor, Users: []string{"spammer"}},
	}, plan.Actions)
}

func TestUnitEvaluatePROpenedSkipAINotRateLimit(t *testing.T) {
	conf := &v1.File{
		Teams: []v1.Team{{Name: "sre", Members: []string{"maya"}}},
		MergeQueue: v1.MergeQueue{
			PRTriage: v1.PRTriage{
				AI: v1.AI{Enabled: true},
				Rules: []v1.Rule{
					{
						Name: "rate-limit",
						When: v1.Clauses{
							{AuthorNotInTeam: []string{"sre"}},
							{MaxPrsOpened: &v1.MaxIssuesOpened{Count: 3, Within: "24h"}},
						},
						Labels: v1.Labels{Add: []string{"spam"}},
						Close:  true,
					},
					{
						Name:   "bug",
						When:   v1.Clauses{{Intention: v1.Intention{Name: "bug"}}},
						Labels: v1.Labels{Add: []string{"bug"}},
					},
				},
			},
		},
	}

	client := &stubClient{prsExceeds: true, intention: v1.Intention{Name: "bug"}}
	plan := EvaluatePROpened(conf, Event{
		Issue: Issue{Author: "maya", Title: "crash"},
	}, client)

	assert.Equal(t, []v1.Intention{{Name: "bug"}}, client.got)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
	}, plan.Actions)

	client = &stubClient{prsExceeds: true, intention: v1.Intention{Name: "bug"}}
	plan = EvaluatePROpened(conf, Event{
		Issue: Issue{Author: "flooder", Title: "crash"},
	}, client)

	assert.Equal(t, []v1.Intention{{Name: "bug"}}, client.got)
	assert.ElementsMatch(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Close},
	}, plan.Actions)
}

func TestUnitEvaluatePROpenedRateLimit(t *testing.T) {
	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			PRTriage: v1.PRTriage{
				Rules: []v1.Rule{
					{
						Name: "rate-limit",
						When: v1.Clauses{{
							MaxPrsOpened: &v1.MaxIssuesOpened{Count: 3, Within: "24h"},
						}},
						Labels:      v1.Labels{Add: []string{"spam"}},
						Close:       true,
						BlockAuthor: true,
					},
				},
			},
		},
	}

	client := &stubClient{prsExceeds: true}
	plan := EvaluatePROpened(conf, Event{
		Issue: Issue{Author: "flooder"},
	}, client)

	assert.Equal(t, "flooder", client.blockedLogin)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Close},
		{Kind: policy.BlockAuthor, Users: []string{"flooder"}},
	}, plan.Actions)
}

func TestUnitEvaluatePROpenedOr(t *testing.T) {
	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			PRTriage: v1.PRTriage{
				Rules: []v1.Rule{
					{
						Name: "trusted",
						When: v1.Clauses{{
							Or: v1.Clauses{
								{AuthorInTeam: []string{"sre"}},
								{AuthorIn: []string{"clivern"}},
							},
						}},
						Labels: v1.Labels{Add: []string{"team/sre"}},
					},
				},
			},
		},
	}

	plan := EvaluatePROpened(conf, Event{
		Issue: Issue{Author: "clivern"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"team/sre"}},
	}, plan.Actions)

	plan = EvaluatePROpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

func TestUnitEvaluatePROpenedDraft(t *testing.T) {
	draft := false
	conf := &v1.File{
		MergeQueue: v1.MergeQueue{
			PRTriage: v1.PRTriage{
				Rules: []v1.Rule{
					{
						Name:   "ready",
						When:   v1.Clauses{{Draft: &draft}},
						Labels: v1.Labels{Add: []string{"bot"}},
					},
				},
			},
		},
	}

	plan := EvaluatePROpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"bot"}},
	}, plan.Actions)

	plan = EvaluatePROpened(conf, Event{
		Issue: Issue{Author: "maya", Draft: true},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

func TestUnitMatchFile(t *testing.T) {
	assert.True(t, MatchFile("api/**", "api/health.go"))
	assert.True(t, MatchFile("**/*.go", "pkg/util/hash.go"))
	assert.False(t, MatchFile("web/**", "api/health.go"))
}
