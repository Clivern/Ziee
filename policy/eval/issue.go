// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"regexp"
	"slices"

	"github.com/clivern/ziee/policy"
	"github.com/clivern/ziee/policy/action"
	v1 "github.com/clivern/ziee/policy/spec/v1"
)

// EvaluateIssueOpened evaluates a new issue against every issue rule.
func EvaluateIssueOpened(file *v1.File, event Event, client Client) action.Plan {
	var plan action.Plan

	if !file.IssueTriage.Enabled {
		return plan
	}

	// TODO: Use client to enrich author teams and AI intention.
	for _, rule := range file.IssueTriage.Rules {
		matched := true

		for _, when := range rule.When {
			titleMatches, _ := regexp.MatchString(when.Title, event.Issue.Title)
			bodyMatches, _ := regexp.MatchString(when.Body, event.Issue.Body)

			if when.Title != "" && !titleMatches {
				matched = false
			}
			if when.Body != "" && !bodyMatches {
				matched = false
			}
			if len(when.AuthorIn) > 0 && !slices.Contains(when.AuthorIn, event.Issue.Author) {
				matched = false
			}
			if slices.Contains(when.AuthorNotIn, event.Issue.Author) {
				matched = false
			}
			if when.Intention != "" && when.Intention != event.Issue.Intention {
				matched = false
			}
			if when.Label != "" && !slices.Contains(event.Issue.Labels, when.Label) {
				matched = false
			}

			if len(when.AuthorInTeam) > 0 {
				authorInTeam := false
				for _, team := range when.AuthorInTeam {
					if slices.Contains(event.Issue.Teams, team) {
						authorInTeam = true
					}
				}
				if !authorInTeam {
					matched = false
				}
			}
			for _, team := range when.AuthorNotInTeam {
				if slices.Contains(event.Issue.Teams, team) {
					matched = false
				}
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
	}

	// TODO: Add an outcome comment according to file.IssueTriage.Comments.

	return plan
}

// EvaluateIssueEdited re-evaluates an edited issue against every issue rule.
func EvaluateIssueEdited(file *v1.File, event Event, client Client) action.Plan {
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
func EvaluateIssueLabelChange(file *v1.File, event Event, client Client) action.Plan {
	// 1. Ignore the event when event.Actor is the Ziee GitHub App.
	// 2. Evaluate every issue rule against event.Issue.Labels, which is the current label set.
	// 3. Compare the matching actions with the issue's current labels and assignees.
	// 4. Add newly required actions and remove actions whose rules no longer match.
	// 5. Return an idempotent plan so webhook redelivery produces no duplicate changes.

	return action.Plan{}
}

// EvaluateIssueComment evaluates an issue comment as a Ziee command.
func EvaluateIssueComment(file *v1.File, event Event, client Client) action.Plan {
	// 1. Ignore the event when event.Actor is the Ziee GitHub App.
	// 2. Ignore comments that do not start with `@ziee`.
	// 3. Parse the command verb and arguments from event.Comment.
	// 4. Find the verb in file.IssueTriage.Commands.
	// 5. Use client to load the actor's repository permission and teams.
	// 6. Allow the command when any permission, team, or user entry matches.
	// 7. Convert the command and arguments into label, assignment, or state actions.
	// 8. Add an outcome comment action when the configured comment mode requires it.
	// 9. Return the complete ordered action plan.

	return action.Plan{}
}
