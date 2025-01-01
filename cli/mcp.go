// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/actx0/ziee/core"
	"github.com/actx0/ziee/mcp"

	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the Ziee MCP server",
	Long:  `Start the Ziee Model Context Protocol server over streamable HTTP.`,
	Run: func(_ *cobra.Command, _ []string) {
		err := core.Load(config)
		if err != nil {
			panic(err.Error())
		}

		err = core.SetupLogging()
		if err != nil {
			panic(err.Error())
		}

		err = mcp.Run(context.Background(), mcp.Options{
			Name:    "ziee",
			Version: Version,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %s\n", err.Error())
			os.Exit(1)
		}
	},
}

// init registers the mcp subcommand.
func init() {
	mcpCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	mcpCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(mcpCmd)
}
