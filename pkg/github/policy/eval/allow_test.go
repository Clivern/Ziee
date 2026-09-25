// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"testing"

	v1 "github.com/clivern/ziee/pkg/github/policy/spec/v1"

	"github.com/stretchr/testify/assert"
)

func TestUnitMatchAllow(t *testing.T) {
	issue := Issue{Author: "author"}
	allow := v1.Allow{
		{Permission: "write"},
		{Teams: []string{"sre"}},
		{Users: []string{"maya"}},
		{Self: true},
	}

	assert.True(t, MatchAllow(allow, Actor{Login: "maya", Permission: "read"}, issue))
	assert.True(t, MatchAllow(allow, Actor{Login: "dev", Permission: "write"}, issue))
	assert.True(t, MatchAllow(allow, Actor{Login: "dev", Teams: []string{"sre"}}, issue))
	assert.True(t, MatchAllow(allow, Actor{Login: "author"}, issue))
	assert.False(t, MatchAllow(allow, Actor{Login: "guest", Permission: "read"}, issue))
}
