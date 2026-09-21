// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package eval

import (
	"strings"

	"github.com/clivern/ziee/pkg/github/app"
	"github.com/samber/lo"
)

// Command is one `@<bot_name>` verb and its arguments from a comment.
type Command struct {
	Verb string
	Args []string
}

// ParseCommand reads the first line of a comment as `@<bot_name> <verb> [args...]`.
// bot_name comes from app.oauth.github.bot_name (default zieeio).
func ParseCommand(comment string) Command {
	line, _, _ := strings.Cut(comment, "\n")
	fields := strings.Fields(line)
	mention := "@" + app.GetConfig().BotName
	if len(fields) < 2 || !strings.EqualFold(fields[0], mention) {
		return Command{}
	}

	cmd := Command{Verb: strings.ToLower(fields[1])}
	if len(fields) > 2 {
		cmd.Args = lo.Map(fields[2:], func(arg string, _ int) string {
			return strings.TrimPrefix(arg, "@")
		})
	}

	return cmd
}
