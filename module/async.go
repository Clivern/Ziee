// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/clivern/actx0/conf"
	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/async"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	// Async task statuses.
	AsyncTaskStatusPending   = "pending"
	AsyncTaskStatusRunning   = "running"
	AsyncTaskStatusCompleted = "completed"
	AsyncTaskStatusFailed    = "failed"

	// Async task priorities. Higher values are claimed first.
	AsyncPriorityHigh   = 100
	AsyncPriorityNormal = 50
	AsyncPriorityLow    = 10
)

var (
	ErrAsyncTaskNotFound       = errors.New("async task not found")
	ErrAsyncTaskHandlerMissing = errors.New("async task handler missing")
	ErrAsyncPoolSaturated      = errors.New("async worker pool saturated")
)

var asmgr *Async

// AsyncTaskHandler runs an async task and optionally returns a result.
type AsyncTaskHandler func(context.Context, *db.AsyncTask) (string, error)

// CreateAsyncTaskOptions is what you pass when creating an async task.
type CreateAsyncTaskOptions struct {
	WorkspaceId db.Id
	Type        string
	Payload     map[string]string
	Priority    int
}

// Async manages persisted async tasks and worker-pool execution.
type Async struct {
	TaskRepository db.AsyncTaskRepository
	Pool           *async.Pool
	wake           chan struct{}

	mu       sync.RWMutex
	handlers map[string]AsyncTaskHandler
}

// Start creates the configured async worker pool.
func Start(atr db.AsyncTaskRepository) (*Async, error) {
	pool, err := async.NewPool(viper.GetInt("app.worker_pool"))
	if err != nil {
		return nil, fmt.Errorf("create async worker pool: %w", err)
	}

	mgr := &Async{
		TaskRepository: atr,
		Pool:           pool,
		wake:           make(chan struct{}, 1),
		handlers:       map[string]AsyncTaskHandler{},
	}
	asmgr = mgr

	go func() {
		ticker := time.NewTicker(conf.AsyncClaimInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
			case <-mgr.wake:
			}
			if err := mgr.RunTasks(); err != nil {
				log.Error().Err(err).Msg("Failed to drain pending async tasks")
			}
		}
	}()

	log.Info().
		Int("workers", pool.Cap()).
		Msg("Async worker pool started")

	return mgr, nil
}

// GetAsyncMgr returns the app-level async manager.
func GetAsyncMgr() *Async {
	return asmgr
}

// Stop releases the worker pool after waiting for running work.
func (a *Async) Stop(timeout time.Duration) error {
	return a.Pool.ReleaseTimeout(timeout)
}

// RegisterHandler registers a task handler by async task type.
func (a *Async) RegisterHandler(taskType string, handler AsyncTaskHandler) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.handlers[taskType] = handler
}

// CreateTask persists a new pending async task.
func (a *Async) CreateTask(options *CreateAsyncTaskOptions) (*db.AsyncTask, error) {
	raw, err := json.Marshal(options.Payload)
	if err != nil {
		return nil, err
	}
	payload := string(raw)

	task := &db.AsyncTask{
		WorkspaceId: options.WorkspaceId,
		Type:        options.Type,
		Status:      AsyncTaskStatusPending,
		Payload:     &payload,
		Attempts:    0,
		Priority:    options.Priority,
	}

	err = a.TaskRepository.Create(task)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("taskId", task.Id.String()).
		Str("taskType", task.Type).
		Int("priority", task.Priority).
		Str("workspaceId", task.WorkspaceId.String()).
		Msg("Async task enqueued")

	a.Poll()

	return task, nil
}

// Poll wakes the async poller to drain pending tasks.
func (a *Async) Poll() {
	select {
	case a.wake <- struct{}{}:
	default:
	}
}

// RunTasks claims and submits one batch of due pending tasks.
func (a *Async) RunTasks() error {
	now := time.Now().UTC()
	limit := viper.GetInt("app.worker_pool")
	if free := a.Pool.Free(); free < limit {
		limit = free
	}
	if limit <= 0 {
		return nil
	}

	tasks, err := a.TaskRepository.ClaimDue(now, limit)
	if err != nil {
		return err
	}
	if len(tasks) > 0 {
		a.CleanupCompletedTasks()
	}
	for i, task := range tasks {
		if err := a.SubmitTask(task.Id); err != nil {
			for _, remaining := range tasks[i:] {
				_ = a.TaskRepository.ReleaseClaim(remaining.Id)
			}
			return err
		}
	}
	return nil
}

// SubmitTask submits a task to the worker pool.
func (a *Async) SubmitTask(taskId db.Id) error {
	err := a.Pool.Submit(func() {
		_ = a.executeTask(context.Background(), taskId)
	})
	if errors.Is(err, async.ErrPoolSaturated) {
		return ErrAsyncPoolSaturated
	}
	return err
}

// executeTask runs a task handler with retries.
func (a *Async) executeTask(ctx context.Context, taskId db.Id) error {
	task, err := a.TaskRepository.GetById(taskId)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrAsyncTaskNotFound
	}

	switch task.Status {
	case AsyncTaskStatusCompleted, AsyncTaskStatusFailed, AsyncTaskStatusPending:
		return nil
	}

	handler, ok := a.handler(task.Type)
	if !ok {
		err := fmt.Errorf("%w: %s", ErrAsyncTaskHandlerMissing, task.Type)
		return a.failTask(task, err)
	}

	log.Info().
		Str("taskId", task.Id.String()).
		Str("taskType", task.Type).
		Int("priority", task.Priority).
		Str("workspaceId", task.WorkspaceId.String()).
		Msg("Async task started")

	var lastErr error
	started := time.Now().UTC()
	for retry := 0; retry <= conf.AsyncTaskMaxRetries; retry++ {
		if retry > 0 {
			delay := time.Second * time.Duration(1<<(retry-1))
			log.Warn().
				Err(lastErr).
				Str("taskId", task.Id.String()).
				Int("retry", retry).
				Dur("delay", delay).
				Msg("Async task failed, retrying")

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		now := time.Now().UTC()
		task.Status = AsyncTaskStatusRunning
		task.Attempts++
		task.LockedAt = &now
		task.Error = nil
		err = a.TaskRepository.Update(task)
		if err != nil {
			return err
		}

		result, err := handler(ctx, task)
		if err == nil {
			now := time.Now().UTC()
			task.Status = AsyncTaskStatusCompleted
			task.Result = stringPtr(result)
			task.Error = nil
			task.CompletedAt = &now

			err = a.TaskRepository.Update(task)
			if err != nil {
				return err
			}
			log.Info().
				Str("taskId", task.Id.String()).
				Str("taskType", task.Type).
				Int("priority", task.Priority).
				Str("workspaceId", task.WorkspaceId.String()).
				Int("attempt", task.Attempts).
				Dur("duration", time.Since(started)).
				Msg("Async task succeeded")
			return nil
		}
		lastErr = err
	}

	return a.failTask(task, lastErr)
}

// handler returns the registered handler for a task type.
func (a *Async) handler(taskType string) (AsyncTaskHandler, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	handler, ok := a.handlers[taskType]
	return handler, ok
}

// CleanupCompletedTasks deletes completed tasks past the retention window.
func (a *Async) CleanupCompletedTasks() {
	cutoff := time.Now().UTC().Add(conf.AsyncCompletedTasksRetention)
	count, err := a.TaskRepository.DeleteCompletedOlderThan(cutoff)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to delete old completed async tasks")
		return
	}
	if count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Deleted old completed async tasks")
	}
}

// failTask marks a task as failed and persists the error.
func (a *Async) failTask(task *db.AsyncTask, taskErr error) error {
	now := time.Now().UTC()
	task.Status = AsyncTaskStatusFailed
	task.LockedAt = nil
	task.CompletedAt = &now

	raw, err := json.Marshal(map[string]string{
		"error": string(taskErr.Error()),
	})
	if err != nil {
		return err
	}
	task.Error = stringPtr(string(raw))

	err = a.TaskRepository.Update(task)
	if err != nil {
		return err
	}

	log.Error().
		Str("taskId", task.Id.String()).
		Str("taskType", task.Type).
		Int("priority", task.Priority).
		Str("workspaceId", task.WorkspaceId.String()).
		Int("attempt", task.Attempts).
		Err(taskErr).
		Msg("Async task failed")

	return taskErr
}

// stringPtr returns a pointer to the given string.
func stringPtr(value string) *string {
	return &value
}
