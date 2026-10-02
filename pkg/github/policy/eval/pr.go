// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"path"
	"strings"

	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// EvaluatePROpened evaluates a new pull request against every pr_triage rule.
func EvaluatePROpened(conf *v1.File, event Event, client Client) action.Plan {
	var plan action.Plan

	if !conf.PRTriage.Enabled {
		return plan
	}

	if strings.ToLower(event.Account.Type) == "organization" {
		event.Org = event.Account.Login
	}

	event.Issue.Teams = MergeTeams(
		event.Issue.Teams,
		GetTeamsFromFile(conf.Teams, event.Issue.Author),
		client.GetTeams(event.Org, event.Issue.Author),
	)

	intentions := GetIntentionsFromRules(conf.PRTriage.Rules)
	ai := conf.PRTriage.AI
	if ai.Enabled && len(intentions) > 0 &&
		!SkipAI(conf.PRTriage.Rules, event.Issue, client) &&
		!SkipAIQuota(ai, event.Issue, client, true) {
		event.Issue.Intention = client.EvaluateIssue(event.Issue, intentions)
	}

	for _, rule := range conf.PRTriage.Rules {
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
		if len(rule.Reviewers) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:  policy.RequestReviewers,
				Users: rule.Reviewers,
			})
		}
		if len(rule.ReviewTeams) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:  policy.RequestReviewTeams,
				Users: rule.ReviewTeams,
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

	body := OutcomeComment(conf.PRTriage.Comments, "pull request", plan.Actions)
	if !lo.IsEmpty(body) {
		plan.Actions = append(plan.Actions, action.Action{
			Kind: policy.Comment,
			Body: body,
		})
	}

	return plan
}

// EvaluatePRComment evaluates a pull request comment as a Ziee command.
func EvaluatePRComment(conf *v1.File, event Event, client Client) action.Plan {
	var plan action.Plan

	if !conf.PRTriage.Enabled || IsAppActor(event.Actor.Login) {
		return plan
	}

	cmd := ParseCommand(event.Comment)
	if lo.IsEmpty(cmd.Verb) {
		return plan
	}

	command, ok := conf.PRTriage.Commands[cmd.Verb]
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
	case "label":
		if len(cmd.Args) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:   policy.AddLabels,
				Labels: cmd.Args,
			})
		}
	case "unlabel":
		if len(cmd.Args) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:   policy.RemoveLabels,
				Labels: cmd.Args,
			})
		}
	case "assign":
		if len(cmd.Args) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:  policy.Assign,
				Users: cmd.Args,
			})
		}
	case "unassign":
		if len(cmd.Args) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:  policy.Unassign,
				Users: cmd.Args,
			})
		}
	case "reviewers":
		if len(cmd.Args) > 0 {
			plan.Actions = append(plan.Actions, action.Action{
				Kind:  policy.RequestReviewers,
				Users: cmd.Args,
			})
		}
	case "close":
		plan.Actions = append(plan.Actions, action.Action{Kind: policy.Close})
	case "reopen":
		plan.Actions = append(plan.Actions, action.Action{Kind: policy.Reopen})
	case "spam":
		plan.Actions = append(plan.Actions, action.Action{
			Kind:   policy.AddLabels,
			Labels: []string{"spam"},
		})
		plan.Actions = append(plan.Actions, action.Action{Kind: policy.Close})
		client.BlockAuthor(event.Issue)
		plan.Actions = append(plan.Actions, action.Action{
			Kind:  policy.BlockAuthor,
			Users: []string{event.Issue.Author},
		})
	}

	body := CommandOutcomeComment(conf.PRTriage.Comments, "pull request", event.Actor.Login, plan.Actions)
	if !lo.IsEmpty(body) {
		plan.Actions = append(plan.Actions, action.Action{
			Kind: policy.Comment,
			Body: body,
		})
	}

	return plan
}

// MatchAnyFile reports whether any file matches any glob-like pattern.
func MatchAnyFile(patterns, files []string) bool {
	for _, file := range files {
		for _, pattern := range patterns {
			if MatchFile(pattern, file) {
				return true
			}
		}
	}

	return false
}

// MatchFile matches a single path against a simple glob (`*`, `**`).
func MatchFile(pattern, file string) bool {
	pattern = strings.ReplaceAll(pattern, "**/", "")
	pattern = strings.ReplaceAll(pattern, "**", "*")
	matched, err := path.Match(pattern, file)
	if err != nil {
		log.Error().
			Err(err).
			Str("pattern", pattern).
			Str("file", file).
			Msg("Failed to match file")
	}
	if matched {
		return true
	}

	matched, err = path.Match(pattern, path.Base(file))
	if err != nil {
		log.Error().
			Err(err).
			Str("pattern", pattern).
			Str("file", file).
			Msg("Failed to match file")
	}

	return matched
}
