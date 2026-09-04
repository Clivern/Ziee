// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"

	"github.com/actx0/ziee/pkg/nats"

	"github.com/rs/zerolog/log"
)

// Handler processes an inbound NATS message.
type Handler func(context.Context, *nats.Msg) error

type registration struct {
	Subject string
	Handler Handler
}

var registrations []registration

// On registers a queue worker handler for subject.
func On(subject string, handler Handler) {
	registrations = append(registrations, registration{
		Subject: subject,
		Handler: handler,
	})
}

// Bind attaches all registered handlers to the NATS client using the queue group.
func Bind(client *nats.Client, queue string) error {
	for _, reg := range registrations {
		reg := reg

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

// Register attaches all worker handlers.
func Register() {
}
