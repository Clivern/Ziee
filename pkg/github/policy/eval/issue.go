// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// EvaluateIssueOpened evaluates a new issue against every issue rule.
func EvaluateIssueOpened(conf *v1.File, event Event, client Client) action.Plan {
	var plan action.Plan

	if !conf.IssueTriage.Enabled {
		return plan
	}

	event.Issue.Teams = MergeTeams(
		event.Issue.Teams,
		GetTeamsFromFile(conf.Teams, event.Issue.Author),
		client.GetTeams(event.Issue.Author),
	)

	intentions := GetIntentionsFromRules(conf.IssueTriage.Rules)
	if conf.IssueTriage.AI.Enabled && len(intentions) > 0 {
		event.Issue.Intentions = client.EvaluateIssue(event.Issue, intentions)
	}

	for _, rule := range conf.IssueTriage.Rules {
		matched := true

		for _, when := range rule.When {
			if !lo.IsEmpty(when.Title) && !MatchPattern(when.Title, event.Issue.Title) {
				matched = false
			}
			if !lo.IsEmpty(when.Body) && !MatchPattern(when.Body, event.Issue.Body) {
				matched = false
			}
			if len(when.AuthorIn) > 0 && !lo.Contains(when.AuthorIn, event.Issue.Author) {
				matched = false
			}
			if lo.Contains(when.AuthorNotIn, event.Issue.Author) {
				matched = false
			}
			if !lo.IsEmpty(when.Intention.Name) && !lo.Contains(event.Issue.Intentions, when.Intention.Name) {
				matched = false
			}
			if !lo.IsEmpty(when.Label) && !lo.Contains(event.Issue.Labels, when.Label) {
				matched = false
			}
			if len(when.AuthorInTeam) > 0 && !lo.Some(event.Issue.Teams, when.AuthorInTeam) {
				matched = false
			}
			if len(when.AuthorNotInTeam) > 0 && lo.Some(event.Issue.Teams, when.AuthorNotInTeam) {
				matched = false
			}
		}
		if !matched {
			continue
		}
		if len(rule.Labels.Add) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:   policy.AddLabels,
				Labels: rule.Labels.Add,
			})
		}
		if len(rule.Labels.Remove) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:   policy.RemoveLabels,
				Labels: rule.Labels.Remove,
			})
		}
		if len(rule.Assign) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:  policy.Assign,
				Users: rule.Assign,
			})
		}
		if !lo.IsEmpty(rule.Comment) {
			plan.Actions = append(plan.Actions, action.Action{
				Kind: policy.Comment,
				Body: rule.Comment,
			})
		}
	}

	body := OutcomeComment(conf.IssueTriage.Comments, plan.Actions)
	if !lo.IsEmpty(body) {
		plan.Actions = append(plan.Actions, action.Action{
			Kind: policy.Comment,
			Body: body,
		})
	}

	return plan
}

// EvaluateIssueEdited re-evaluates an edited issue against every issue rule.
func EvaluateIssueEdited(conf *v1.File, event Event, client Client) action.Plan {
	// 1. Ignore the event when event.Actor is the Ziee GitHub App.
	// 2. Enrich the event with author teams and AI intention.
	// 3. Evaluate every issue rule against the updated issue.
	// 4. Compare the matching actions with the issue's current labels and assignees.
	// 5. Add actions that are newly required.
	// 6. Remove managed labels or assignees whose rules no longer match.
	// 7. Add an outcome comment action when the configured comment mode requires it.
	// 8. Return the reconciliation plan.

	return action.Plan{}
}

// EvaluateIssueLabelChange re-evaluates an issue after a label is added or removed.
func EvaluateIssueLabelChange(conf *v1.File, event Event, client Client) action.Plan {
	// 1. Ignore the event when event.Actor is the Ziee GitHub App.
	// 2. Evaluate every issue rule against event.Issue.Labels, which is the current label set.
	// 3. Compare the matching actions with the issue's current labels and assignees.
	// 4. Add newly required actions and remove actions whose rules no longer match.
	// 5. Return an idempotent plan so webhook redelivery produces no duplicate changes.

	return action.Plan{}
}

// EvaluateIssueComment evaluates an issue comment as a Ziee command.
func EvaluateIssueComment(conf *v1.File, event Event, client Client) action.Plan {
	// 1. Ignore the event when event.Actor is the Ziee GitHub App.
	// 2. Ignore comments that do not start with `@ziee`.
	// 3. Parse the command verb and arguments from event.Comment.
	// 4. Find the verb in conf.IssueTriage.Commands.
	// 5. Use client to load the actor's repository permission and teams.
	// 6. Allow the command when any permission, team, or user entry matches.
	// 7. Convert the command and arguments into label, assignment, or state actions.
	// 8. Add an outcome comment action when the configured comment mode requires it.
	// 9. Return the complete ordered action plan.

	return action.Plan{}
}
