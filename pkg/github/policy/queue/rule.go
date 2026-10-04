// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package queue

import (
	"strings"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// SelectQueueRule picks the first matching queue rule for a queue/requeue command.
func SelectQueueRule(
	rules []v1.QueueRule,
	event eval.Event,
	client eval.Client,
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
		if len(rule.Allow) > 0 && !eval.MatchAllow(rule.Allow, event.Actor, event.Issue) {
			continue
		}
		when := QueueWhenForEntry(rule.QueueWhen, queuedLabel)
		if !eval.MatchClauses(when, event.Issue, client) {
			continue
		}

		return rule, true
	}

	return v1.QueueRule{}, false
}

// SelectEntryRule picks the first rule whose queue_when matches for automatic entry.
func SelectEntryRule(
	rules []v1.QueueRule,
	issue eval.Issue,
	client eval.Client,
	queuedLabel string,
) (v1.QueueRule, bool) {
	for _, rule := range rules {
		when := QueueWhenForEntry(rule.QueueWhen, queuedLabel)
		if !eval.MatchClauses(when, issue, client) {
			continue
		}

		return rule, true
	}

	return v1.QueueRule{}, false
}

// SelectStayRule picks the first rule whose full queue_when still matches while queued.
func SelectStayRule(rules []v1.QueueRule, issue eval.Issue, client eval.Client) (v1.QueueRule, bool) {
	for _, rule := range rules {
		if !eval.MatchClauses(rule.QueueWhen, issue, client) {
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
func SelectPriority(rules []v1.PriorityRule, issue eval.Issue, client eval.Client) string {
	for _, rule := range rules {
		if eval.MatchClauses(rule.When, issue, client) {
			return lo.Ternary(lo.IsEmpty(rule.Priority), db.PQueuePriorityMedium, rule.Priority)
		}
	}

	return db.PQueuePriorityMedium
}
