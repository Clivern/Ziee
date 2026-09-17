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

	client := &stubClient{intention: v1.Intention{Name: "bug"}}

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
		Issue: Issue{Author: "newcomer"},
	}, &stubClient{first: true})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"first-contribution"}},
		{Kind: policy.Comment, Body: "Welcome — thanks for opening your first issue here! A maintainer will take a look soon."},
	}, plan.Actions)

	plan = EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

func TestUnitEvaluateIssueOpenedClose(t *testing.T) {
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			Rules: []v1.Rule{
				{
					Name:    "spam",
					When:    v1.Clauses{{AuthorIn: []string{"spammer"}}},
					Labels:  v1.Labels{Add: []string{"spam"}},
					Close:   true,
					Comment: "Closed as spam.",
				},
			},
		},
	}

	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "spammer"},
	}, &stubClient{})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Comment, Body: "Closed as spam."},
		{Kind: policy.Close},
	}, plan.Actions)
}

func TestUnitEvaluateIssueOpenedBlockAuthor(t *testing.T) {
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
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
	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "spammer", Title: "buy crypto"},
	}, client)

	assert.Equal(t, "spammer", client.blockedLogin)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Close},
		{Kind: policy.BlockAuthor, Users: []string{"spammer"}},
	}, plan.Actions)
}

func TestUnitEvaluateIssueOpenedSkipAI(t *testing.T) {
	blocked := true
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			AI:      v1.AI{Enabled: true},
			Rules: []v1.Rule{
				{
					Name:   "spam-blocklist",
					When:   v1.Clauses{{AuthorBlocked: &blocked}},
					Labels: v1.Labels{Add: []string{"spam"}},
					Close:  true,
				},
				{
					Name:   "spam",
					When:   v1.Clauses{{Intention: v1.Intention{Name: "spam"}}},
					Labels: v1.Labels{Add: []string{"spam"}},
				},
			},
		},
	}

	client := &stubClient{blocked: true, intention: v1.Intention{Name: "spam"}}
	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "spammer", Title: "again"},
	}, client)

	assert.Nil(t, client.got)
	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Close},
	}, plan.Actions)
}

func TestUnitEvaluateIssueOpenedRateLimit(t *testing.T) {
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			Rules: []v1.Rule{
				{
					Name: "rate-limit",
					When: v1.Clauses{{
						MaxIssuesOpened: &v1.MaxIssuesOpened{Count: 3, Within: "24h"},
					}},
					Labels: v1.Labels{Add: []string{"spam"}},
					Close:  true,
				},
			},
		},
	}

	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "flooder"},
	}, &stubClient{exceeds: true})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Close},
	}, plan.Actions)

	plan = EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

func TestUnitEvaluateIssueOpenedAuthorBlocked(t *testing.T) {
	blocked := true
	conf := &v1.File{
		IssueTriage: v1.IssueTriage{
			Enabled: true,
			Rules: []v1.Rule{
				{
					Name:   "spam-blocklist",
					When:   v1.Clauses{{AuthorBlocked: &blocked}},
					Labels: v1.Labels{Add: []string{"spam"}},
					Close:  true,
				},
			},
		},
	}

	plan := EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "blocked-user"},
	}, &stubClient{blocked: true})

	assert.Equal(t, []action.Action{
		{Kind: policy.AddLabels, Labels: []string{"spam"}},
		{Kind: policy.Close},
	}, plan.Actions)

	plan = EvaluateIssueOpened(conf, Event{
		Issue: Issue{Author: "maya"},
	}, &stubClient{})

	assert.Empty(t, plan.Actions)
}

type stubClient struct {
	got          []v1.Intention
	intention    v1.Intention
	teams        []string
	org          string
	first        bool
	exceeds      bool
	prsExceeds   bool
	blocked      bool
	blockedLogin string
}

func (s *stubClient) EvaluateIssue(_ Issue, intentions []v1.Intention) v1.Intention {
	s.got = intentions

	return s.intention
}

func (s *stubClient) GetTeams(org, _ string) []string {
	s.org = org

	return s.teams
}

func (s *stubClient) IsFirstContribution(Issue) bool {
	return s.first
}

func (s *stubClient) IssuesOpenedExceeds(Issue, int, string) bool {
	return s.exceeds
}

func (s *stubClient) PrsOpenedExceeds(Issue, int, string) bool {
	return s.prsExceeds
}

func (s *stubClient) IsAuthorBlocked(Issue) bool {
	return s.blocked
}

func (s *stubClient) BlockAuthor(login string) {
	s.blockedLogin = login
	s.blocked = true
}
