// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"
	"fmt"
	"time"

	"github.com/clivern/ziee/db"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

var (
	ErrPQueueNotFound     = errors.New("merge queue entry not found")
	ErrFailedCreatePQueue = errors.New("failed create merge queue entry")
	ErrFailedUpdatePQueue = errors.New("failed update merge queue entry")
	ErrFailedListPQueue   = errors.New("failed list merge queue")
	ErrFailedGetPQueue    = errors.New("failed get merge queue entry")
)

// PQueue is the module for repository merge-queue entries.
type PQueue struct {
	PQueueRepository db.PQueueRepository
}

// AppendPRRequest is what you pass when appending a PR to a repo queue.
type AppendPRRequest struct {
	RemoteId int64  `json:"remoteId"`
	Priority string `json:"priority"`
	Checksum string `json:"checksum"`
	Meta     string `json:"meta"`
}

// UpdatePRRequest patches selected fields on a merge-queue entry.
type UpdatePRRequest struct {
	Status   *string    `json:"status"`
	Checksum *string    `json:"checksum"`
	MergedAt *time.Time `json:"mergedAt"`
}

// PR is one entry in a working-queue snapshot, including PRs ahead of it.
type PR struct {
	Id       db.Id
	RemoteId int64
	Checksum string
	Status   string
	Front    []PR
}

// NewPQueue creates a merge-queue module with the given repository.
func NewPQueue(items db.PQueueRepository) *PQueue {
	return &PQueue{PQueueRepository: items}
}

// Append adds a pull request to the end of a repository queue.
func (p *PQueue) Append(repoId db.Id, req *AppendPRRequest) (*db.PQueue, error) {
	items, err := p.PQueueRepository.ListUnmergedByRepoId(repoId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPQueue, err)
	}

	rank := 1
	for _, existing := range items {
		if existing.Status == db.PQueueStatusQueued || existing.Status == db.PQueueStatusChecking {
			if existing.Rank >= rank {
				rank = existing.Rank + 1
			}
		}
	}

	priority := lo.Ternary(lo.IsNotEmpty(req.Priority), req.Priority, db.PQueuePriorityMedium)
	meta := lo.Ternary(lo.IsNotEmpty(req.Meta), req.Meta, "{}")
	now := time.Now().UTC()

	item := &db.PQueue{
		RepoId:   repoId,
		RemoteId: req.RemoteId,
		Priority: priority,
		Rank:     rank,
		Status:   db.PQueueStatusQueued,
		Checksum: req.Checksum,
		Meta:     meta,
		OpenedAt: &now,
	}

	err = p.PQueueRepository.Create(item)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePQueue, err)
	}

	log.Info().
		Str("repoId", repoId.String()).
		Int64("remoteId", req.RemoteId).
		Int("rank", rank).
		Msg("Pull request appended to merge queue")

	return item, nil
}

// Update patches status, checksum, and/or merged_at for a queued PR.
func (p *PQueue) Update(repoId db.Id, remoteId int64, req *UpdatePRRequest) (*db.PQueue, error) {
	item, err := p.PQueueRepository.GetByRepoIdAndRemoteId(repoId, remoteId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPQueue, err)
	}
	if item == nil {
		return nil, ErrPQueueNotFound
	}

	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.Checksum != nil {
		item.Checksum = *req.Checksum
	}
	if req.MergedAt != nil {
		item.MergedAt = req.MergedAt
	}

	err = p.PQueueRepository.Update(item)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdatePQueue, err)
	}

	log.Info().
		Str("repoId", repoId.String()).
		Int64("remoteId", remoteId).
		Str("status", item.Status).
		Msg("Merge queue entry updated")

	return item, nil
}

// ListQueued returns the top limit queued PRs for a repository.
func (p *PQueue) ListQueued(repoId db.Id, limit int) ([]*db.PQueue, error) {
	items, err := p.PQueueRepository.ListQueuedByRepoId(repoId, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPQueue, err)
	}

	return items, nil
}

// Get returns one merge-queue entry so callers can compare checksums.
func (p *PQueue) Get(repoId db.Id, remoteId int64) (*db.PQueue, error) {
	item, err := p.PQueueRepository.GetByRepoIdAndRemoteId(repoId, remoteId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPQueue, err)
	}
	if item == nil {
		return nil, ErrPQueueNotFound
	}

	return item, nil
}

// IsStillValid reports whether the working-batch snapshot is still current.
func (p *PQueue) IsStillValid(repoId db.Id, snapshot []PR) (bool, error) {
	// Check if the PRs in the snapshot still have the same checksum
	// in DB and have not moved to failed or dequeued.
	for _, pr := range snapshot {
		item, err := p.PQueueRepository.GetById(pr.Id)
		if err != nil {
			return false, fmt.Errorf("%w: %v", ErrFailedGetPQueue, err)
		}

		if item.Checksum != pr.Checksum ||
			item.Status == db.PQueueStatusFailed ||
			item.Status == db.PQueueStatusDequeued {
			return false, nil
		}
	}

	// Check if the PRs in the snapshot still have the same front in DB.
	live, err := p.PQueueRepository.ListQueuedByRepoId(repoId, 1000)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrFailedListPQueue, err)
	}

	fromDB := SnapshotWorkingQueue(live)
	byId := make(map[db.Id]PR, len(fromDB))
	for _, pr := range fromDB {
		byId[pr.Id] = pr
	}

	for _, pr := range snapshot {
		live, ok := byId[pr.Id]
		if !ok {
			continue
		}

		if HasNewInFront(pr.Front, live.Front) {
			return false, nil
		}
	}

	return true, nil
}

// HasNewInFront reports whether live Front gained a PR that was not in saved Front.
func HasNewInFront(saved, live []PR) bool {
	known := make(map[db.Id]struct{}, len(saved))
	for _, pr := range saved {
		known[pr.Id] = struct{}{}
	}

	for _, pr := range live {
		if _, ok := known[pr.Id]; !ok {
			return true
		}
	}

	return false
}

// SnapshotWorkingQueue builds the saved slice from a queue just loaded from DB.
func SnapshotWorkingQueue(snapshot []*db.PQueue) []PR {
	saved := make([]PR, 0, len(snapshot))
	for i, pr := range snapshot {
		front := make([]PR, 0, i)
		for j := 0; j < i; j++ {
			front = append(front, PR{
				Id:       snapshot[j].Id,
				RemoteId: snapshot[j].RemoteId,
				Checksum: snapshot[j].Checksum,
				Status:   snapshot[j].Status,
			})
		}
		saved = append(saved, PR{
			Id:       pr.Id,
			RemoteId: pr.RemoteId,
			Checksum: pr.Checksum,
			Status:   pr.Status,
			Front:    front,
		})
	}

	return saved
}
