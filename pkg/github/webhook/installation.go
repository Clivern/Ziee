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

// InstallationRepositoriesEvent is a GitHub App installation_repositories webhook payload.
type InstallationRepositoriesEvent struct {
	Action              string              `json:"action"`
	Installation        InstallationPayload `json:"installation"`
	RepositorySelection string              `json:"repository_selection"`
	RepositoriesAdded   []InstallationRepo  `json:"repositories_added"`
	RepositoriesRemoved []InstallationRepo  `json:"repositories_removed"`
}

// InstallationRepo is a repo added or removed from a GitHub App installation.
type InstallationRepo struct {
	ID       int64  `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}
