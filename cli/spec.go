// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"

	"github.com/clivern/ziee/policy/spec"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var specCmd = &cobra.Command{
	Use:   "spec [path]",
	Short: "Parse a .ziee.yml file and print it",
	Args:  cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		path := ".ziee.yml"
		if len(args) == 1 {
			path = args[0]
		}

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file, err := spec.Parse(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		out, err := yaml.Marshal(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print(string(out))
	},
}

func init() {
	rootCmd.AddCommand(specCmd)
}
