// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package queue

import (
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// QueueActions builds enter-queue actions for a `@zieeai queue` command.
func QueueActions(conf *v1.File, event eval.Event, client eval.Client, args []string) []action.Action {
	rule, ok := SelectQueueRule(conf.MergeQueue.QueueRules, event, client, args, conf.MergeQueue.Labels.Queued)
	if !ok {
		return nil
	}

	return actionsForRule(conf, event, client, rule)
}

// EnterActions builds enter-queue actions when queue_when matches on a PR event.
func EnterActions(conf *v1.File, event eval.Event, client eval.Client) []action.Action {
	rule, ok := SelectEntryRule(conf.MergeQueue.QueueRules, event.Issue, client, conf.MergeQueue.Labels.Queued)
	if !ok {
		return nil
	}

	return actionsForRule(conf, event, client, rule)
}

// RefreshActions re-evaluates queue_when for an already-queued PR (`@zieeai refresh`).
func RefreshActions(conf *v1.File, event eval.Event, client eval.Client) []action.Action {
	if !eval.ContainsFold(event.Issue.Labels, conf.MergeQueue.Labels.Queued) &&
		!eval.ContainsFold(event.Issue.Labels, conf.MergeQueue.Labels.Checking) {
		return nil
	}

	rule, ok := SelectStayRule(conf.MergeQueue.QueueRules, event.Issue, client)
	if !ok {
		return eval.DequeueActions(conf)
	}

	priority := SelectPriority(conf.MergeQueue.PriorityRules, event.Issue, client)

	return []action.Action{{
		Kind:     policy.Queue,
		Priority: priority,
		Rule:     rule.Name,
	}}
}

func actionsForRule(conf *v1.File, event eval.Event, client eval.Client, rule v1.QueueRule) []action.Action {
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
