// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"path"
	"strings"

	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// EvaluatePROpened evaluates a new pull request against every pr_triage rule.
func EvaluatePROpened(conf *v1.File, event Event, client Client) action.Plan {
	var plan action.Plan

	if strings.ToLower(event.Account.Type) == "organization" {
		event.Org = event.Account.Login
	}

	event.Issue.Teams = MergeTeams(
		event.Issue.Teams,
		GetTeamsFromFile(conf.Teams, event.Issue.Author),
		client.GetTeams(event.Org, event.Issue.Author),
	)

	intentions := GetIntentionsFromRules(conf.MergeQueue.PRTriage.Rules)
	if conf.MergeQueue.PRTriage.AI.Enabled && len(intentions) > 0 && !SkipAI(conf.MergeQueue.PRTriage.Rules, event.Issue, client) {
		event.Issue.Intention = client.EvaluateIssue(event.Issue, intentions)
	}

	for _, rule := range conf.MergeQueue.PRTriage.Rules {
		matched := true
		for _, when := range rule.When {
			if len(when.Files) > 0 && !MatchAnyFile(when.Files, event.Issue.Files) {
				matched = false
			}
			if when.MaxFilesChanged != nil && len(event.Issue.Files) > *when.MaxFilesChanged {
				matched = false
			}
			if when.MinFilesChanged != nil && len(event.Issue.Files) < *when.MinFilesChanged {
				matched = false
			}
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
			if !lo.IsEmpty(when.Intention.Name) && event.Issue.Intention.Name != when.Intention.Name {
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
			if when.FirstContribution != nil && *when.FirstContribution != client.IsFirstContribution(event.Issue) {
				matched = false
			}
			if when.MaxPrsOpened != nil && !client.PrsOpenedExceeds(
				event.Issue,
				when.MaxPrsOpened.Count,
				when.MaxPrsOpened.Within,
			) {
				matched = false
			}
			if when.AuthorBlocked != nil && *when.AuthorBlocked != client.IsAuthorBlocked(event.Issue) {
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
			client.BlockAuthor(event.Issue.Author)
			plan.Actions = append(plan.Actions, action.Action{
				Kind:  policy.BlockAuthor,
				Users: []string{event.Issue.Author},
			})
		}
	}

	body := OutcomeComment(conf.MergeQueue.Comments, plan.Actions)
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
	matched, _ := path.Match(pattern, file)
	if matched {
		return true
	}

	matched, _ = path.Match(pattern, path.Base(file))

	return matched
}
