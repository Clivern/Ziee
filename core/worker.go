// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/ai"
	"github.com/clivern/ziee/pkg/github"
	"github.com/clivern/ziee/pkg/nats"
	"github.com/clivern/ziee/pkg/qdrant"
	"github.com/clivern/ziee/pkg/storage"
	"github.com/clivern/ziee/service/knowledge"
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

	err = github.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize github app: %w", err)
	}

	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to initialize document storage: %w", err)
	}

	vdb, err := qdrant.New()
	if err != nil {
		return fmt.Errorf("failed to initialize qdrant: %w", err)
	}

	defer func() {
		err := vdb.Close()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error closing qdrant client")
		}
	}()

	ksvc := knowledge.New(knowledge.Dependencies{
		Documents: db.NewWorkspaceDocumentRepository(
			db.GetDB(true),
		),
		Embed:         ai.NewEmbedClient(),
		Vectors:       vdb,
		Store:         store,
		Usage:         db.NewUsageRepository(db.GetDB(false)),
		Subscriptions: db.NewSubscriptionRepository(db.GetDB(false)),
	})

	client, err := nats.New()
	if err != nil {
		return fmt.Errorf("failed to connect to nats: %w", err)
	}

	defer client.Close()

	cfg := client.Config()

	worker.Register(worker.Dependencies{
		Knowledge: ksvc,
	})

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
