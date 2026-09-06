// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
)

const UsageTypeAITokens = "ai_tokens"
const UsageTypeAICost = "ai_cost"
const UsageUnitTokens = "tokens"
const UsageUnitNanoUSD = "nano_usd"

// WorkspaceStats holds aggregate metrics for a workspace dashboard.
type WorkspaceStats struct {
	DocumentsStored int64
}

// WorkspaceStatsRepository loads workspace dashboard metrics.
type WorkspaceStatsRepository interface {
	GetByWorkspaceId(workspaceId Id) (*WorkspaceStats, error)
}

type WorkspaceStatsRepositoryPostgres struct {
	db *sql.DB
}

// NewWorkspaceStatsRepository returns a workspace stats repository.
func NewWorkspaceStatsRepository(db *sql.DB) WorkspaceStatsRepository {
	return &WorkspaceStatsRepositoryPostgres{db: db}
}

// GetByWorkspaceId returns a workspace stats by workspace id.
func (r *WorkspaceStatsRepositoryPostgres) GetByWorkspaceId(workspaceId Id) (*WorkspaceStats, error) {
	stats := &WorkspaceStats{}

	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM workspace_documents
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&stats.DocumentsStored)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
