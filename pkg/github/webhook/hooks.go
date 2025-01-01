// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"github.com/actx0/ziee/pkg/github/event"
)

// Hooks holds GitHub webhook callbacks. Unset hooks are no-ops.
type Hooks struct {
	Raw                      func(event.Raw) error
	Status                   func(event.Status) error
	Watch                    func(event.Watch) error
	Issues                   func(event.Issues) error
	IssueComment             func(event.IssueComment) error
	Push                     func(event.Push) error
	Create                   func(event.Create) error
	Label                    func(event.Label) error
	Delete                   func(event.Delete) error
	Milestone                func(event.Milestone) error
	PullRequest              func(event.PullRequest) error
	PullRequestReview        func(event.PullRequestReview) error
	PullRequestReviewComment func(event.PullRequestReviewComment) error
	IssuesCommand            map[string]func(event.Command, event.Issues) error
	IssueCommentCommand      map[string]func(event.Command, event.IssueComment) error
}

// DefaultHooks returns no-op hooks suitable as a starting point.
func DefaultHooks() Hooks {
	noop := func() error { return nil }
	return Hooks{
		Raw:                      func(event.Raw) error { return noop() },
		Status:                   func(event.Status) error { return noop() },
		Watch:                    func(event.Watch) error { return noop() },
		Issues:                   func(event.Issues) error { return noop() },
		IssueComment:             func(event.IssueComment) error { return noop() },
		Push:                     func(event.Push) error { return noop() },
		Create:                   func(event.Create) error { return noop() },
		Label:                    func(event.Label) error { return noop() },
		Delete:                   func(event.Delete) error { return noop() },
		Milestone:                func(event.Milestone) error { return noop() },
		PullRequest:              func(event.PullRequest) error { return noop() },
		PullRequestReview:        func(event.PullRequestReview) error { return noop() },
		PullRequestReviewComment: func(event.PullRequestReviewComment) error { return noop() },
		IssuesCommand:            map[string]func(event.Command, event.Issues) error{},
		IssueCommentCommand:      map[string]func(event.Command, event.IssueComment) error{},
	}
}
