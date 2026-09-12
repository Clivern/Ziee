// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

// IssueEvent is a GitHub issues webhook payload.
type IssueEvent struct {
	Action       string            `json:"action"`
	Issue        IssuePayload      `json:"issue"`
	Repository   RepositoryPayload `json:"repository"`
	Installation InstallationRef   `json:"installation"`
	Sender       UserPayload       `json:"sender"`
}

// IssueCommentEvent is a GitHub issue_comment webhook payload.
type IssueCommentEvent struct {
	Action       string            `json:"action"`
	Issue        IssuePayload      `json:"issue"`
	Comment      CommentPayload    `json:"comment"`
	Repository   RepositoryPayload `json:"repository"`
	Installation InstallationRef   `json:"installation"`
	Sender       UserPayload       `json:"sender"`
}

// InstallationRef is the installation id on a repository webhook.
type InstallationRef struct {
	ID int64 `json:"id"`
}

// RepositoryPayload is the repository object on a webhook.
type RepositoryPayload struct {
	ID            int64       `json:"id"`
	Name          string      `json:"name"`
	FullName      string      `json:"full_name"`
	Owner         UserPayload `json:"owner"`
	DefaultBranch string      `json:"default_branch"`
}

// IssuePayload is the issue object on a webhook.
type IssuePayload struct {
	Number            int           `json:"number"`
	Title             string        `json:"title"`
	Body              string        `json:"body"`
	User              UserPayload   `json:"user"`
	AuthorAssociation string        `json:"author_association"`
	Labels            []NamePayload `json:"labels"`
	Assignees         []UserPayload `json:"assignees"`
}

// CommentPayload is the comment object on an issue_comment webhook.
type CommentPayload struct {
	Body string      `json:"body"`
	User UserPayload `json:"user"`
}

// UserPayload is a GitHub login and optional account type.
type UserPayload struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

// NamePayload is a named GitHub object such as a label.
type NamePayload struct {
	Name string `json:"name"`
}
