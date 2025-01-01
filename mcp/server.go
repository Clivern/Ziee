// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package mcp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Options configures the MCP server.
type Options struct {
	Name    string
	Version string
}

// Run starts the Ziee MCP server over streamable HTTP.
func Run(ctx context.Context, opts Options) error {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    opts.Name,
		Version: opts.Version,
	}, nil)

	RegisterMemoryTools(server)

	path := viper.GetString("app.mcp.path")
	port := viper.GetInt("app.mcp.port")
	timeout := time.Duration(viper.GetInt("app.timeout")) * time.Second

	handler := mcpsdk.NewStreamableHTTPHandler(func(*http.Request) *mcpsdk.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		MaxHeaderBytes:    230 * 1024,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadTimeout:       timeout + 5*time.Second,
		WriteTimeout:      timeout + 5*time.Second,
	}

	serr := make(chan error, 1)

	go func() {
		log.Info().
			Int("port", port).
			Str("path", path).
			Str("name", opts.Name).
			Str("version", opts.Version).
			Msg("Starting MCP HTTP server")

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serr:
		return fmt.Errorf("mcp server error: %w", err)
	case sig := <-quit:
		log.Info().
			Str("signal", sig.String()).
			Msg("Received MCP shutdown signal")

		shutdownTimeout := 30 * time.Second
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		log.Info().
			Dur("timeout", shutdownTimeout).
			Msg("Gracefully shutting down MCP server")

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			return fmt.Errorf("mcp server forced to shutdown: %w", err)
		}

		log.Info().Msg("MCP server shutdown complete")
		return nil
	}
}
