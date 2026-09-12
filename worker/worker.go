// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"errors"

	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/broker"

	"github.com/rs/zerolog/log"
)

var ErrInvalidPayload = errors.New("invalid worker payload")

// Handler processes an inbound NATS message.
type Handler func(context.Context, *broker.Msg) error

// Registration binds a subject to a handler.
type Registration struct {
	Subject string
	Handler Handler
}

// Registrations is a list of all registered handlers.
var Registrations []Registration

// Knowledge indexes and deletes workspace documents.
type Knowledge interface {
	Index(ctx context.Context, documentId db.Id) error
	Delete(ctx context.Context, documentId db.Id, internalId string) error
}

// Dependencies are the services required by worker handlers.
type Dependencies struct {
	Knowledge Knowledge
	Tasks     db.AsyncTaskRepository
}

type handlers struct {
	knowledge Knowledge
	tasks     db.AsyncTaskRepository
}

// Register attaches all worker handlers.
func Register(deps Dependencies) {
	h := &handlers{knowledge: deps.Knowledge, tasks: deps.Tasks}

	On(conf.NATSSubjectDocIndex, h.HandleDocumentIndex)
	On(conf.NATSSubjectDocDelete, h.HandleDocumentDelete)
	On(conf.NATSSubjectRepoBootstrap, h.HandleRepositoryBootstrap)
	On(conf.NATSSubjectGitHubIssue, h.HandleGitHubIssue)
	On(conf.NATSSubjectRepoLabels, h.HandleRepositoryLabels)
}

// On registers a queue worker handler for subject.
func On(subject string, handler Handler) {
	Registrations = append(Registrations, Registration{
		Subject: subject,
		Handler: handler,
	})
}

// Bind attaches all registered handlers to the NATS client using the queue group.
func Bind(client *broker.Client, queue string) error {
	for _, reg := range Registrations {
		_, err := client.QueueSubscribe(reg.Subject, queue, func(msg *broker.Msg) {
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
