// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

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
	Number            int           `json:"number"`
	Title             string        `json:"title"`
	Body              string        `json:"body"`
	User              UserPayload   `json:"user"`
	AuthorAssociation string        `json:"author_association"`
	Labels            []NamePayload `json:"labels"`
	Assignees         []UserPayload `json:"assignees"`
}
