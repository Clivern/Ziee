// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"fmt"

	"github.com/clivern/ziee/conf"
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/broker"

	"github.com/rs/zerolog/log"
)

var bus *broker.Client

// StartBus connects the shared NATS client used to publish work.
func StartBus() error {
	client, err := broker.New()
	if err != nil {
		return fmt.Errorf("connect nats bus: %w", err)
	}

	bus = client

	log.Info().
		Str("url", client.Config().NATS.URL).
		Str("name", client.Config().NATS.Name).
		Msg("NATS bus connected")

	return nil
}

// GetBus returns the shared NATS publisher.
func GetBus() *broker.Client {
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

// EnqueueTask records a pending async task and publishes it on NATS.
func EnqueueTask(taskType, subject string, payload map[string]string, workspaceId db.Id) error {
	taskId, err := db.NewId()
	if err != nil {
		return err
	}

	payload["taskId"] = taskId.String()
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body := string(raw)

	err = db.NewAsyncTaskRepository(db.GetDB()).Create(&db.AsyncTask{
		Id:          taskId,
		WorkspaceId: workspaceId,
		Type:        taskType,
		Status:      db.AsyncTaskStatusPending,
		Payload:     &body,
	})
	if err != nil {
		return err
	}

	err = GetBus().Publish(subject, raw)
	if err != nil {
		return err
	}

	return GetBus().Flush()
}

// RepublishPendingTasks publishes pending async tasks to NATS again.
func RepublishPendingTasks() (int, error) {
	tasks, err := db.NewAsyncTaskRepository(db.GetDB()).ListByStatus(db.AsyncTaskStatusPending)
	if err != nil {
		return 0, err
	}

	for n, task := range tasks {
		var subject string
		switch task.Type {
		case db.AsyncTaskTypeDocIndex:
			subject = conf.NATSSubjectDocIndex
		case db.AsyncTaskTypeDocDelete:
			subject = conf.NATSSubjectDocDelete
		case db.AsyncTaskTypeRepoBootstrap:
			subject = conf.NATSSubjectRepoBootstrap
		case db.AsyncTaskTypeGitHubIssue:
			subject = conf.NATSSubjectGitHubIssue
		}

		err = GetBus().Publish(subject, []byte(*task.Payload))
		if err != nil {
			return n, err
		}

		log.Info().
			Str("id", task.Id.String()).
			Str("type", task.Type).
			Str("subject", subject).
			Msg("Pending task republished")
	}

	err = GetBus().Flush()
	if err != nil {
		return 0, err
	}

	return len(tasks), nil
}
