// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/broker"
	"github.com/clivern/ziee/pkg/github/app"
)

// HandleRepositoryBootstrap creates a setup issue and .ziee.yml PR for a repository.
func (h *handlers) HandleRepositoryBootstrap(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])
	h.tasks.MarkRunning(taskId)

	installationId, err := strconv.ParseInt(payload["installationId"], 10, 64)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	githubId, err := strconv.ParseInt(payload["githubId"], 10, 64)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	owner, _, _ := strings.Cut(payload["fullName"], "/")
	repo := payload["name"]

	exists, err := app.Get().FileExists(ctx, installationId, owner, repo, ".ziee.yaml", ".ziee.yml")
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}
	if exists {
		return h.tasks.Complete(taskId, "")
	}

	install := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewRepositoriesRepository(db.GetDB()),
	)

	issue, err := app.Get().CreateIssue(ctx, installationId, owner, repo, conf.SetupIssueTitle, conf.SetupIssueBody)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	meta, err := json.Marshal(map[string]int64{
		"id":     issue.ID,
		"number": int64(issue.Number),
	})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	err = install.SetRepoMeta(githubId, db.RepositoryMetaSetupIssue, string(meta))
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	pr, err := app.Get().CreatePullRequest(ctx, installationId, owner, repo, app.NewPullRequest{
		Branch:  "feat/ziee-init",
		Path:    ".ziee.yaml",
		Content: conf.DefaultZieeYML,
		Title:   conf.SetupPullRequestTitle,
		Body:    conf.SetupPullRequestBody,
	})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	prMeta, err := json.Marshal(map[string]int64{
		"id":     pr.ID,
		"number": int64(pr.Number),
	})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	err = install.SetRepoMeta(githubId, db.RepositoryMetaSetupPR, string(prMeta))
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	return h.tasks.Complete(taskId, "")
}
