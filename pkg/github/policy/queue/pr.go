// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package queue

import (
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"
)

// EvaluatePR decides whether a PR should enter the queue, stay, or be merge-ready.
func EvaluatePR(conf *v1.File, event eval.Event, client eval.Client) action.Plan {
	var plan action.Plan

	if event.Issue.Closed {
		return plan
	}

	labels := conf.MergeQueue.Labels
	inQueue := eval.ContainsFold(event.Issue.Labels, labels.Queued) ||
		eval.ContainsFold(event.Issue.Labels, labels.Checking)

	if event.Issue.Draft {
		if inQueue {
			plan.Actions = append(plan.Actions, eval.DequeueActions(conf)...)
		}

		return plan
	}

	if inQueue {
		plan.Actions = append(plan.Actions, EvaluateQueued(conf, event, client)...)

		return plan
	}

	plan.Actions = append(plan.Actions, EnterActions(conf, event, client)...)

	return plan
}

// EvaluateQueued re-checks queue_when / merge_when for a PR already in the line.
func EvaluateQueued(conf *v1.File, event eval.Event, client eval.Client) []action.Action {
	rule, ok := SelectStayRule(conf.MergeQueue.QueueRules, event.Issue, client)
	if !ok {
		return eval.DequeueActions(conf)
	}

	if len(rule.MergeWhen) == 0 || !eval.MatchClauses(rule.MergeWhen, event.Issue, client) {
		return nil
	}

	priority := SelectPriority(conf.MergeQueue.PriorityRules, event.Issue, client)

	return []action.Action{{
		Kind:     policy.Queue,
		Priority: priority,
		Rule:     rule.Name,
	}}
}
