// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"testing"

	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitEvaluatePROpenedSpam(t *testing.T) {
	conf := &v1.File{
		PRTriage: v1.PRTriage{
			Enabled: true,
			AI:      v1.AI{Enabled: true},
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
	}

	client := &stubClient{intention: v1.Intention{Name: "spam"}}
	plan := EvaluatePullRequestOpened(conf, Event{
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
		PRTriage: v1.PRTriage{
			Enabled: true,
			AI:      v1.AI{Enabled: true},
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
	}

	client := &stubClient{prsExceeds: true, intention: v1.Intention{Name: "bug"}}
	plan := EvaluatePullRequestOpened(conf, Event{
		Issue: Issue{Author: "maya", Title: "crash"},
	}, client)

	assert.Equal(t, []v1.Intention{{Name: "bug"}}, client.got)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
	}, plan.Actions)

	client = &stubClient{prsExceeds: true, intention: v1.Intention{Name: "bug"}}
	plan = EvaluatePullRequestOpened(conf, Event{
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
		PRTriage: v1.PRTriage{
			Enabled: true,
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
	}

	client := &stubClient{prsExceeds: true}
	plan := EvaluatePullRequestOpened(conf, Event{
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
		PRTriage: v1.PRTriage{
			Enabled: true,
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
	}

	plan := EvaluatePullRequestOpened(conf, Event{
		Issue: Issue{Author: "clivern"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"team/sre"}},
	}, plan.Actions)

	plan = EvaluatePullRequestOpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

func TestUnitEvaluatePROpenedDraft(t *testing.T) {
	draft := false
	conf := &v1.File{
		PRTriage: v1.PRTriage{
			Enabled: true,
			Rules: []v1.Rule{
				{
					Name:   "ready",
					When:   v1.Clauses{{Draft: &draft}},
					Labels: v1.Labels{Add: []string{"bot"}},
				},
			},
		},
	}

	plan := EvaluatePullRequestOpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"bot"}},
	}, plan.Actions)

	plan = EvaluatePullRequestOpened(conf, Event{
		Issue: Issue{Author: "maya", Draft: true},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

func TestUnitMatchFile(t *testing.T) {
	assert.True(t, MatchFile("api/**", "api/health.go"))
	assert.True(t, MatchFile("**/*.go", "pkg/util/hash.go"))
	assert.False(t, MatchFile("web/**", "api/health.go"))
}

func TestUnitEvaluatePRComment(t *testing.T) {
	viper.Set("app.oauth.github.bot_name", "zieeai")

	conf := &v1.File{
		PRTriage: v1.PRTriage{
			Enabled:  true,
			Comments: policy.CommentsOutcomes,
			Commands: v1.Commands{
				"label":     {Allow: v1.Allow{{Users: []string{"maya"}}}},
				"reviewers": {Allow: v1.Allow{{Users: []string{"maya"}}}},
				"close":     {Allow: v1.Allow{{Users: []string{"clivern"}}}},
				"spam":      {Allow: v1.Allow{{Users: []string{"clivern"}}}},
			},
		},
	}

	plan := EvaluatePullRequestComment(conf, Event{
		Comment: "@zieeai label bug",
		Actor:   Actor{Login: "maya"},
		Issue:   Issue{Author: "guest"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
		{Kind: policy.Comment, Body: "Labeled `bug` as requested by @maya."},
	}, plan.Actions)

	plan = EvaluatePullRequestComment(conf, Event{
		Comment: "@zieeai reviewers clivern",
		Actor:   Actor{Login: "maya"},
		Issue:   Issue{Author: "guest"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.RequestReviewers, Users: []string{"clivern"}},
		{Kind: policy.Comment, Body: "Requested review from @clivern as requested by @maya."},
	}, plan.Actions)

	plan = EvaluatePullRequestComment(conf, Event{
		Comment: "@zieeai close",
		Actor:   Actor{Login: "clivern"},
		Issue:   Issue{Author: "guest"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.Close},
		{Kind: policy.Comment, Body: "Closed this pull request as requested by @clivern."},
	}, plan.Actions)

	client := &stubClient{}
	plan = EvaluatePullRequestComment(conf, Event{
		Comment: "@zieeai spam",
		Actor:   Actor{Login: "clivern"},
		Issue:   Issue{Author: "spammer"},
	}, client)

	assert.Equal(t, "spammer", client.blockedLogin)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Close},
		{Kind: policy.BlockAuthor, Users: []string{"spammer"}},
		{Kind: policy.Comment, Body: "Labeled `spam`, closed this pull request, blocked the author as requested by @clivern."},
	}, plan.Actions)

	plan = EvaluatePullRequestComment(conf, Event{
		Comment: "@zieeai label bug",
		Actor:   Actor{Login: "guest"},
		Issue:   Issue{Author: "guest"},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

func TestUnitEvaluatePRUpdated(t *testing.T) {
	first := true
	conf := &v1.File{
		PRTriage: v1.PRTriage{
			Enabled:  true,
			Comments: policy.CommentsOutcomes,
			Rules: []v1.Rule{
				{
					Name:   "area-api",
					When:   v1.Clauses{{Files: []string{"api/**"}}},
					Labels: v1.Labels{Add: []string{"area/api"}},
					Assign: []string{"clivern"},
				},
				{
					Name:    "first",
					When:    v1.Clauses{{FirstContribution: &first}},
					Comment: "Welcome",
				},
			},
		},
	}

	plan := EvaluatePullRequestUpdated(conf, Event{
		Issue: Issue{
			Author:    "maya",
			Files:     []string{"api/health.go"},
			Labels:    []string{"area/api"},
			Assignees: []string{"clivern"},
		},
		Actor: Actor{Login: "maya"},
	}, &stubClient{first: true})

	assert.Empty(t, plan.Actions)

	plan = EvaluatePullRequestUpdated(conf, Event{
		Issue: Issue{
			Author: "maya",
			Files:  []string{"api/health.go"},
		},
		Actor: Actor{Login: "maya"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"area/api"}},
		{Kind: policy.Assign, Users: []string{"clivern"}},
		{Kind: policy.Comment, Body: "Triaged this pull request: labeled `area/api`, assigned @clivern."},
	}, plan.Actions)
}
