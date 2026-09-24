// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package action

import "context"

// Action is one GitHub change to apply.
type Action struct {
	Kind   string   `json:"kind,omitempty"`
	Labels []string `json:"labels,omitempty"`
	Users  []string `json:"users,omitempty"`
	Body   string   `json:"body,omitempty"`
}

// Plan is the ordered list of actions for one event.
type Plan struct {
	Actions []Action
}

// Repo is the GitHub issue/PR the actions apply to.
type Repo struct {
	Owner  string
	Name   string
	Number int
}

// Client performs planned actions on GitHub.
type Client interface {
	AddLabels(ctx context.Context, repo Repo, labels []string) error
	RemoveLabels(ctx context.Context, repo Repo, labels []string) error
	Assign(ctx context.Context, repo Repo, users []string) error
	Unassign(ctx context.Context, repo Repo, users []string) error
	Comment(ctx context.Context, repo Repo, body string) error
	Close(ctx context.Context, repo Repo) error
	Reopen(ctx context.Context, repo Repo) error
	RequestReviewers(ctx context.Context, repo Repo, users []string) error
	RequestReviewTeams(ctx context.Context, repo Repo, teams []string) error
	BlockAuthor(ctx context.Context, repo Repo, users []string) error
}
