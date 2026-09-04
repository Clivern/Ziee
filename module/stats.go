// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"github.com/clivern/ziee/db"
)

// Stats is the module for workspace dashboard metrics.
type Stats struct {
	WorkspaceRepository db.WorkspaceRepository
	StatsRepository     db.WorkspaceStatsRepository
}

// NewStats creates a stats module with the given repositories.
func NewStats(workspaces db.WorkspaceRepository, stats db.WorkspaceStatsRepository) *Stats {
	return &Stats{
		WorkspaceRepository: workspaces,
		StatsRepository:     stats,
	}
}

// WorkspaceStatsResponse is workspace metrics shaped for API responses.
type WorkspaceStatsResponse struct {
	APICallsMonth   int64 `json:"apiCallsMonth"`
	DocumentsStored int64 `json:"documentsStored"`
}

// GetWorkspaceStats returns dashboard metrics for a workspace.
func (s *Stats) GetWorkspaceStats(workspaceId db.Id) (*WorkspaceStatsResponse, error) {
	workspace, err := s.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	stats, err := s.StatsRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}

	return &WorkspaceStatsResponse{
		APICallsMonth:   stats.APICallsMonth,
		DocumentsStored: stats.DocumentsStored,
	}, nil
}
