// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const UsageTypeAPICalls = "api_calls"
const UsageTypeAITokens = "ai_tokens"
const UsageTypeAICost = "ai_cost"
const UsageUnitCalls = "calls"
const UsageUnitTokens = "tokens"
const UsageUnitNanoUSD = "nano_usd"

// WorkspaceStats holds aggregate metrics for a workspace dashboard.
type WorkspaceStats struct {
	MemoriesStored int64
	APICallsMonth  int64
	ActiveAgents   int64
	Prompts        int64
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

	// Memories Stored
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM session_memories sm
		INNER JOIN sessions s ON s.id = sm.session_id
		WHERE s.workspace_id = $1`,
		workspaceId.String(),
	).Scan(&stats.MemoriesStored)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)

	// APICalls Month
	err = r.db.QueryRow(
		`SELECT COALESCE(SUM(quantity), 0)
		FROM usage
		WHERE workspace_id = $1
			AND type = $2
			AND period_start >= $3
			AND period_start < $4`,
		workspaceId.String(),
		UsageTypeAPICalls,
		monthStart,
		monthEnd,
	).Scan(&stats.APICallsMonth)
	if err != nil {
		return nil, err
	}

	// Active Agents
	err = r.db.QueryRow(
		`SELECT COUNT(*)
		FROM agents
		WHERE workspace_id = $1 AND status = 'active'`,
		workspaceId.String(),
	).Scan(&stats.ActiveAgents)
	if err != nil {
		return nil, err
	}

	// Prompts
	err = r.db.QueryRow(
		`SELECT COUNT(*)
		FROM prompts
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&stats.Prompts)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
