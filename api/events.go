// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/event"
	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/pkg/github/webhook"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// Webhook is emitted for every verified GitHub App delivery.
var Webhook = event.New[webhook.Delivery]("github.webhook")

func init() {
	Webhook.On(dump)
	Webhook.On(installation)
	Webhook.On(installationRepositories)
	Webhook.On(issues)
	Webhook.On(issueComment)
	Webhook.On(detectConfChanges)
}

func dump(_ context.Context, d webhook.Delivery) {
	err := os.MkdirAll("events", 0o755)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GitHub webhook events dir")
		return
	}

	path := filepath.Join(
		"events",
		fmt.Sprintf("gh-%s-%s.json", d.Event, d.ID),
	)

	var pretty bytes.Buffer
	err = json.Indent(&pretty, d.Body, "", "  ")
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to format GitHub webhook")
		return
	}

	pretty.WriteByte('\n')
	err = os.WriteFile(path, pretty.Bytes(), 0o644)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to dump GitHub webhook")
		return
	}

	log.Info().
		Str("event", d.Event).
		Str("deliveryId", d.ID).
		Str("path", path).
		Msg("GitHub webhook dumped")
}

func installation(_ context.Context, d webhook.Delivery) {
	if d.Event != "installation" {
		return
	}

	var payload webhook.InstallationEvent
	json.Unmarshal(d.Body, &payload)

	i := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewRepositoriesRepository(db.GetDB()),
	)

	if payload.Action == "deleted" {
		err := i.Delete(payload.Installation.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to persist GitHub installation")
			return
		}

		log.Info().
			Str("deliveryId", d.ID).
			Str("action", payload.Action).
			Int64("githubId", payload.Installation.ID).
			Msg("GitHub installation webhook handled")

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
		return
	}

	log.Info().
		Str("deliveryId", d.ID).
		Str("action", payload.Action).
		Int64("githubId", payload.Installation.ID).
		Str("accountLogin", payload.Installation.Account.Login).
		Str("accountType", payload.Installation.Account.Type).
		Msg("GitHub installation webhook handled")
}

func installationRepositories(_ context.Context, d webhook.Delivery) {
	if d.Event != "installation_repositories" {
		return
	}

	var payload webhook.InstallationRepositoriesEvent
	json.Unmarshal(d.Body, &payload)

	i := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewRepositoriesRepository(db.GetDB()),
	)

	added := make([]app.Repository, len(payload.RepositoriesAdded))
	for n, repo := range payload.RepositoriesAdded {
		added[n] = app.Repository{
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
		return
	}

	log.Info().
		Str("deliveryId", d.ID).
		Str("action", payload.Action).
		Int64("githubId", payload.Installation.ID).
		Int("added", len(added)).
		Int("removed", len(removed)).
		Msg("GitHub installation repositories webhook handled")
}

func issues(_ context.Context, d webhook.Delivery) {
	if d.Event != "issues" {
		return
	}

	var payload webhook.IssueEvent
	json.Unmarshal(d.Body, &payload)

	switch payload.Action {
	case "opened", "edited", "labeled", "unlabeled":
	default:
		return
	}

	i := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewRepositoriesRepository(db.GetDB()),
	)

	installation, err := i.GetByGitHubId(payload.Installation.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to enqueue GitHub issue webhook")
		return
	}
	if lo.IsEmpty(lo.FromPtr(installation).WorkspaceId) {
		return
	}

	err = module.EnqueueTask(db.AsyncTaskTypeGitHubIssue, map[string]string{
		"deliveryId":     d.ID,
		"event":          d.Event,
		"action":         payload.Action,
		"body":           string(d.Body),
		"installationId": strconv.FormatInt(payload.Installation.ID, 10),
		"owner":          payload.Repository.Owner.Login,
		"repo":           payload.Repository.Name,
		"number":         strconv.Itoa(payload.Issue.Number),
	}, installation.WorkspaceId)

	if err != nil {
		log.Error().Err(err).Msg("Failed to enqueue GitHub issue webhook")
		return
	}

	log.Info().
		Str("deliveryId", d.ID).
		Str("action", payload.Action).
		Int64("githubId", payload.Installation.ID).
		Str("owner", payload.Repository.Owner.Login).
		Str("repo", payload.Repository.Name).
		Int("number", payload.Issue.Number).
		Msg("GitHub issue webhook handled")
}

func issueComment(_ context.Context, d webhook.Delivery) {
	if d.Event != "issue_comment" {
		return
	}

	var payload webhook.IssueCommentEvent
	json.Unmarshal(d.Body, &payload)

	if payload.Action != "created" {
		return
	}

	i := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewRepositoriesRepository(db.GetDB()),
	)

	installation, err := i.GetByGitHubId(payload.Installation.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to enqueue GitHub issue comment webhook")
		return
	}
	if lo.IsEmpty(lo.FromPtr(installation).WorkspaceId) {
		return
	}

	err = module.EnqueueTask(db.AsyncTaskTypeGitHubIssue, map[string]string{
		"deliveryId":     d.ID,
		"event":          d.Event,
		"action":         payload.Action,
		"body":           string(d.Body),
		"installationId": strconv.FormatInt(payload.Installation.ID, 10),
		"owner":          payload.Repository.Owner.Login,
		"repo":           payload.Repository.Name,
		"number":         strconv.Itoa(payload.Issue.Number),
	}, installation.WorkspaceId)

	if err != nil {
		log.Error().Err(err).Msg("Failed to enqueue GitHub issue comment webhook")
		return
	}

	log.Info().
		Str("deliveryId", d.ID).
		Str("action", payload.Action).
		Int64("githubId", payload.Installation.ID).
		Str("owner", payload.Repository.Owner.Login).
		Str("repo", payload.Repository.Name).
		Int("number", payload.Issue.Number).
		Msg("GitHub issue comment webhook handled")
}

func detectConfChanges(_ context.Context, d webhook.Delivery) {
	if d.Event != "push" {
		return
	}

	var payload webhook.PushEvent
	json.Unmarshal(d.Body, &payload)

	if payload.Ref != fmt.Sprintf("refs/heads/%s", payload.Repository.DefaultBranch) {
		return
	}

	path := payload.ChangedPath(".ziee.yml", ".ziee.yaml")
	if path == "" {
		return
	}

	i := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewRepositoriesRepository(db.GetDB()),
	)

	installation, err := i.GetByGitHubId(payload.Installation.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to enqueue GitHub push webhook")
		return
	}
	if lo.IsEmpty(lo.FromPtr(installation).WorkspaceId) {
		return
	}

	err = i.SetConfigPath(payload.Repository.ID, path)
	if err != nil {
		log.Error().Err(err).Msg("Failed to store GitHub config path")
		return
	}

	err = module.EnqueueTask(db.AsyncTaskTypeRepoLabels, map[string]string{
		"deliveryId":     d.ID,
		"installationId": strconv.FormatInt(payload.Installation.ID, 10),
		"fullName":       payload.Repository.FullName,
		"name":           payload.Repository.Name,
		"path":           path,
		"sha":            payload.After,
	}, installation.WorkspaceId)

	if err != nil {
		log.Error().Err(err).Msg("Failed to enqueue GitHub push webhook")
		return
	}

	log.Info().
		Str("deliveryId", d.ID).
		Int64("githubId", payload.Installation.ID).
		Str("owner", payload.Repository.Owner.Login).
		Str("repo", payload.Repository.Name).
		Str("path", path).
		Str("ref", payload.Ref).
		Msg("GitHub push webhook handled")
}
