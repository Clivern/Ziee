// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.
//
// Package queue evaluates merge-queue entry and merge readiness from `.ziee.yml`.
//
//	spec.Parse → queue.Run → action.Apply
package queue

import (
	"github.com/clivern/ziee/pkg/github/policy"
	"github.com/clivern/ziee/pkg/github/policy/action"
	"github.com/clivern/ziee/pkg/github/policy/eval"
	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

// Run evaluates merge-queue commands and PR enter / merge readiness.
// Triage lives in eval.Run — callers run that separately.
func Run(conf *v1.File, event eval.Event, client eval.Client) action.Plan {
	if !conf.MergeQueue.Enabled {
		return action.Plan{}
	}

	switch event.Kind {
	case policy.KindPullRequestComment:
		cmd := eval.ParseCommand(event.Comment)
		if lo.IsEmpty(cmd.Verb) {
			return action.Plan{}
		}
		if _, ok := conf.MergeQueue.Commands[cmd.Verb]; ok {
			return EvaluateComment(conf, event, client)
		}
		return action.Plan{}
	case policy.KindPullRequestOpened, policy.KindPullRequestEdited, policy.KindPullRequestSynchronize,
		policy.KindCheckRunCompleted:
		return EvaluatePR(conf, event, client)
	default:
		return action.Plan{}
	}
}
