// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"strings"

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

	if strings.ToLower(event.Account.Type) == "organization" {
		event.Org = event.Account.Login
	}

	// TODO: Github Org Teams Fetching
	event.Issue.Teams = MergeTeams(
		event.Issue.Teams,
		GetTeamsFromFile(conf.Teams, event.Issue.Author),
		client.GetTeams(event.Org, event.Issue.Author),
	)

	intentions := GetIntentionsFromRules(conf.IssueTriage.Rules)
	ai := conf.IssueTriage.AI
	if ai.Enabled && len(intentions) > 0 &&
		!SkipAI(conf.IssueTriage.Rules, event.Issue, client) &&
		!SkipAIQuota(ai, event.Issue, client, false) {
		event.Issue.Intention = client.EvaluateIssue(event.Issue, intentions)
	}

	for _, rule := range conf.IssueTriage.Rules {
		if !MatchClauses(rule.When, event.Issue, client) {
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
		if rule.Close {
			plan.Actions = append(plan.Actions, action.Action{
				Kind: policy.Close,
			})
		}
		if rule.BlockAuthor {
			client.BlockAuthor(event.Issue)
			plan.Actions = append(plan.Actions, action.Action{
				Kind:  policy.BlockAuthor,
				Users: []string{event.Issue.Author},
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
	// 2. cmd := ParseCommand(event.Comment); ignore when Verb is empty.
	// 3. Find cmd.Verb in conf.IssueTriage.Commands (label, unlabel, assign, unassign, close, reopen, spam, summarize).
	// 4. Use client to load the actor's repository permission and teams.
	// 5. Allow the command when any permission, team, or user entry matches.
	// 6. Convert cmd.Verb and cmd.Args into label, assignment, or state actions.
	//    For `spam`: add spam label, close, block the issue author (BlockAuthor), and comment.
	// 7. Add an outcome comment action when the configured comment mode requires it.
	// 8. Return the complete ordered action plan.

	return action.Plan{}
}
