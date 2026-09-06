// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.
//
// Package eval turns a parsed `.ziee.yml` and an event into an action plan.
//
//	spec.Parse → eval.Run → action.Apply
package eval

import (
	"github.com/clivern/ziee/policy/action"
	v1 "github.com/clivern/ziee/policy/spec/v1"
)

// Run evaluates a parsed spec against an event.
func Run(file *v1.File, event Event, client Client) action.Plan {
	return action.Plan{}
}
