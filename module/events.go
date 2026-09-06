// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/github"
	"github.com/clivern/ziee/pkg/github/webhook"

	"github.com/rs/zerolog/log"
)

// RegisterEventListeners wires application event handlers.
func RegisterEventListeners() {
	webhook.Installation.On(func(_ context.Context, payload webhook.InstallationEvent) {
		i := NewInstallation(
			db.NewGitHubInstallationRepository(db.GetDB()),
			db.NewWorkspaceGitHubRepoRepository(db.GetDB()),
		)

		if payload.Action == "deleted" {
			err := i.Delete(payload.Installation.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to persist GitHub installation")
			}
			return
		}

		meta, err := json.Marshal(payload.Installation)
		if err != nil {
			log.Error().Err(err).Msg("Failed to persist GitHub installation")
			return
		}
		raw := string(meta)

		err = i.Upsert(&db.GitHubInstallation{
			GitHubId:            payload.Installation.ID,
			GitHubUserId:        strconv.FormatInt(payload.Sender.ID, 10),
			AccountId:           payload.Installation.Account.ID,
			AccountLogin:        payload.Installation.Account.Login,
			AccountType:         payload.Installation.Account.Type,
			RepositorySelection: payload.Installation.RepositorySelection,
			HTMLURL:             payload.Installation.HTMLURL,
			Meta:                &raw,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to persist GitHub installation")
		}
	})

	webhook.InstallationRepositories.On(func(_ context.Context, payload webhook.InstallationRepositoriesEvent) {
		i := NewInstallation(
			db.NewGitHubInstallationRepository(db.GetDB()),
			db.NewWorkspaceGitHubRepoRepository(db.GetDB()),
		)

		added := make([]github.Repository, len(payload.RepositoriesAdded))
		for n, repo := range payload.RepositoriesAdded {
			added[n] = github.Repository{
				ID:       repo.ID,
				NodeID:   repo.NodeID,
				Name:     repo.Name,
				FullName: repo.FullName,
				Private:  repo.Private,
			}
		}

		removed := make([]int64, len(payload.RepositoriesRemoved))
		for n, repo := range payload.RepositoriesRemoved {
			removed[n] = repo.ID
		}

		err := i.UpdateRepositories(payload.Installation.ID, added, removed)
		if err != nil {
			log.Error().Err(err).Msg("Failed to persist GitHub installation repositories")
		}
	})
}
