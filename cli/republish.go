// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/clivern/ziee/core"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var republishCmd = &cobra.Command{
	Use:   "republish",
	Short: "Republish pending async tasks to NATS",
	Run: func(_ *cobra.Command, _ []string) {
		err := core.Load(config)
		if err != nil {
			panic(err.Error())
		}

		err = core.SetupLogging()
		if err != nil {
			panic(err.Error())
		}

		err = db.InitDB(core.ReadWriteDatabase(), core.ReadOnlyDatabase()...)
		if err != nil {
			panic(fmt.Sprintf("Database error: %s", err.Error()))
		}

		defer db.CloseDB()

		err = module.StartBus()
		if err != nil {
			panic(fmt.Sprintf("NATS error: %s", err.Error()))
		}

		defer module.StopBus()

		count, err := module.RepublishPendingTasks()
		if err != nil {
			panic(fmt.Sprintf("Republish error: %s", err.Error()))
		}

		log.Info().
			Int("count", count).
			Msg("Pending tasks republished")
	},
}

func init() {
	republishCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"mgmt_config.prod.yml",
		"Absolute path to config file (required)",
	)
	republishCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(republishCmd)
}
