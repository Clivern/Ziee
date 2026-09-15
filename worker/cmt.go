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

// HandleGitHubComment evaluates an issue comment webhook.
func (h *handlers) HandleGitHubComment(ctx context.Context, msg *broker.Msg) error {
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

	var comment webhook.IssueCommentEvent
	err = json.Unmarshal([]byte(payload["body"]), &comment)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	repos := module.NewRepository(
		db.NewRepositoriesRepository(db.GetDB()),
		db.NewRepositoryMetaRepository(db.GetDB()),
	)

	path, err := repos.GetConfigPath(comment.Repository.ID)
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
			Association: comment.Issue.AuthorAssociation,
			Labels:      labels,
			Assignees:   assignees,
		},
		Comment: comment.Comment.Body,
		Actor: eval.Actor{
			Login: comment.Sender.Login,
		},
	}, NewIssueClient(ctx, installationId, payload["owner"], payload["repo"]))

	log.Info().
		Str("deliveryId", payload["deliveryId"]).
		Str("owner", payload["owner"]).
		Str("repo", payload["repo"]).
		Str("number", payload["number"]).
		Str("action", payload["action"]).
		Interface("actions", plan.Actions).
		Msg("Comment plan")

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
