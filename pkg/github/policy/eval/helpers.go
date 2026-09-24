// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"regexp"
	"strings"

	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// GetIntentionsFromRules returns unique intentions from triage rules.
func GetIntentionsFromRules(rules []v1.Rule) []v1.Intention {
	var intentions []v1.Intention

	for _, rule := range rules {
		CollectIntentions(rule.When, &intentions)
	}

	return lo.UniqBy(intentions, func(intention v1.Intention) string {
		return intention.Name
	})
}

// CollectIntentions collects intentions from clauses.
func CollectIntentions(clauses v1.Clauses, intentions *[]v1.Intention) {
	for _, when := range clauses {
		if !lo.IsEmpty(when.Intention.Name) {
			*intentions = append(*intentions, when.Intention)
		}
		CollectIntentions(when.And, intentions)
		CollectIntentions(when.Or, intentions)
	}
}

// MatchPattern reports whether value matches pattern, ignoring case.
func MatchPattern(pattern, value string) bool {
	matched, err := regexp.MatchString("(?i)"+pattern, value)
	if err != nil {
		log.Error().
			Err(err).
			Str("pattern", pattern).
			Msg("Failed to match pattern")
	}

	return matched
}

// ContainsFold reports whether items includes value, ignoring case.
func ContainsFold(items []string, value string) bool {
	return lo.ContainsBy(items, func(item string) bool {
		return strings.EqualFold(item, value)
	})
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

// CommandOutcomeComment returns the command acknowledgement for comments mode.
func CommandOutcomeComment(mode, actor string, actions []action.Action) string {
	if mode != policy.CommentsAll && mode != policy.CommentsOutcomes {
		return ""
	}

	var add, remove, assign, unassign []string
	var closed, reopened, blocked bool

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
		case policy.BlockAuthor:
			blocked = true
		case policy.Comment:
		}
	}

	by := " as requested by @" + actor + "."
	if closed && len(add)+len(remove)+len(assign)+len(unassign)+lo.Ternary(reopened, 1, 0)+lo.Ternary(blocked, 1, 0) == 0 {
		return "Closed this issue" + by
	}
	if reopened && len(add)+len(remove)+len(assign)+len(unassign)+lo.Ternary(closed, 1, 0)+lo.Ternary(blocked, 1, 0) == 0 {
		return "Reopened this issue" + by
	}

	parts := lo.Compact([]string{
		JoinClause("Labeled", "`", "`", add),
		JoinClause("Removed", "`", "`", remove),
		JoinClause("Assigned", "@", "", assign),
		JoinClause("Unassigned", "@", "", unassign),
		lo.Ternary(closed, "closed this issue", ""),
		lo.Ternary(reopened, "reopened this issue", ""),
		lo.Ternary(blocked, "blocked the author", ""),
	})
	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ", ") + by
}

// IsAppActor reports whether login is the configured GitHub App bot.
func IsAppActor(login string) bool {
	bot := strings.ToLower(app.GetConfig().BotName)
	login = strings.ToLower(login)

	return login == bot || login == bot+"[bot]"
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

// SkipAI reports whether classification is skipped before triage rules run.
func SkipAI(_ []v1.Rule, issue Issue, client Client) bool {
	return client.IsAuthorBlocked(issue)
}

// SkipAIQuota reports whether ai.max_* skips classification (except overrides the cap).
func SkipAIQuota(ai v1.AI, issue Issue, client Client, pr bool) bool {
	if len(ai.Except) > 0 && MatchClauses(ai.Except, issue, client) {
		return false
	}
	if pr && ai.MaxPrsOpened != nil {
		return client.PrsOpenedExceeds(issue, ai.MaxPrsOpened.Count, ai.MaxPrsOpened.Within)
	}
	if !pr && ai.MaxIssuesOpened != nil {
		return client.IssuesOpenedExceeds(issue, ai.MaxIssuesOpened.Count, ai.MaxIssuesOpened.Within)
	}

	return false
}
