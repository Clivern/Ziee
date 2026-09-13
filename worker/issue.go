// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/broker"
	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/pkg/github/policy/action"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	"github.com/clivern/ziee/pkg/github/policy/spec"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"
	"github.com/clivern/ziee/pkg/github/webhook"

	"github.com/rs/zerolog/log"
)

type noopClient struct{}

func (noopClient) GetTeams(string, string) []string { return nil }

func (noopClient) EvaluateIssue(eval.Issue, []v1.Intention) []string { return nil }

// HandleGitHubIssue evaluates an issue or issue comment webhook.
func (h *handlers) HandleGitHubIssue(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])
	h.tasks.MarkRunning(taskId)

	var plan action.Plan
	if payload["event"] != "issues" {
		result, err := json.Marshal(map[string]any{
			"deliveryId": payload["deliveryId"],
			"event":      payload["event"],
			"action":     payload["action"],
			"owner":      payload["owner"],
			"repo":       payload["repo"],
			"number":     payload["number"],
			"actions":    plan.Actions,
		})
		if err != nil {
			h.tasks.Fail(taskId, err.Error())
			return err
		}

		return h.tasks.Complete(taskId, string(result))
	}

	installationId, err := strconv.ParseInt(payload["installationId"], 10, 64)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	var issue webhook.IssueEvent
	err = json.Unmarshal([]byte(payload["body"]), &issue)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	repos := module.NewRepository(
		db.NewRepositoriesRepository(db.GetDB()),
		db.NewRepositoryMetaRepository(db.GetDB()),
	)

	path, err := repos.GetConfigPath(issue.Repository.ID)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	data, err := app.Get().GetFile(ctx, installationId, payload["owner"], payload["repo"], path)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	file, err := spec.Parse(data)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	labels := make([]string, len(issue.Issue.Labels))
	for i, label := range issue.Issue.Labels {
		labels[i] = label.Name
	}

	assignees := make([]string, len(issue.Issue.Assignees))
	for i, user := range issue.Issue.Assignees {
		assignees[i] = user.Login
	}

	plan = eval.Run(file, eval.Event{
		Kind: "issue." + payload["action"],
		Account: eval.Account{
			Login: issue.Repository.Owner.Login,
			Type:  issue.Repository.Owner.Type,
		},
		Issue: eval.Issue{
			Number:      issue.Issue.Number,
			Title:       issue.Issue.Title,
			Body:        issue.Issue.Body,
			Author:      issue.Issue.User.Login,
			Association: issue.Issue.AuthorAssociation,
			Labels:      labels,
			Assignees:   assignees,
		},
		Actor: eval.Actor{
			Login: issue.Sender.Login,
		},
	}, noopClient{})

	log.Info().
		Str("deliveryId", payload["deliveryId"]).
		Str("owner", payload["owner"]).
		Str("repo", payload["repo"]).
		Str("number", payload["number"]).
		Str("action", payload["action"]).
		Interface("actions", plan.Actions).
		Msg("Issue plan")

	result, err := json.Marshal(map[string]any{
		"deliveryId": payload["deliveryId"],
		"event":      payload["event"],
		"action":     payload["action"],
		"owner":      payload["owner"],
		"repo":       payload["repo"],
		"number":     payload["number"],
		"actions":    plan.Actions,
	})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	return h.tasks.Complete(taskId, string(result))
}
