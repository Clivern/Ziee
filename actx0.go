// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

//go:build ignore

package main

import (
	"embed"
	"io/fs"

	"github.com/actx0/ziee/cli"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/pkg/resend"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

//go:embed web/dist/*
var static embed.FS

//go:embed locale/locales
var localeFS embed.FS

//go:embed mails/*.html
var mailsFS embed.FS

// main is the application entry point.
func main() {
	cli.Version = version
	cli.Commit = commit
	cli.Date = date
	cli.BuiltBy = builtBy
	cli.Static = static

	// Load locales
	if sub, err := fs.Sub(localeFS, "locale/locales"); err == nil {
		err = locale.Load(sub)
		if err != nil {
			panic(err)
		}
	}

	// Load mail templates
	if sub, err := fs.Sub(mailsFS, "mails"); err == nil {
		err = resend.LoadTemplates(sub)
		if err != nil {
			panic(err)
		}
	}

	cli.Execute()
}
