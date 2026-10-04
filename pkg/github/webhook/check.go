// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

// CheckRunEvent is a GitHub check_run webhook payload.
type CheckRunEvent struct {
	Action       string            `json:"action"`
	CheckRun     CheckRunPayload   `json:"check_run"`
	Repository   RepositoryPayload `json:"repository"`
	Installation InstallationRef   `json:"installation"`
	Sender       UserPayload       `json:"sender"`
}

// CheckRunPayload is the check run object on a webhook.
type CheckRunPayload struct {
	ID           int64             `json:"id"`
	Name         string            `json:"name"`
	HeadSHA      string            `json:"head_sha"`
	Status       string            `json:"status"`
	Conclusion   string            `json:"conclusion"`
	PullRequests []CheckRunPullRef `json:"pull_requests"`
}

// CheckRunPullRef is a pull request linked to a check run.
type CheckRunPullRef struct {
	Number int `json:"number"`
}
