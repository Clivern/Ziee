// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

// Event is one GitHub delivery eval can evaluate.
type Event struct {
	Kind    string
	Issue   Issue
	Comment string
	Actor   Actor
}

// Issue is the ticket the event is about.
type Issue struct {
	Number    int
	Title     string
	Body      string
	Author    string
	Labels    []string
	Assignees []string
	Teams     []string
	Files     []string
	Intention string
}

// Actor is who opened the issue or typed the comment.
type Actor struct {
	Login      string
	Permission string
	Teams      []string
}

// Client is the GitHub client.
type Client interface {
}
