// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"
)

// Event is one GitHub delivery eval can evaluate.
type Event struct {
	Kind    string
	Org     string
	Account Account
	Issue   Issue
	Comment string
	Actor   Actor
}

// Account is the user or org the GitHub App is installed on.
type Account struct {
	Login string
	Type  string
}

// Issue is the ticket the event is about.
type Issue struct {
	Number      int
	Title       string
	Body        string
	Author      string
	AuthorId    int64
	Association string
	Labels      []string
	Assignees   []string
	Teams       []string
	Files       []string
	Intention   v1.Intention
}

// Actor is who opened the issue or typed the comment.
type Actor struct {
	Login      string
	Permission string
	Teams      []string
}

// Client is the GitHub and AI client.
type Client interface {
	// GetTeams returns GitHub team slugs in org that include login.
	GetTeams(org, login string) []string
	// EvaluateIssue classifies issue intention from title and body.
	// Intention descriptions are added to the classify prompt.
	EvaluateIssue(issue Issue, intentions []v1.Intention) v1.Intention
	// IsFirstContribution reports whether the issue author is a first-time contributor.
	IsFirstContribution(issue Issue) bool
	// IssuesOpenedExceeds reports whether the author opened more than count issues within the duration.
	IssuesOpenedExceeds(issue Issue, count int, within string) bool
	// PrsOpenedExceeds reports whether the author opened more than count pull requests within the duration.
	PrsOpenedExceeds(issue Issue, count int, within string) bool
	// IsAuthorBlocked reports whether the author is on the repository spam blocklist.
	IsAuthorBlocked(issue Issue) bool
	// BlockAuthor adds the author to the repository spam blocklist.
	BlockAuthor(issue Issue)
}
