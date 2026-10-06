// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/broker"
	"github.com/clivern/ziee/pkg/github/app"
	"github.com/clivern/ziee/sandbox"
	"github.com/clivern/ziee/worker"

	"github.com/rs/zerolog/log"
)

// RunWorker starts the NATS worker and blocks until shutdown.
func RunWorker() error {
	err := db.InitDB(ReadWriteDatabase(), ReadOnlyDatabase()...)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	defer func() {
		err := db.CloseDB()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error closing database connection")
		}
	}()

	err = app.Init(module.NewCache(db.NewKVRepository(db.GetDB())))
	if err != nil {
		return fmt.Errorf("failed to initialize github app: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	check := sandbox.IsDockerRunning(ctx)
	cancel()

	if !check {
		log.Warn().Msg("Docker is not running, sandbox can't run!")
	}

	w := worker.New()

	defer func() {
		err := w.Close()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error closing worker resources")
		}
	}()

	client, err := broker.New()
	if err != nil {
		return fmt.Errorf("failed to connect to nats: %w", err)
	}

	defer client.Close()

	nats := client.Config().NATS

	err = worker.Bind(client, nats.Queue)
	if err != nil {
		return fmt.Errorf("failed to bind workers: %w", err)
	}

	log.Info().
		Str("url", nats.URL).
		Str("name", nats.Name).
		Str("queue", nats.Queue).
		Msg("Starting NATS worker")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit

	log.Info().
		Str("signal", sig.String()).
		Msg("Received shutdown signal")

	err = client.Conn().Drain()
	if err != nil {
		return fmt.Errorf("failed to drain nats connection: %w", err)
	}

	deadline := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if client.Conn().IsClosed() {
			break
		}

		select {
		case <-deadline:
			return fmt.Errorf("nats drain timed out")
		case <-ticker.C:
		}
	}

	log.Info().Msg("Worker shutdown complete")

	return nil
}
