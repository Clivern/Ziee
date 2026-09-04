// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"

	"github.com/clivern/ziee/pkg/nats"

	"github.com/rs/zerolog/log"
)

// Handler processes an inbound NATS message.
type Handler func(context.Context, *nats.Msg) error

// Registration binds a subject to a handler.
type Registration struct {
	Subject string
	Handler Handler
}

// Registrations is a list of all registered handlers.
var Registrations []Registration

// On registers a queue worker handler for subject.
func On(subject string, handler Handler) {
	Registrations = append(Registrations, Registration{
		Subject: subject,
		Handler: handler,
	})
}

// Bind attaches all registered handlers to the NATS client using the queue group.
func Bind(client *nats.Client, queue string) error {
	for _, reg := range Registrations {
		_, err := client.QueueSubscribe(reg.Subject, queue, func(msg *nats.Msg) {
			err := reg.Handler(context.Background(), msg)
			if err != nil {
				log.Error().
					Err(err).
					Str("subject", msg.Subject).
					Msg("Worker handler failed")
			}
		})
		if err != nil {
			return err
		}

		log.Info().
			Str("subject", reg.Subject).
			Str("queue", queue).
			Msg("Worker subscribed")
	}

	return nil
}
