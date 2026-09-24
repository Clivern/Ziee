// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package action

import (
	"context"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/github/app"
)

type githubClient struct {
	installationId int64
	githubRepoId   int64
	authorId       int64
	repos          *module.Repository
}

// NewClient returns a GitHub client that performs planned triage actions.
func NewClient(installationId, githubRepoId, authorId int64) Client {
	return githubClient{
		installationId: installationId,
		githubRepoId:   githubRepoId,
		authorId:       authorId,
		repos: module.NewRepository(
			db.NewRepositoriesRepository(db.GetDB()),
			db.NewRepositoryMetaRepository(db.GetDB()),
			db.NewRepositorySpamUserRepository(db.GetDB()),
		),
	}
}

// AddLabels adds labels to the issue or pull request.
func (a githubClient) AddLabels(ctx context.Context, repo Repo, labels []string) error {
	return app.Get().AddLabels(ctx, a.installationId, repo.Owner, repo.Name, repo.Number, labels)
}

// RemoveLabels removes labels from the issue or pull request.
func (a githubClient) RemoveLabels(ctx context.Context, repo Repo, labels []string) error {
	return app.Get().RemoveLabels(ctx, a.installationId, repo.Owner, repo.Name, repo.Number, labels)
}

// Assign assigns users to the issue or pull request.
func (a githubClient) Assign(ctx context.Context, repo Repo, users []string) error {
	return app.Get().AddAssignees(ctx, a.installationId, repo.Owner, repo.Name, repo.Number, users)
}

// Unassign removes users from the issue or pull request.
func (a githubClient) Unassign(ctx context.Context, repo Repo, users []string) error {
	return app.Get().RemoveAssignees(ctx, a.installationId, repo.Owner, repo.Name, repo.Number, users)
}

// Comment adds a comment on the issue or pull request.
func (a githubClient) Comment(ctx context.Context, repo Repo, body string) error {
	_, err := app.Get().CreateComment(ctx, a.installationId, repo.Owner, repo.Name, repo.Number, body)

	return err
}

// Close closes the issue or pull request.
func (a githubClient) Close(ctx context.Context, repo Repo) error {
	return app.Get().CloseIssue(ctx, a.installationId, repo.Owner, repo.Name, repo.Number)
}

// Reopen reopens the issue or pull request.
func (a githubClient) Reopen(ctx context.Context, repo Repo) error {
	return app.Get().ReopenIssue(ctx, a.installationId, repo.Owner, repo.Name, repo.Number)
}

// RequestReviewers asks users to review the pull request.
func (a githubClient) RequestReviewers(ctx context.Context, repo Repo, users []string) error {
	return app.Get().RequestReviewers(ctx, a.installationId, repo.Owner, repo.Name, repo.Number, users, nil)
}

// RequestReviewTeams asks teams to review the pull request.
func (a githubClient) RequestReviewTeams(ctx context.Context, repo Repo, teams []string) error {
	return app.Get().RequestReviewers(ctx, a.installationId, repo.Owner, repo.Name, repo.Number, nil, teams)
}

// BlockAuthor adds the issue author to the repository spam blocklist.
func (a githubClient) BlockAuthor(ctx context.Context, repo Repo, users []string) error {
	a.repos.BlockAuthor(a.githubRepoId, users[0], a.authorId)

	return nil
}
