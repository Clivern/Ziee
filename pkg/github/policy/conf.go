// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.
//
// Package policy is the framework that enforces `.ziee.yml`.
//
//	spec.Parse → eval.Run → action.Apply
package policy

const (
	KindIssueOpened    = "issue.opened"
	KindIssueEdited    = "issue.edited"
	KindIssueLabeled   = "issue.labeled"
	KindIssueUnlabeled = "issue.unlabeled"
	KindComment        = "issue.comment"

	CommentsOutcomes = "outcomes"
	CommentsAll      = "all"
	CommentsNone     = "none"

	AddLabels    = "add_labels"
	RemoveLabels = "remove_labels"
	Assign       = "assign"
	Unassign     = "unassign"
	Comment      = "comment"
	Close        = "close"
	Reopen       = "reopen"
)
