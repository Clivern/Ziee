// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitParseCommand(t *testing.T) {
	assert.Equal(t, Command{Verb: "label", Args: []string{"bug", "security"}}, ParseCommand("@zieeio label bug security"))
	assert.Equal(t, Command{Verb: "unlabel", Args: []string{"bot"}}, ParseCommand("@zieeio unlabel bot"))
	assert.Equal(t, Command{Verb: "assign", Args: []string{"clivern"}}, ParseCommand("@zieeio assign clivern"))
	assert.Equal(t, Command{Verb: "assign", Args: []string{"clivern"}}, ParseCommand("@zieeio assign @clivern"))
	assert.Equal(t, Command{Verb: "unassign", Args: []string{"clivern"}}, ParseCommand("@zieeio unassign clivern"))
	assert.Equal(t, Command{Verb: "close"}, ParseCommand("@zieeio close"))
	assert.Equal(t, Command{Verb: "reopen"}, ParseCommand("@zieeio reopen"))
	assert.Equal(t, Command{Verb: "spam"}, ParseCommand("@zieeio spam"))
	assert.Equal(t, Command{Verb: "summarize"}, ParseCommand("@zieeio summarize"))
	assert.Equal(t, Command{Verb: "queue", Args: []string{"hotfix"}}, ParseCommand("@zieeio queue hotfix"))
	assert.Equal(t, Command{Verb: "label", Args: []string{"bug"}}, ParseCommand("  @ZIEEIO Label bug  "))
	assert.Equal(t, Command{Verb: "close"}, ParseCommand("@zieeio close\nthis is a duplicate"))
	assert.Equal(t, Command{}, ParseCommand("please @zieeio close"))
	assert.Equal(t, Command{}, ParseCommand("@zieeio"))
	assert.Equal(t, Command{}, ParseCommand("@ziee close"))
	assert.Equal(t, Command{}, ParseCommand("thanks"))
}
