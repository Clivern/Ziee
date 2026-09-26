// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitParseCommand(t *testing.T) {
	viper.Set("app.oauth.github.bot_name", "zieeai")

	assert.Equal(t, Command{Verb: "label", Args: []string{"bug", "security"}}, ParseCommand("@zieeai label bug security"))
	assert.Equal(t, Command{Verb: "unlabel", Args: []string{"bot"}}, ParseCommand("@zieeai unlabel bot"))
	assert.Equal(t, Command{Verb: "assign", Args: []string{"clivern"}}, ParseCommand("@zieeai assign clivern"))
	assert.Equal(t, Command{Verb: "assign", Args: []string{"clivern"}}, ParseCommand("@zieeai assign @clivern"))
	assert.Equal(t, Command{Verb: "unassign", Args: []string{"clivern"}}, ParseCommand("@zieeai unassign clivern"))
	assert.Equal(t, Command{Verb: "close", Args: []string{}}, ParseCommand("@zieeai close"))
	assert.Equal(t, Command{Verb: "reopen", Args: []string{}}, ParseCommand("@zieeai reopen"))
	assert.Equal(t, Command{Verb: "spam", Args: []string{}}, ParseCommand("@zieeai spam"))
	assert.Equal(t, Command{Verb: "summarize", Args: []string{}}, ParseCommand("@zieeai summarize"))
	assert.Equal(t, Command{Verb: "queue", Args: []string{"hotfix"}}, ParseCommand("@zieeai queue hotfix"))
	assert.Equal(t, Command{Verb: "label", Args: []string{"bug"}}, ParseCommand("  @ZIEEAI Label bug  "))
	assert.Equal(t, Command{Verb: "close", Args: []string{}}, ParseCommand("@zieeai close\nthis is a duplicate"))
	assert.Equal(t, Command{}, ParseCommand("please @zieeai close"))
	assert.Equal(t, Command{}, ParseCommand("@zieeai"))
	assert.Equal(t, Command{}, ParseCommand("@ziee close"))
	assert.Equal(t, Command{}, ParseCommand("thanks"))
}

func TestUnitParseCommandCustomBotName(t *testing.T) {
	viper.Set("app.oauth.github.bot_name", "acmebot")
	t.Cleanup(func() { viper.Set("app.oauth.github.bot_name", "zieeai") })

	assert.Equal(t, Command{Verb: "queue", Args: []string{}}, ParseCommand("@acmebot queue"))
	assert.Equal(t, Command{Verb: "label", Args: []string{"bug"}}, ParseCommand("@AcmeBot label bug"))
	assert.Equal(t, Command{}, ParseCommand("@zieeai queue"))
}
