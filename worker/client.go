// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"time"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/ai"
	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// IssueClient is a GitHub client for issue events.
type IssueClient struct {
	ctx            context.Context
	installationId int64
	githubRepoId   int64
	owner          string
	repo           string
	repos          *module.Repository
}

// PullRequestClient is a GitHub client for pull request events.
type PullRequestClient struct {
	ctx            context.Context
	installationId int64
	githubRepoId   int64
	owner          string
	repo           string
	repos          *module.Repository
}

// NewIssueClient returns a GitHub client for issue events.
func NewIssueClient(ctx context.Context, installationId, githubRepoId int64, owner, repo string) IssueClient {
	return IssueClient{
		ctx:            ctx,
		installationId: installationId,
		githubRepoId:   githubRepoId,
		owner:          owner,
		repo:           repo,
		repos: module.NewRepository(
			db.NewRepositoriesRepository(db.GetDB()),
			db.NewRepositoryMetaRepository(db.GetDB()),
			db.NewRepositorySpamUserRepository(db.GetDB()),
		),
	}
}

// GetTeams returns GitHub team slugs in org that include login.
func (c IssueClient) GetTeams(org, login string) []string {
	if org == "" || login == "" {
		return []string{}
	}

	teams, err := app.Get().ListUserTeams(
		c.ctx,
		c.installationId,
		org,
		login,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Str("org", org).
			Str("login", login).
			Msg("Failed to list issue author teams")
	}

	slugs := make([]string, len(teams))
	for i, team := range teams {
		slugs[i] = team.Slug
	}

	return slugs
}

// EvaluateIssue classifies issue intention from title and body.
func (c IssueClient) EvaluateIssue(issue eval.Issue, intentions []v1.Intention) v1.Intention {
	usage := module.NewUsage()

	options := make([]ai.ClassifyOption, len(intentions))
	for i, intention := range intentions {
		options[i] = ai.ClassifyOption{
			Name:        intention.Name,
			Description: intention.Description,
		}
	}

	name, aiUsage, err := ai.NewClassifyClient().Classify(c.ctx, issue.Title, issue.Body, options)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Int("number", issue.Number).
			Msg("Failed to classify issue")
	}

	// record usage
	usage.IncrementAIUsage(
		db.NewUsageRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		c.repos.WorkspaceId(c.githubRepoId),
		aiUsage.TotalTokens,
		aiUsage.Cost,
	)

	if lo.IsEmpty(name) {
		return v1.Intention{}
	}

	matched, _ := lo.Find(intentions, func(intention v1.Intention) bool {
		return intention.Name == name
	})

	return matched
}

// IsFirstContribution reports whether the issue author is a first-time contributor.
func (c IssueClient) IsFirstContribution(issue eval.Issue) bool {
	first, err := app.Get().IsFirstIssue(
		c.ctx,
		c.installationId,
		c.owner,
		c.repo,
		issue.Author,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Int("number", issue.Number).
			Str("author", issue.Author).
			Msg("Failed to check first issue")
	}

	return first
}

// IssuesOpenedExceeds reports whether the author opened more than count issues within the duration.
func (c IssueClient) IssuesOpenedExceeds(issue eval.Issue, count int, within string) bool {
	window, err := time.ParseDuration(within)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Int("number", issue.Number).
			Str("within", within).
			Msg("Failed to parse issue window")
	}

	opened, err := app.Get().CountIssuesOpened(
		c.ctx,
		c.installationId,
		c.owner,
		c.repo,
		issue.Author,
		time.Now().UTC().Add(-window),
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Int("number", issue.Number).
			Str("author", issue.Author).
			Msg("Failed to count issues opened")
	}

	return opened > count
}

// PrsOpenedExceeds is unused for issue events.
func (c IssueClient) PrsOpenedExceeds(eval.Issue, int, string) bool {
	return false
}

// IsAuthorBlocked reports whether the author is on the repository spam blocklist.
func (c IssueClient) IsAuthorBlocked(issue eval.Issue) bool {
	return c.repos.IsAuthorBlocked(c.githubRepoId, issue.AuthorId)
}

// BlockAuthor adds the author to the repository spam blocklist.
func (c IssueClient) BlockAuthor(issue eval.Issue) {
	c.repos.BlockAuthor(c.githubRepoId, issue.Author, issue.AuthorId)
}

// NewPullRequestClient returns a GitHub client for pull request events.
func NewPullRequestClient(ctx context.Context, installationId, githubRepoId int64, owner, repo string) PullRequestClient {
	return PullRequestClient{
		ctx:            ctx,
		installationId: installationId,
		githubRepoId:   githubRepoId,
		owner:          owner,
		repo:           repo,
		repos: module.NewRepository(
			db.NewRepositoriesRepository(db.GetDB()),
			db.NewRepositoryMetaRepository(db.GetDB()),
			db.NewRepositorySpamUserRepository(db.GetDB()),
		),
	}
}

// GetTeams returns GitHub team slugs in org that include login.
func (c PullRequestClient) GetTeams(org, login string) []string {
	if org == "" || login == "" {
		return []string{}
	}

	teams, err := app.Get().ListUserTeams(
		c.ctx,
		c.installationId,
		org,
		login,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Str("org", org).
			Str("login", login).
			Msg("Failed to list pull request author teams")
	}

	slugs := make([]string, len(teams))
	for i, team := range teams {
		slugs[i] = team.Slug
	}

	return slugs
}

// EvaluateIssue classifies pull request intention from title and body.
func (c PullRequestClient) EvaluateIssue(issue eval.Issue, intentions []v1.Intention) v1.Intention {
	usage := module.NewUsage()

	options := make([]ai.ClassifyOption, len(intentions))
	for i, intention := range intentions {
		options[i] = ai.ClassifyOption{
			Name:        intention.Name,
			Description: intention.Description,
		}
	}

	name, aiUsage, err := ai.NewClassifyClient().Classify(c.ctx, issue.Title, issue.Body, options)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Int("number", issue.Number).
			Msg("Failed to classify pull request")
	}

	// record usage
	usage.IncrementAIUsage(
		db.NewUsageRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		c.repos.WorkspaceId(c.githubRepoId),
		aiUsage.TotalTokens,
		aiUsage.Cost,
	)

	if lo.IsEmpty(name) {
		return v1.Intention{}
	}

	matched, _ := lo.Find(intentions, func(intention v1.Intention) bool {
		return intention.Name == name
	})

	return matched
}

// IsFirstContribution reports whether the author is opening their first pull request.
func (c PullRequestClient) IsFirstContribution(issue eval.Issue) bool {
	first, err := app.Get().IsFirstPullRequest(
		c.ctx,
		c.installationId,
		c.owner,
		c.repo,
		issue.Author,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Int("number", issue.Number).
			Str("author", issue.Author).
			Msg("Failed to check first pull request")
	}

	return first
}

// IssuesOpenedExceeds is unused for pull request events.
func (c PullRequestClient) IssuesOpenedExceeds(eval.Issue, int, string) bool {
	return false
}

// PrsOpenedExceeds reports whether the author opened more than count pull requests within the duration.
func (c PullRequestClient) PrsOpenedExceeds(issue eval.Issue, count int, within string) bool {
	window, err := time.ParseDuration(within)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Int("number", issue.Number).
			Str("within", within).
			Msg("Failed to parse pull request window")
	}

	opened, err := app.Get().CountPullRequestsOpened(
		c.ctx,
		c.installationId,
		c.owner,
		c.repo,
		issue.Author,
		time.Now().UTC().Add(-window),
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", c.owner).
			Str("repo", c.repo).
			Int("number", issue.Number).
			Str("author", issue.Author).
			Msg("Failed to count pull requests opened")
	}

	return opened > count
}

// IsAuthorBlocked reports whether the author is on the repository spam blocklist.
func (c PullRequestClient) IsAuthorBlocked(issue eval.Issue) bool {
	return c.repos.IsAuthorBlocked(c.githubRepoId, issue.AuthorId)
}

// BlockAuthor adds the author to the repository spam blocklist.
func (c PullRequestClient) BlockAuthor(issue eval.Issue) {
	c.repos.BlockAuthor(c.githubRepoId, issue.Author, issue.AuthorId)
}
