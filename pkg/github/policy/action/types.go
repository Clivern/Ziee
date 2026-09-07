// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package action

import "context"

// Action is one GitHub change to apply.
type Action struct {
	Kind   string
	Labels []string
	Users  []string
	Body   string
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

// Client performs planned actions. The GitHub installation client implements this later.
type Client interface {
	AddLabels(ctx context.Context, repo Repo, labels []string) error
	RemoveLabels(ctx context.Context, repo Repo, labels []string) error
	Assign(ctx context.Context, repo Repo, users []string) error
	Unassign(ctx context.Context, repo Repo, users []string) error
	Comment(ctx context.Context, repo Repo, body string) error
	Close(ctx context.Context, repo Repo) error
	Reopen(ctx context.Context, repo Repo) error
}
