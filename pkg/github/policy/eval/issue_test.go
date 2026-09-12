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

func TestUnitEvaluateIssueOpened(t *testing.T) {
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			Rules: []v1.Rule{
				{
					Name:    "hotfix",
					When:    v1.Clauses{{Title: "hotfix"}},
					Labels:  v1.Labels{Add: []string{"hotfix"}},
					Assign:  []string{"clivern"},
					Comment: "Hotfix acknowledged.",
				},
				{
					Name:   "bot",
					When:   v1.Clauses{{AuthorIn: []string{"dependabot"}}},
					Labels: v1.Labels{Add: []string{"bot"}},
				},
			},
		},
	}

	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{
			Title:  "HOTFIX: cache",
			Author: "maya",
		},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"hotfix"}},
		{Kind: policy.Assign, Users: []string{"clivern"}},
		{Kind: policy.Comment, Body: "Hotfix acknowledged."},
	}, plan.Actions)
}

func TestUnitEvaluateIssueOpenedIntentions(t *testing.T) {
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			AI:      v1.AI{Enabled: true},
			Rules: []v1.Rule{
				{
					Name:   "bug",
					When:   v1.Clauses{{Intention: v1.Intention{Name: "bug", Description: "A defect that needs a fix."}}},
					Labels: v1.Labels{Add: []string{"bug"}},
				},
				{
					Name:   "docs",
					When:   v1.Clauses{{Intention: v1.Intention{Name: "docs"}}},
					Labels: v1.Labels{Add: []string{"docs"}},
				},
			},
		},
	}

	client := &stubClient{intentions: []string{"bug"}}

	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Title: "crash on save"},
	}, client)

	assert.Equal(t, []v1.Intention{
		{Name: "bug", Description: "A defect that needs a fix."},
		{Name: "docs"},
	}, client.got)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
	}, plan.Actions)
}

func TestUnitEvaluateIssueOpenedTeams(t *testing.T) {
	conf := &v1.File{
		Teams: []v1.Team{
			{Name: "sre", Members: []string{"maya"}},
		},
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			Rules: []v1.Rule{
				{
					Name:   "from-sre",
					When:   v1.Clauses{{AuthorInTeam: []string{"sre"}}},
					Labels: v1.Labels{Add: []string{"team/sre"}},
				},
				{
					Name:   "from-core",
					When:   v1.Clauses{{AuthorInTeam: []string{"core"}}},
					Labels: v1.Labels{Add: []string{"area/api"}},
				},
				{
					Name:   "not-sre",
					When:   v1.Clauses{{AuthorNotInTeam: []string{"sre"}}},
					Labels: v1.Labels{Add: []string{"community"}},
				},
			},
		},
	}

	client := &stubClient{teams: []string{"core"}}

	plan := EvaluateIssueOpened(conf, Event{
		Account: Account{Type: "Organization", Login: "acme"},
		Issue:   Issue{Author: "maya"},
	}, client)

	assert.Equal(t, "acme", client.org)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"team/sre"}},
		{Kind: policy.AddLabels, Labels: []string{"area/api"}},
	}, plan.Actions)

	client = &stubClient{}

	plan = EvaluateIssueOpened(conf, Event{
		Account: Account{Type: "User", Login: "maya"},
		Issue:   Issue{Author: "maya"},
	}, client)

	assert.Empty(t, client.org)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"team/sre"}},
	}, plan.Actions)
}

func TestUnitEvaluateIssueOpenedOutcome(t *testing.T) {
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled:  true,
			Comments: policy.CommentsOutcomes,
			Rules: []v1.Rule{
				{
					Name:   "bug",
					When:   v1.Clauses{{Title: "crash"}},
					Labels: v1.Labels{Add: []string{"bug"}},
					Assign: []string{"clivern"},
				},
			},
		},
	}

	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Title: "crash on save"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
		{Kind: policy.Assign, Users: []string{"clivern"}},
		{Kind: policy.Comment, Body: "Triaged this issue: labeled `bug`, assigned @clivern."},
	}, plan.Actions)
}

func TestUnitEvaluateIssueOpenedOutcomeAllNone(t *testing.T) {
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled:  true,
			Comments: policy.CommentsAll,
			Rules: []v1.Rule{
				{
					Name:   "bot",
					When:   v1.Clauses{{AuthorIn: []string{"dependabot"}}},
					Labels: v1.Labels{Add: []string{"bot"}},
				},
			},
		},
	}

	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.Comment, Body: "Triage ran; no rules matched."},
	}, plan.Actions)

	conf.IssueTriage.Comments = policy.CommentsNone

	plan = EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

func TestUnitEvaluateIssueOpenedFirstContribution(t *testing.T) {
	first := true
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			Rules: []v1.Rule{
				{
					Name:    "first-contribution",
					When:    v1.Clauses{{FirstContribution: &first}},
					Labels:  v1.Labels{Add: []string{"first-contribution"}},
					Comment: "Welcome — thanks for opening your first issue here! A maintainer will take a look soon.",
				},
			},
		},
	}

	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "newcomer", Association: "FIRST_TIME_CONTRIBUTOR"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"first-contribution"}},
		{Kind: policy.Comment, Body: "Welcome — thanks for opening your first issue here! A maintainer will take a look soon."},
	}, plan.Actions)

	plan = EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "maya", Association: "MEMBER"},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

type stubClient struct {
	got        []v1.Intention
	intentions []string
	teams      []string
	org        string
}

func (s *stubClient) EvaluateIssue(_ Issue, intentions []v1.Intention) []string {
	s.got = intentions

	return s.intentions
}

func (s *stubClient) GetTeams(org, _ string) []string {
	s.org = org

	return s.teams
}
