// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"github.com/clivern/ziee/policy/action"
	v1 "github.com/clivern/ziee/policy/spec/v1"
)

// EvaluateIssueChange evaluates an issue change against every issue rule.
func EvaluateIssueChange(file *v1.File, event Event, client Client) action.Plan {
	// 1. Enrich the event with author teams and AI intention when configured.
	// 2. Visit every file.IssueTriage.Rules entry.
	// 3. Match every clause in rule.When against event.Issue.
	// 4. When all clauses match, append the rule's label and assignment actions.
	// 5. Add an outcome comment action when the configured comment mode requires it.
	// 6. Return the complete ordered action plan.

	return action.Plan{}
}
