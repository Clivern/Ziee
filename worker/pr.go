// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/broker"
	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	"github.com/clivern/ziee/pkg/github/policy/spec"
	"github.com/clivern/ziee/pkg/github/webhook"

	"github.com/rs/zerolog/log"
)

// HandleGitHubPullRequest evaluates a pull request webhook.
func (h *handlers) HandleGitHubPullRequest(ctx context.Context, msg *broker.Msg) error {
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

	var pull webhook.PullRequestEvent
	err = json.Unmarshal([]byte(payload["body"]), &pull)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	repos := module.NewRepository(
		db.NewRepositoriesRepository(db.GetDB()),
		db.NewRepositoryMetaRepository(db.GetDB()),
		db.NewRepositorySpamUserRepository(db.GetDB()),
	)

	path, err := repos.GetConfigPath(pull.Repository.ID)
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

	files, err := app.Get().ListPullRequestFiles(
		ctx,
		installationId,
		payload["owner"],
		payload["repo"],
		pull.PullRequest.Number,
	)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	labels := make([]string, len(pull.PullRequest.Labels))
	for i, label := range pull.PullRequest.Labels {
		labels[i] = label.Name
	}

	assignees := make([]string, len(pull.PullRequest.Assignees))
	for i, user := range pull.PullRequest.Assignees {
		assignees[i] = user.Login
	}

	plan := eval.Run(file, eval.Event{
		Kind: fmt.Sprintf("pull_request.%s", payload["action"]),
		Account: eval.Account{
			Login: pull.Repository.Owner.Login,
			Type:  pull.Repository.Owner.Type,
		},
		Issue: eval.Issue{
			Number:      pull.PullRequest.Number,
			Title:       pull.PullRequest.Title,
			Body:        pull.PullRequest.Body,
			Author:      pull.PullRequest.User.Login,
			AuthorId:    pull.PullRequest.User.ID,
			Association: pull.PullRequest.AuthorAssociation,
			Labels:      labels,
			Assignees:   assignees,
			Files:       files,
			Draft:       pull.PullRequest.Draft,
			Conflict:    pull.PullRequest.Mergeable != nil && !*pull.PullRequest.Mergeable,
			Closed:      pull.PullRequest.State == "closed",
		},
		Actor: eval.Actor{
			Login: pull.Sender.Login,
		},
	}, NewPullRequestClient(ctx, installationId, pull.Repository.ID, payload["owner"], payload["repo"]))

	log.Info().
		Str("deliveryId", payload["deliveryId"]).
		Str("owner", payload["owner"]).
		Str("repo", payload["repo"]).
		Str("number", payload["number"]).
		Str("action", payload["action"]).
		Interface("actions", plan.Actions).
		Msg("Pull request plan")

	err = action.Apply(ctx, action.NewClient(
		installationId,
		pull.Repository.ID,
		pull.PullRequest.User.ID,
	), action.Repo{
		Owner:  payload["owner"],
		Name:   payload["repo"],
		Number: pull.PullRequest.Number,
	}, plan)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

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

// HandleGitHubPullRequestComment evaluates an `@bot` command from a pull request comment.
func (h *handlers) HandleGitHubPullRequestComment(ctx context.Context, msg *broker.Msg) error {
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

	cmd := eval.ParseCommand(comment.Comment.Body)

	_, err = app.Get().CreateIssueCommentReaction(
		ctx,
		installationId,
		payload["owner"],
		payload["repo"],
		comment.Comment.ID,
		app.ReactionEyes,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("deliveryId", payload["deliveryId"]).
			Int64("commentId", comment.Comment.ID).
			Msg("Failed to react to pull request command comment")
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
		Kind: policy.KindPullRequestComment,
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
	}, NewPullRequestClient(ctx, installationId, comment.Repository.ID, payload["owner"], payload["repo"]))

	log.Info().
		Str("deliveryId", payload["deliveryId"]).
		Str("owner", payload["owner"]).
		Str("repo", payload["repo"]).
		Str("number", payload["number"]).
		Str("action", payload["action"]).
		Str("verb", cmd.Verb).
		Strs("args", cmd.Args).
		Interface("actions", plan.Actions).
		Msg("Pull request command plan")

	err = action.Apply(ctx, action.NewClient(
		installationId,
		comment.Repository.ID,
		comment.Issue.User.ID,
	), action.Repo{
		Owner:  payload["owner"],
		Name:   payload["repo"],
		Number: comment.Issue.Number,
	}, plan)
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	result, err := json.Marshal(map[string]any{
		"deliveryId": payload["deliveryId"],
		"event":      payload["event"],
		"action":     payload["action"],
		"owner":      payload["owner"],
		"repo":       payload["repo"],
		"number":     payload["number"],
		"verb":       cmd.Verb,
		"args":       cmd.Args,
		"actions":    plan.Actions,
	})
	if err != nil {
		h.tasks.Fail(taskId, err.Error())
		return err
	}

	return h.tasks.Complete(taskId, string(result))
}
