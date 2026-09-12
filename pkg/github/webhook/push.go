// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import "strings"

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

// ChangedPath returns the first added or modified file that contains one of the names.
func (e PushEvent) ChangedPath(names ...string) string {
	for _, commit := range e.Commits {
		for _, file := range append(commit.Added, commit.Modified...) {
			for _, name := range names {
				if strings.Contains(file, name) {
					return file
				}
			}
		}
	}

	return ""
}
