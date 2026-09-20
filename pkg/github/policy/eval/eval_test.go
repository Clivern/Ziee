// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"
	"github.com/clivern/ziee/pkg/github/webhook"

	"github.com/stretchr/testify/assert"
)

func TestUnitRunIssueOpenedPayload(t *testing.T) {
	body, err := os.ReadFile("../../../../testdata/issues_opened.json")
	assert.NoError(t, err)

	var payload webhook.IssueEvent
	err = json.Unmarshal(body, &payload)
	assert.NoError(t, err)

	labels := make([]string, len(payload.Issue.Labels))
	for i, label := range payload.Issue.Labels {
		labels[i] = label.Name
	}

	assignees := make([]string, len(payload.Issue.Assignees))
	for i, user := range payload.Issue.Assignees {
		assignees[i] = user.Login
	}

	event := Event{
		Kind: "issue." + payload.Action,
		Account: Account{
			Login: payload.Repository.Owner.Login,
			Type:  payload.Repository.Owner.Type,
		},
		Issue: Issue{
			Number:      payload.Issue.Number,
			Title:       payload.Issue.Title,
			Body:        payload.Issue.Body,
			Author:      payload.Issue.User.Login,
			AuthorId:    payload.Issue.User.ID,
			Association: payload.Issue.AuthorAssociation,
			Labels:      labels,
			Assignees:   assignees,
			Closed:      payload.Issue.State == "closed",
		},
		Actor: Actor{
			Login: payload.Sender.Login,
		},
	}

	assert.Equal(t, policy.KindIssueOpened, event.Kind)
	assert.Equal(t, "Clivern", event.Issue.Author)
	assert.Equal(t, "Test the issue request", event.Issue.Title)
	assert.False(t, event.Issue.Closed)

	t.Run("disabled", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{Enabled: false},
		}, event, &stubClient{})

		assert.Empty(t, plan.Actions)
	})

	t.Run("unknown kind", func(t *testing.T) {
		other := event
		other.Kind = "issue.assigned"

		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{Enabled: true},
		}, other, &stubClient{})

		assert.Empty(t, plan.Actions)
	})

	t.Run("no rules match", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled:  true,
				Comments: policy.CommentsAll,
				Rules: []v1.Rule{{
					Name:   "hotfix",
					When:   v1.Clauses{{Title: "hotfix"}},
					Labels: v1.Labels{Add: []string{"hotfix"}},
				}},
			},
		}, event, &stubClient{})

		assert.Equal(t, []action.Action{
			{Kind: policy.Comment, Body: "Triage ran; no rules matched."},
		}, plan.Actions)
	})

	t.Run("title matches", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				Rules: []v1.Rule{{
					Name:    "test",
					When:    v1.Clauses{{Title: "test the issue"}},
					Labels:  v1.Labels{Add: []string{"needs-triage"}},
					Comment: "Thanks for opening this.",
				}},
			},
		}, event, &stubClient{})

		assert.Equal(t, []action.Action{
			{Kind: policy.AddLabels, Labels: []string{"needs-triage"}},
			{Kind: policy.Comment, Body: "Thanks for opening this."},
		}, plan.Actions)
	})

	t.Run("author in", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				Rules: []v1.Rule{{
					Name:   "owner",
					When:   v1.Clauses{{AuthorIn: []string{"Clivern"}}},
					Labels: v1.Labels{Add: []string{"from-owner"}},
				}},
			},
		}, event, &stubClient{})

		assert.Equal(t, []action.Action{
			{Kind: policy.AddLabels, Labels: []string{"from-owner"}},
		}, plan.Actions)
	})

	t.Run("author not bot", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				Rules: []v1.Rule{{
					Name:   "bot",
					When:   v1.Clauses{{AuthorIn: []string{"dependabot", "ziee-bot"}}},
					Labels: v1.Labels{Add: []string{"bot"}},
				}},
			},
		}, event, &stubClient{})

		assert.Empty(t, plan.Actions)
	})

	t.Run("not sre", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				Rules: []v1.Rule{{
					Name:    "not-sre",
					When:    v1.Clauses{{AuthorNotInTeam: []string{"sre"}}},
					Comment: "Thanks for opening this. Someone from the team will take a look.",
				}},
			},
		}, event, &stubClient{})

		assert.Equal(t, []action.Action{
			{Kind: policy.Comment, Body: "Thanks for opening this. Someone from the team will take a look."},
		}, plan.Actions)
	})

	t.Run("from sre", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				Rules: []v1.Rule{{
					Name:   "from-sre",
					When:   v1.Clauses{{AuthorInTeam: []string{"sre"}}},
					Labels: v1.Labels{Add: []string{"team/sre"}},
				}},
			},
		}, event, &stubClient{})

		assert.Empty(t, plan.Actions)
	})

	t.Run("first contribution", func(t *testing.T) {
		first := true
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				Rules: []v1.Rule{{
					Name:    "first-contribution",
					When:    v1.Clauses{{FirstContribution: &first}},
					Labels:  v1.Labels{Add: []string{"first-contribution"}},
					Comment: "Welcome — thanks for opening your first issue here! A maintainer will take a look soon.",
				}},
			},
		}, event, &stubClient{first: true})

		assert.Equal(t, []action.Action{
			{Kind: policy.AddLabels, Labels: []string{"first-contribution"}},
			{Kind: policy.Comment, Body: "Welcome — thanks for opening your first issue here! A maintainer will take a look soon."},
		}, plan.Actions)
	})

	t.Run("rate limit", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				Rules: []v1.Rule{{
					Name: "rate-limit",
					When: v1.Clauses{
						{AuthorNotIn: []string{"dependabot", "ziee-bot"}},
						{MaxIssuesOpened: &v1.MaxIssuesOpened{Count: 3, Within: "24h"}},
					},
					Labels:      v1.Labels{Add: []string{"spam"}},
					Close:       true,
					BlockAuthor: true,
					Comment:     "Too many issues in a short time. Please wait before opening another.",
				}},
			},
		}, event, &stubClient{exceeds: true})

		assert.Equal(t, []action.Action{
			{Kind: policy.AddLabels, Labels: []string{"spam"}},
			{Kind: policy.Comment, Body: "Too many issues in a short time. Please wait before opening another."},
			{Kind: policy.Close},
			{Kind: policy.BlockAuthor, Users: []string{"Clivern"}},
		}, plan.Actions)
	})

	t.Run("blocked author", func(t *testing.T) {
		blocked := true
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				Rules: []v1.Rule{{
					Name:    "spam-blocklist",
					When:    v1.Clauses{{AuthorBlocked: &blocked}},
					Labels:  v1.Labels{Add: []string{"spam"}},
					Close:   true,
					Comment: "This account is blocked from opening issues here.",
				}},
			},
		}, event, &stubClient{blocked: true})

		assert.Equal(t, []action.Action{
			{Kind: policy.AddLabels, Labels: []string{"spam"}},
			{Kind: policy.Comment, Body: "This account is blocked from opening issues here."},
			{Kind: policy.Close},
		}, plan.Actions)
	})

	t.Run("intention bug", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled: true,
				AI:      v1.AI{Enabled: true},
				Rules: []v1.Rule{{
					Name: "bug",
					When: v1.Clauses{
						{AuthorNotIn: []string{"dependabot", "ziee-bot"}},
						{Intention: v1.Intention{Name: "bug", Description: "A defect or unexpected behavior that needs a fix."}},
					},
					Labels:  v1.Labels{Add: []string{"bug"}},
					Assign:  []string{"clivern"},
					Comment: "Thanks for the report. We'll investigate this bug.",
				}},
			},
		}, event, &stubClient{intention: v1.Intention{Name: "bug"}})

		assert.Equal(t, []action.Action{
			{Kind: policy.AddLabels, Labels: []string{"bug"}},
			{Kind: policy.Assign, Users: []string{"clivern"}},
			{Kind: policy.Comment, Body: "Thanks for the report. We'll investigate this bug."},
		}, plan.Actions)
	})

	t.Run("outcomes comment", func(t *testing.T) {
		plan := Run(&v1.File{
			IssueTriage: v1.IssueTriage{
				Enabled:  true,
				Comments: policy.CommentsOutcomes,
				Rules: []v1.Rule{{
					Name:   "test",
					When:   v1.Clauses{{Body: "issue request"}},
					Labels: v1.Labels{Add: []string{"needs-triage"}},
				}},
			},
		}, event, &stubClient{})

		assert.Equal(t, []action.Action{
			{Kind: policy.AddLabels, Labels: []string{"needs-triage"}},
			{Kind: policy.Comment, Body: "Triaged this issue: labeled `needs-triage`."},
		}, plan.Actions)
	})
}
