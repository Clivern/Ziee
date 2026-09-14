// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"

	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"
)

// IssueClient is a GitHub client for issue events.
type IssueClient struct {
	ctx            context.Context
	installationId int64
	owner          string
	repo           string
}

// GetTeams returns GitHub team slugs in org that include login.
func (c IssueClient) GetTeams(org, login string) []string {
	if org == "" || login == "" {
		return []string{}
	}

	teams, _ := app.Get().ListUserTeams(
		c.ctx,
		c.installationId,
		org,
		login,
	)

	slugs := make([]string, len(teams))
	for i, team := range teams {
		slugs[i] = team.Slug
	}

	return slugs
}

// EvaluateIssue classifies issue intention from title and body.
func (c IssueClient) EvaluateIssue(issue eval.Issue, intentions []v1.Intention) []string {
	return []string{}
}

// IsFirstContribution reports whether the issue author is a first-time contributor.
func (c IssueClient) IsFirstContribution(issue eval.Issue) bool {
	first, _ := app.Get().IsFirstIssue(
		c.ctx,
		c.installationId,
		c.owner,
		c.repo,
		issue.Author,
	)

	return first
}
