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

// ChangedPath returns the first added or modified path that matches.
func (e PushEvent) ChangedPath(paths ...string) string {
	for _, commit := range e.Commits {
		files := lo.Union(commit.Added, commit.Modified)
		for _, path := range paths {
			if lo.Contains(files, path) {
				return path
			}
		}
	}

	return ""
}
