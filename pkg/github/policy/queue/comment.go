// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package queue

import (
	"strings"

	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// EvaluateComment evaluates a merge-queue command on a pull request.
func EvaluateComment(conf *v1.File, event eval.Event, client eval.Client) action.Plan {
	var plan action.Plan

	if eval.IsAppActor(event.Actor.Login) {
		return plan
	}

	cmd := eval.ParseCommand(event.Comment)
	command := conf.MergeQueue.Commands[cmd.Verb]

	if strings.ToLower(event.Account.Type) == "organization" {
		event.Org = event.Account.Login
	}

	event.Actor.Permission = client.GetPermission(event.Actor.Login)
	event.Actor.Teams = eval.MergeTeams(
		event.Actor.Teams,
		eval.GetTeamsFromFile(conf.Teams, event.Actor.Login),
		client.GetTeams(event.Org, event.Actor.Login),
	)

	if !eval.MatchAllow(command.Allow, event.Actor, event.Issue) {
		return plan
	}

	switch cmd.Verb {
	case "queue":
		plan.Actions = append(plan.Actions, QueueActions(conf, event, client, cmd.Args)...)
	case "dequeue":
		plan.Actions = append(plan.Actions, eval.DequeueActions(conf)...)
	case "requeue":
		plan.Actions = append(plan.Actions, eval.DequeueActions(conf)...)
		plan.Actions = append(plan.Actions, QueueActions(conf, event, client, cmd.Args)...)
	case "refresh":
		plan.Actions = append(plan.Actions, RefreshActions(conf, event, client)...)
	}

	body := AckComment(conf.MergeQueue.Comments, event.Actor.Login, plan.Actions)
	if !lo.IsEmpty(body) {
		plan.Actions = append(plan.Actions, action.Action{
			Kind: policy.Comment,
			Body: body,
		})
	}

	return plan
}

// AckComment returns the merge-queue command acknowledgement.
func AckComment(mode, actor string, actions []action.Action) string {
	if mode != policy.CommentsAll && mode != policy.CommentsOutcomes {
		return ""
	}

	var queued, dequeued bool
	for _, a := range actions {
		switch a.Kind {
		case policy.Queue:
			queued = true
		case policy.Dequeue:
			dequeued = true
		}
	}

	by := " as requested by @" + actor + "."
	if queued && dequeued {
		return "Requeued this pull request" + by
	}
	if queued {
		return "Queued this pull request" + by
	}
	if dequeued {
		return "Dequeued this pull request" + by
	}

	return ""
}
