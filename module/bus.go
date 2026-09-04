// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"fmt"

	"github.com/clivern/ziee/pkg/nats"

	"github.com/rs/zerolog/log"
)

var bus *nats.Client

// StartBus connects the shared NATS client used to publish work.
func StartBus() error {
	client, err := nats.New()
	if err != nil {
		return fmt.Errorf("connect nats bus: %w", err)
	}

	bus = client

	log.Info().
		Str("url", client.Config().URL).
		Str("name", client.Config().Name).
		Msg("NATS bus connected")

	return nil
}

// GetBus returns the shared NATS publisher.
func GetBus() *nats.Client {
	return bus
}

// StopBus drains and closes the shared NATS client.
func StopBus() {
	err := bus.Conn().Drain()
	if err != nil {
		log.Error().Err(err).Msg("Error draining NATS bus")
		bus.Close()
		return
	}

	bus = nil
}
