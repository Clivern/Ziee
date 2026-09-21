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
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	"github.com/clivern/ziee/pkg/github/policy/spec"
	"github.com/clivern/ziee/pkg/github/webhook"

	"github.com/rs/zerolog/log"
)

type commandPayload struct {
	TaskId         string   `json:"taskId"`
	DeliveryId     string   `json:"deliveryId"`
	Event          string   `json:"event"`
	Action         string   `json:"action"`
	Body           string   `json:"body"`
	InstallationId string   `json:"installationId"`
	Owner          string   `json:"owner"`
	Repo           string   `json:"repo"`
	Number         string   `json:"number"`
	Verb           string   `json:"verb"`
	Args           []string `json:"args"`
}

// HandleGitHubCommand evaluates an `@bot` command from an issue or PR comment.
func (h *handlers) HandleGitHubCommand(ctx context.Context, msg *broker.Msg) error {
	var payload commandPayload
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload.TaskId)

	h.tasks.MarkRunning(taskId)

	installationId, err := strconv.ParseInt(payload.InstallationId, 10, 64)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	var comment webhook.IssueCommentEvent
	err = json.Unmarshal([]byte(payload.Body), &comment)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	repos := module.NewRepository(
		db.NewRepositoriesRepository(db.GetDB()),
		db.NewRepositoryMetaRepository(db.GetDB()),
		db.NewRepositorySpamUserRepository(db.GetDB()),
	)

	path, err := repos.GetConfigPath(comment.Repository.ID)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	data, err := app.Get().GetFile(ctx, installationId, payload.Owner, payload.Repo, path)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	file, err := spec.Parse(data)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	labels := make([]string, len(comment.Issue.Labels))
	for i, label := range comment.Issue.Labels {
		labels[i] = label.Name
	}

	assignees := make([]string, len(comment.Issue.Assignees))
	for i, user := range comment.Issue.Assignees {
		assignees[i] = user.Login
	}

	plan := eval.Run(file, eval.Event{
		Kind: policy.KindComment,
		Account: eval.Account{
			Login: comment.Repository.Owner.Login,
			Type:  comment.Repository.Owner.Type,
		},
		Issue: eval.Issue{
			Number:      comment.Issue.Number,
			Title:       comment.Issue.Title,
			Body:        comment.Issue.Body,
			Author:      comment.Issue.User.Login,
			AuthorId:    comment.Issue.User.ID,
			Association: comment.Issue.AuthorAssociation,
			Labels:      labels,
			Assignees:   assignees,
			Closed:      comment.Issue.State == "closed",
		},
		Comment: comment.Comment.Body,
		Actor: eval.Actor{
			Login: comment.Sender.Login,
		},
	}, NewIssueClient(ctx, installationId, comment.Repository.ID, payload.Owner, payload.Repo))

	log.Info().
		Str("deliveryId", payload.DeliveryId).
		Str("owner", payload.Owner).
		Str("repo", payload.Repo).
		Str("number", payload.Number).
		Str("action", payload.Action).
		Str("verb", payload.Verb).
		Strs("args", payload.Args).
		Interface("actions", plan.Actions).
		Msg("Command plan")

	result, err := json.Marshal(map[string]any{
		"deliveryId": payload.DeliveryId,
		"event":      payload.Event,
		"action":     payload.Action,
		"owner":      payload.Owner,
		"repo":       payload.Repo,
		"number":     payload.Number,
		"verb":       payload.Verb,
		"args":       payload.Args,
		"actions":    plan.Actions,
	})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	return h.tasks.Complete(taskId, string(result))
}
