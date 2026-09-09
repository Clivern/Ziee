// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"
)

// Event is one GitHub delivery eval can evaluate.
type Event struct {
	Kind    string
	Issue   Issue
	Comment string
	Actor   Actor
}

// Issue is the ticket the event is about.
type Issue struct {
	Number     int
	Title      string
	Body       string
	Author     string
	Labels     []string
	Assignees  []string
	Teams      []string
	Files      []string
	Intentions []string
}

// Actor is who opened the issue or typed the comment.
type Actor struct {
	Login      string
	Permission string
	Teams      []string
}

// Client is the GitHub and AI client.
type Client interface {
	// GetTeams returns GitHub org team slugs for login.
	GetTeams(login string) []string
	// EvaluateIssue classifies issue intention from title and body.
	// Intention descriptions are added to the classify prompt.
	EvaluateIssue(issue Issue, intentions []v1.Intention) []string
}
