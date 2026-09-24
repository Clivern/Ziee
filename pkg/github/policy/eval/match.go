// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// MatchClauses reports whether every clause matches.
func MatchClauses(clauses v1.Clauses, issue Issue, client Client) bool {
	for _, clause := range clauses {
		if !MatchClause(clause, issue, client) {
			return false
		}
	}

	return true
}

// MatchClause reports whether one when matcher matches.
func MatchClause(when v1.Clause, issue Issue, client Client) bool {
	if len(when.Files) > 0 && !MatchAnyFile(when.Files, issue.Files) {
		return false
	}
	if when.MaxFilesChanged != nil && len(issue.Files) > *when.MaxFilesChanged {
		return false
	}
	if when.MinFilesChanged != nil && len(issue.Files) < *when.MinFilesChanged {
		return false
	}
	if !lo.IsEmpty(when.Title) && !MatchPattern(when.Title, issue.Title) {
		return false
	}
	if !lo.IsEmpty(when.Body) && !MatchPattern(when.Body, issue.Body) {
		return false
	}
	if len(when.AuthorIn) > 0 && !ContainsFold(when.AuthorIn, issue.Author) {
		return false
	}
	if ContainsFold(when.AuthorNotIn, issue.Author) {
		return false
	}
	if !lo.IsEmpty(when.Intention.Name) && issue.Intention.Name != when.Intention.Name {
		return false
	}
	if !lo.IsEmpty(when.Label) && !lo.Contains(issue.Labels, when.Label) {
		return false
	}
	if len(when.AuthorInTeam) > 0 && !lo.Some(issue.Teams, when.AuthorInTeam) {
		return false
	}
	if len(when.AuthorNotInTeam) > 0 && lo.Some(issue.Teams, when.AuthorNotInTeam) {
		return false
	}
	if when.FirstContribution != nil && *when.FirstContribution != client.IsFirstContribution(issue) {
		return false
	}
	if when.MaxIssuesOpened != nil && !client.IssuesOpenedExceeds(
		issue,
		when.MaxIssuesOpened.Count,
		when.MaxIssuesOpened.Within,
	) {
		return false
	}
	if when.MaxPrsOpened != nil && !client.PrsOpenedExceeds(
		issue,
		when.MaxPrsOpened.Count,
		when.MaxPrsOpened.Within,
	) {
		return false
	}
	if when.AuthorBlocked != nil && *when.AuthorBlocked != client.IsAuthorBlocked(issue) {
		return false
	}
	if when.Draft != nil && *when.Draft != issue.Draft {
		return false
	}
	if when.Conflict != nil && *when.Conflict != issue.Conflict {
		return false
	}
	if when.Closed != nil && *when.Closed != issue.Closed {
		return false
	}
	if len(when.And) > 0 && !MatchClauses(when.And, issue, client) {
		return false
	}
	if len(when.Or) > 0 && !MatchAnyClause(when.Or, issue, client) {
		return false
	}

	return true
}

// MatchAnyClause reports whether any clause matches.
func MatchAnyClause(clauses v1.Clauses, issue Issue, client Client) bool {
	for _, clause := range clauses {
		if MatchClause(clause, issue, client) {
			return true
		}
	}

	return false
}
