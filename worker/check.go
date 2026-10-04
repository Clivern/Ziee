// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/broker"
	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	"github.com/clivern/ziee/pkg/github/policy/queue"
	"github.com/clivern/ziee/pkg/github/policy/spec"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"
	"github.com/clivern/ziee/pkg/github/webhook"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// HandleGitHubCheckRun re-evaluates merge-queue entry/merge when a check completes.
func (h *handlers) HandleGitHubCheckRun(ctx context.Context, msg *broker.Msg) error {
	var payload map[string]string
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return ErrInvalidPayload
	}

	taskId := db.Id(payload["taskId"])

	h.Tasks.MarkRunning(taskId)

	installationId, err := strconv.ParseInt(payload["installationId"], 10, 64)
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	var check webhook.CheckRunEvent
	err = json.Unmarshal([]byte(payload["body"]), &check)
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	path, err := h.Repository.GetConfigPath(check.Repository.ID)
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	data, err := app.Get().GetFile(ctx, installationId, payload["owner"], payload["repo"], path)
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	file, err := spec.Parse(data)
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	if !file.MergeQueue.Enabled {
		return h.Tasks.Complete(taskId, `{"skipped":"merge_queue_disabled"}`)
	}

	numbers := lo.Map(check.CheckRun.PullRequests, func(pr webhook.CheckRunPullRef, _ int) int {
		return pr.Number
	})
	if len(numbers) == 0 {
		pulls, err := app.Get().ListPullRequestsForCommit(
			ctx,
			installationId,
			payload["owner"],
			payload["repo"],
			check.CheckRun.HeadSHA,
		)
		if err != nil {
			h.Tasks.Fail(taskId, err.Error())
			return err
		}
		numbers = lo.Map(pulls, func(pr app.PullRequest, _ int) int {
			return pr.Number
		})
	}

	for _, number := range numbers {
		err = h.evaluateCheckRunPR(ctx, file, installationId, payload, check, number)
		if err != nil {
			h.Tasks.Fail(taskId, err.Error())
			return err
		}
	}

	result, err := json.Marshal(map[string]any{
		"deliveryId": payload["deliveryId"],
		"event":      payload["event"],
		"action":     payload["action"],
		"check":      check.CheckRun.Name,
		"numbers":    numbers,
	})
	if err != nil {
		h.Tasks.Fail(taskId, err.Error())
		return err
	}

	return h.Tasks.Complete(taskId, string(result))
}

func (h *handlers) evaluateCheckRunPR(
	ctx context.Context,
	file *v1.File,
	installationId int64,
	payload map[string]string,
	check webhook.CheckRunEvent,
	number int,
) error {
	pull, err := app.Get().GetPullRequest(ctx, installationId, payload["owner"], payload["repo"], number)
	if err != nil {
		return err
	}

	files, err := app.Get().ListPullRequestFiles(ctx, installationId, payload["owner"], payload["repo"], number)
	if err != nil {
		return err
	}

	labels := make([]string, len(pull.Labels))
	for i, label := range pull.Labels {
		labels[i] = label.Name
	}

	assignees := make([]string, len(pull.Assignees))
	for i, user := range pull.Assignees {
		assignees[i] = user.Login
	}

	event := eval.Event{
		Kind: policy.KindCheckRunCompleted,
		Account: eval.Account{
			Login: check.Repository.Owner.Login,
			Type:  check.Repository.Owner.Type,
		},
		Issue: eval.Issue{
			Number:      pull.Number,
			Title:       pull.Title,
			Body:        pull.Body,
			Author:      pull.User.Login,
			AuthorId:    pull.User.ID,
			Association: pull.AuthorAssociation,
			Labels:      labels,
			Assignees:   assignees,
			Files:       files,
			Draft:       pull.Draft,
			Conflict:    pull.Mergeable != nil && !*pull.Mergeable,
			Closed:      pull.State == "closed",
		},
		Actor: eval.Actor{
			Login: check.Sender.Login,
		},
	}
	client := NewPullRequestClient(ctx, installationId, check.Repository.ID, payload["owner"], payload["repo"])

	plan := queue.Run(file, event, client)

	log.Info().
		Str("deliveryId", payload["deliveryId"]).
		Str("owner", payload["owner"]).
		Str("repo", payload["repo"]).
		Int("number", number).
		Str("check", check.CheckRun.Name).
		Str("conclusion", check.CheckRun.Conclusion).
		Interface("actions", plan.Actions).
		Msg("Check run queue plan")

	return action.Apply(ctx, action.NewClient(
		installationId,
		check.Repository.ID,
		pull.User.ID,
	), action.Repo{
		Owner:  payload["owner"],
		Name:   payload["repo"],
		Number: number,
	}, plan)
}
