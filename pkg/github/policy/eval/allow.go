// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"strings"

	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/samber/lo"
)

var PermissionRank = map[string]int{
	"admin":    4,
	"maintain": 3,
	"write":    2,
	"triage":   1,
	"read":     0,
}

// MatchAllow reports whether actor may run a command from its allow list.
func MatchAllow(allow v1.Allow, actor Actor, issue Issue) bool {
	for _, entry := range allow {
		if !lo.IsEmpty(entry.Permission) &&
			PermissionRank[strings.ToLower(actor.Permission)] >= PermissionRank[strings.ToLower(entry.Permission)] {
			return true
		}
		if len(entry.Teams) > 0 && lo.Some(actor.Teams, entry.Teams) {
			return true
		}
		if len(entry.Users) > 0 && ContainsFold(entry.Users, actor.Login) {
			return true
		}
		if entry.Self && strings.EqualFold(actor.Login, issue.Author) {
			return true
		}
	}

	return false
}
