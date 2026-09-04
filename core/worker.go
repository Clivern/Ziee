// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/actx0/ziee/pkg/nats"
	"github.com/actx0/ziee/worker"

	"github.com/rs/zerolog/log"
)

// RunWorker starts the NATS worker and blocks until shutdown.
func RunWorker() error {
	client, err := nats.New()
	if err != nil {
		return fmt.Errorf("failed to connect to nats: %w", err)
	}

	defer client.Close()

	cfg := client.Config()

	worker.Register()

	err = worker.Bind(client, cfg.Queue)
	if err != nil {
		return fmt.Errorf("failed to bind workers: %w", err)
	}

	log.Info().
		Str("url", cfg.URL).
		Str("name", cfg.Name).
		Str("queue", cfg.Queue).
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
