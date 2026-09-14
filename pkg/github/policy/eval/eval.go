// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.
//
// Package eval turns a parsed `.ziee.yml` and an event into an action plan.
//
//	spec.Parse → eval.Run → action.Apply
package eval

import (
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"
)

// Run evaluates a parsed spec against an event.
func Run(conf *v1.File, event Event, client Client) action.Plan {
	switch event.Kind {
	case policy.KindIssueOpened:
		return EvaluateIssueOpened(conf, event, client)
	case policy.KindIssueEdited:
		return EvaluateIssueEdited(conf, event, client)
	case policy.KindIssueLabeled, policy.KindIssueUnlabeled:
		return EvaluateIssueLabelChange(conf, event, client)
	case policy.KindComment:
		return EvaluateIssueComment(conf, event, client)
	default:
		return action.Plan{}
	}
}
