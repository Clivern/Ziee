// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package worker

import (
	"context"
	"errors"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/broker"

	"github.com/rs/zerolog/log"
)

var ErrInvalidPayload = errors.New("invalid worker payload")

// Instance is the process-wide worker wired at startup.
var Instance *Worker

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

// Worker holds shared modules and repositories for queue handlers.
type Worker struct {
	Knowledge     Knowledge
	Tasks         db.AsyncTaskRepository
	Repository    *module.Repository
	PQueue        *module.PQueue
	Usage         db.UsageRepository
	Subscriptions db.SubscriptionRepository
}

type handlers struct {
	*Worker
}

// New wires repositories and modules once for the worker process.
func New(knowledge Knowledge) *Worker {
	conn := db.GetDB()

	repos := db.NewRepositoriesRepository(conn)
	repoMeta := db.NewRepositoryMetaRepository(conn)
	spam := db.NewRepositorySpamUserRepository(conn)

	w := &Worker{
		Knowledge:     knowledge,
		Tasks:         db.NewAsyncTaskRepository(conn),
		Repository:    module.NewRepository(repos, repoMeta, spam),
		PQueue:        module.NewPQueue(db.NewPQueueRepository(conn)),
		Usage:         db.NewUsageRepository(conn),
		Subscriptions: db.NewSubscriptionRepository(conn),
	}

	Instance = w
	w.register()

	return w
}

func (w *Worker) register() {
	h := &handlers{Worker: w}

	On(db.AsyncTaskTypeDocIndex, h.HandleDocumentIndex)
	On(db.AsyncTaskTypeDocDelete, h.HandleDocumentDelete)
	On(db.AsyncTaskTypeRepoBootstrap, h.HandleRepositoryBootstrap)
	On(db.AsyncTaskTypeGitHubIssue, h.HandleGitHubIssue)
	On(db.AsyncTaskTypeGitHubIssueComment, h.HandleGitHubIssueComment)
	On(db.AsyncTaskTypeGitHubPullRequestComment, h.HandleGitHubPullRequestComment)
	On(db.AsyncTaskTypeGitHubPullRequest, h.HandleGitHubPullRequest)
	On(db.AsyncTaskTypeRepoLabels, h.HandleRepositoryLabels)
	On(db.AsyncTaskTypeRepoMerge, h.HandleRepositoryMerge)
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
