// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"regexp"
	"strings"

	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// GetIntentionsFromRules returns unique intentions from triage rules.
func GetIntentionsFromRules(rules []v1.Rule) []v1.Intention {
	var intentions []v1.Intention

	for _, rule := range rules {
		for _, when := range rule.When {
			if !lo.IsEmpty(when.Intention.Name) {
				intentions = append(intentions, when.Intention)
			}
		}
	}

	return lo.UniqBy(intentions, func(intention v1.Intention) string {
		return intention.Name
	})
}

// MatchPattern reports whether value matches pattern, ignoring case.
func MatchPattern(pattern, value string) bool {
	matched, _ := regexp.MatchString("(?i)"+pattern, value)
	return matched
}

// GetTeamsFromFile returns `.ziee.yml` team names that include login.
func GetTeamsFromFile(teams []v1.Team, login string) []string {
	var names []string

	for _, team := range teams {
		if lo.ContainsBy(team.Members, func(member string) bool {
			return strings.EqualFold(member, login)
		}) {
			names = append(names, team.Name)
		}
	}

	return names
}

// MergeTeams returns unique team names from the given lists.
func MergeTeams(groups ...[]string) []string {
	return lo.Uniq(lo.Compact(lo.Flatten(groups)))
}

// OutcomeComment returns the triage outcome body for comments mode.
func OutcomeComment(mode string, actions []action.Action) string {
	if mode != policy.CommentsAll && mode != policy.CommentsOutcomes {
		return ""
	}

	var add, remove, assign, unassign []string
	var closed, reopened, commented bool

	for _, a := range actions {
		switch a.Kind {
		case policy.AddLabels:
			add = append(add, a.Labels...)
		case policy.RemoveLabels:
			remove = append(remove, a.Labels...)
		case policy.Assign:
			assign = append(assign, a.Users...)
		case policy.Unassign:
			unassign = append(unassign, a.Users...)
		case policy.Close:
			closed = true
		case policy.Reopen:
			reopened = true
		case policy.Comment:
			commented = true
		}
	}

	parts := lo.Compact([]string{
		JoinClause("labeled", "`", "`", add),
		JoinClause("removed", "`", "`", remove),
		JoinClause("assigned", "@", "", assign),
		JoinClause("unassigned", "@", "", unassign),
		lo.Ternary(closed, "closed this issue", ""),
		lo.Ternary(reopened, "reopened this issue", ""),
	})

	if len(parts) > 0 {
		return "Triaged this issue: " + strings.Join(parts, ", ") + "."
	}
	if commented {
		return "Triaged this issue."
	}
	if mode == policy.CommentsAll {
		return "Triage ran; no rules matched."
	}

	return ""
}

func JoinClause(verb, left, right string, items []string) string {
	items = lo.Uniq(items)
	if len(items) == 0 {
		return ""
	}

	return verb + " " + strings.Join(lo.Map(items, func(item string, _ int) string {
		return left + item + right
	}), ", ")
}
