// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

// InstallationEvent is a GitHub App installation webhook payload.
type InstallationEvent struct {
	Action       string              `json:"action"`
	Installation InstallationPayload `json:"installation"`
	Sender       struct {
		ID int64 `json:"id"`
	} `json:"sender"`
}

// InstallationPayload is the installation object on a webhook.
type InstallationPayload struct {
	ID                  int64  `json:"id"`
	HTMLURL             string `json:"html_url"`
	RepositorySelection string `json:"repository_selection"`
	Account             struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Type  string `json:"type"`
	} `json:"account"`
}
