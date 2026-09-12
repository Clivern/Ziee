// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package main

import (
	"embed"
	"io/fs"

	"github.com/clivern/ziee/cli"
	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/locale"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
	edition = "oss"
)

//go:embed web/dist/*
var static embed.FS

//go:embed locale/locales
var localeFS embed.FS

// main is the application entry point.
func main() {
	cli.Version = version
	cli.Commit = commit
	cli.Date = date
	cli.BuiltBy = builtBy
	cli.Edition = edition
	cli.Static = static
	conf.BuiltEdition = edition

	// Load locales
	if sub, err := fs.Sub(localeFS, "locale/locales"); err == nil {
		err = locale.Load(sub)
		if err != nil {
			panic(err)
		}
	}

	cli.Execute()
}
