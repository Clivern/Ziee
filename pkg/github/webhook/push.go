// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import "github.com/samber/lo"

// PushEvent is a GitHub push webhook payload.
type PushEvent struct {
	Ref          string            `json:"ref"`
	Repository   RepositoryPayload `json:"repository"`
	Installation InstallationRef   `json:"installation"`
	Commits      []PushCommit      `json:"commits"`
}

// PushCommit is one commit on a push webhook.
type PushCommit struct {
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

// Changed reports whether any commit added, modified, or removed one of the paths.
func (e PushEvent) Changed(paths ...string) bool {
	for _, commit := range e.Commits {
		if lo.Some(lo.Union(commit.Added, commit.Removed, commit.Modified), paths) {
			return true
		}
	}

	return false
}
