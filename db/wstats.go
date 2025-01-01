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
	APICallsMonth   int64
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

	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)

	err := r.db.QueryRow(
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

	err = r.db.QueryRow(
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
