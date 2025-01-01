// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var config string

var rootCmd = &cobra.Command{
	Use: "ziee",
	Short: `🐀 The Autonomous Merge Layer for Agent-Scale Delivery


If you have any suggestions, bug reports, or annoyances please report
them to our issue tracker at <https://github.com/actx0/ziee/issues>`,
}

// Execute runs cmd tool
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
