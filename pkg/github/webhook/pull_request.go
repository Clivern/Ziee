// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import "time"

// PullRequestEvent is a GitHub pull_request webhook payload.
type PullRequestEvent struct {
	Action       string             `json:"action"`
	Number       int                `json:"number"`
	PullRequest  PullRequestPayload `json:"pull_request"`
	Repository   RepositoryPayload  `json:"repository"`
	Installation InstallationRef    `json:"installation"`
	Sender       UserPayload        `json:"sender"`
}

// PullRequestPayload is the pull request object on a webhook.
type PullRequestPayload struct {
	ID                int64         `json:"id"`
	Number            int           `json:"number"`
	Title             string        `json:"title"`
	Body              string        `json:"body"`
	Draft             bool          `json:"draft"`
	Merged            bool          `json:"merged"`
	Mergeable         *bool         `json:"mergeable"`
	State             string        `json:"state"`
	User              UserPayload   `json:"user"`
	AuthorAssociation string        `json:"author_association"`
	Labels            []NamePayload `json:"labels"`
	Assignees         []UserPayload `json:"assignees"`
	CreatedAt         time.Time     `json:"created_at"`
	MergedAt          *time.Time    `json:"merged_at"`
}
