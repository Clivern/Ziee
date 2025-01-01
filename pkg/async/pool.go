// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package async

import (
	"time"

	"github.com/panjf2000/ants/v2"
)

// ErrPoolSaturated is returned when the worker pool cannot accept more tasks.
var ErrPoolSaturated = ants.ErrPoolOverload

// Option configures a pool.
type Option = ants.Option

// Task is a unit of async work.
type Task func()

// Pool wraps the ants goroutine pool.
type Pool struct {
	pool *ants.Pool
}

// NewPool creates a goroutine pool with the given capacity.
func NewPool(size int, options ...Option) (*Pool, error) {
	options = append(options, ants.WithNonblocking(true))
	pool, err := ants.NewPool(size, options...)
	if err != nil {
		return nil, err
	}

	return &Pool{pool: pool}, nil
}

// Submit schedules a task for execution.
func (p *Pool) Submit(task Task) error {
	return p.pool.Submit(task)
}

// Tune changes the pool capacity at runtime.
func (p *Pool) Tune(size int) {
	p.pool.Tune(size)
}

// Cap returns the current pool capacity.
func (p *Pool) Cap() int {
	return p.pool.Cap()
}

// Running returns the number of running workers.
func (p *Pool) Running() int {
	return p.pool.Running()
}

// Free returns the number of available workers.
func (p *Pool) Free() int {
	return p.pool.Free()
}

// Reboot reactivates a released pool.
func (p *Pool) Reboot() {
	p.pool.Reboot()
}

// Release closes the pool and discards idle workers.
func (p *Pool) Release() {
	p.pool.Release()
}

// ReleaseTimeout closes the pool and waits for running tasks.
func (p *Pool) ReleaseTimeout(timeout time.Duration) error {
	return p.pool.ReleaseTimeout(timeout)
}
