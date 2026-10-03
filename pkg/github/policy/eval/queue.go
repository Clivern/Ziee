// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"strings"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// EvaluateMergeQueueComment evaluates a merge-queue command on a pull request.
func EvaluateMergeQueueComment(conf *v1.File, event Event, client Client) action.Plan {
	var plan action.Plan

	if !conf.MergeQueue.Enabled || IsAppActor(event.Actor.Login) {
		return plan
	}

	cmd := ParseCommand(event.Comment)
	if lo.IsEmpty(cmd.Verb) {
		return plan
	}

	command, ok := conf.MergeQueue.Commands[cmd.Verb]
	if !ok {
		return plan
	}

	if strings.ToLower(event.Account.Type) == "organization" {
		event.Org = event.Account.Login
	}

	event.Actor.Permission = client.GetPermission(event.Actor.Login)
	event.Actor.Teams = MergeTeams(
		event.Actor.Teams,
		GetTeamsFromFile(conf.Teams, event.Actor.Login),
		client.GetTeams(event.Org, event.Actor.Login),
	)

	if !MatchAllow(command.Allow, event.Actor, event.Issue) {
		return plan
	}

	switch cmd.Verb {
	case "queue":
		plan.Actions = append(plan.Actions, queueActions(conf, event, client, cmd.Args)...)
	case "dequeue":
		plan.Actions = append(plan.Actions, dequeueActions(conf)...)
	case "requeue":
		plan.Actions = append(plan.Actions, dequeueActions(conf)...)
		plan.Actions = append(plan.Actions, queueActions(conf, event, client, cmd.Args)...)
	case "refresh":
		plan.Actions = append(plan.Actions, refreshActions(conf, event, client)...)
	}

	body := MergeQueueOutcomeComment(conf.MergeQueue.Comments, event.Actor.Login, plan.Actions)
	if !lo.IsEmpty(body) {
		plan.Actions = append(plan.Actions, action.Action{
			Kind: policy.Comment,
			Body: body,
		})
	}

	return plan
}

func queueActions(conf *v1.File, event Event, client Client, args []string) []action.Action {
	rule, ok := SelectQueueRule(conf.MergeQueue.QueueRules, event, client, args, conf.MergeQueue.Labels.Queued)
	if !ok {
		return nil
	}

	priority := SelectPriority(conf.MergeQueue.PriorityRules, event.Issue, client)
	labels := conf.MergeQueue.Labels

	actions := []action.Action{{
		Kind:     policy.Queue,
		Priority: priority,
		Rule:     rule.Name,
	}}

	add := lo.Compact([]string{labels.Queued})
	remove := lo.Compact([]string{labels.Dequeued, labels.Checking})
	if len(add) > 0 {
		actions = append(actions, action.Action{Kind: policy.AddLabels, Labels: add})
	}
	if len(remove) > 0 {
		actions = append(actions, action.Action{Kind: policy.RemoveLabels, Labels: remove})
	}

	return actions
}

func dequeueActions(conf *v1.File) []action.Action {
	labels := conf.MergeQueue.Labels

	actions := []action.Action{{Kind: policy.Dequeue}}

	add := lo.Compact([]string{labels.Dequeued})
	remove := lo.Compact([]string{labels.Queued, labels.Checking})
	if len(add) > 0 {
		actions = append(actions, action.Action{Kind: policy.AddLabels, Labels: add})
	}
	if len(remove) > 0 {
		actions = append(actions, action.Action{Kind: policy.RemoveLabels, Labels: remove})
	}

	return actions
}

func refreshActions(conf *v1.File, event Event, client Client) []action.Action {
	if !ContainsFold(event.Issue.Labels, conf.MergeQueue.Labels.Queued) &&
		!ContainsFold(event.Issue.Labels, conf.MergeQueue.Labels.Checking) {
		return nil
	}

	_, ok := SelectQueueRule(conf.MergeQueue.QueueRules, event, client, nil, "")
	if !ok {
		return dequeueActions(conf)
	}

	priority := SelectPriority(conf.MergeQueue.PriorityRules, event.Issue, client)

	return []action.Action{{
		Kind:     policy.Queue,
		Priority: priority,
	}}
}

// SelectQueueRule picks the queue rule for a queue/requeue command.
func SelectQueueRule(
	rules []v1.QueueRule,
	event Event,
	client Client,
	args []string,
	queuedLabel string,
) (v1.QueueRule, bool) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	for _, rule := range rules {
		if !lo.IsEmpty(name) && !strings.EqualFold(rule.Name, name) {
			continue
		}
		if len(rule.Allow) > 0 && !MatchAllow(rule.Allow, event.Actor, event.Issue) {
			continue
		}
		when := QueueWhenForEntry(rule.QueueWhen, queuedLabel)
		if !MatchClauses(when, event.Issue, client) {
			continue
		}

		return rule, true
	}

	return v1.QueueRule{}, false
}

// QueueWhenForEntry drops the queued-state label so entry is not chicken-and-egg.
func QueueWhenForEntry(clauses v1.Clauses, queuedLabel string) v1.Clauses {
	if lo.IsEmpty(queuedLabel) {
		return clauses
	}

	return lo.Filter(clauses, func(clause v1.Clause, _ int) bool {
		return lo.IsEmpty(clause.Label) || !strings.EqualFold(clause.Label, queuedLabel)
	})
}

// SelectPriority returns the first matching priority rule, or medium.
func SelectPriority(rules []v1.PriorityRule, issue Issue, client Client) string {
	for _, rule := range rules {
		if MatchClauses(rule.When, issue, client) {
			return lo.Ternary(lo.IsEmpty(rule.Priority), db.PQueuePriorityMedium, rule.Priority)
		}
	}

	return db.PQueuePriorityMedium
}
