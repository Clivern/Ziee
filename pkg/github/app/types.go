// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

// User is a GitHub user on an issue, pull request, or assignee list.
type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

// Issue is a GitHub issue.
type Issue struct {
	ID        int64   `json:"id"`
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	State     string  `json:"state"`
	Body      string  `json:"body"`
	HTMLURL   string  `json:"html_url"`
	User      User    `json:"user"`
	Labels    []Label `json:"labels"`
	Assignees []User  `json:"assignees"`
}

// Comment is a GitHub issue or pull request comment.
type Comment struct {
	ID      int64  `json:"id"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    User   `json:"user"`
}
