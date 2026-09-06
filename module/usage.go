// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/pkg/ai"
	"github.com/clivern/ziee/pkg/util"
)

const bytesPerGB = 1024 * 1024 * 1024

// Usage provides workspace consumption helpers.
type Usage struct{}

// WorkspaceUsageMetrics is current workspace consumption for billing.
type WorkspaceUsageMetrics struct {
	WorkspaceMembers int64   `json:"workspaceMembers"`
	DocumentsCount   int64   `json:"documentsCount"`
	StorageGB        float64 `json:"storageGB"`
	AITokens         int64   `json:"aiTokens"`
	AICost           float64 `json:"aiCost"`
}

// UsageSnapshotDeps holds repositories needed to load workspace usage.
type UsageSnapshotDeps struct {
	WorkspaceUserRepository     db.WorkspaceUserRepository
	WorkspaceDocumentRepository db.WorkspaceDocumentRepository
	UsageRepository             db.UsageRepository
}

// NewUsage returns a usage module.
func NewUsage() *Usage {
	return &Usage{}
}

// MembersCount returns the number of members in a workspace.
func (u *Usage) MembersCount(workspaceUsers db.WorkspaceUserRepository, workspaceId db.Id) (int64, error) {
	return workspaceUsers.CountByWorkspaceId(workspaceId)
}

// DocumentsCount returns the number of documents in a workspace.
func (u *Usage) DocumentsCount(documents db.WorkspaceDocumentRepository, workspaceId db.Id) (int64, error) {
	return documents.CountByWorkspaceId(workspaceId)
}

// StorageUsed returns total document storage in bytes for a workspace.
func (u *Usage) StorageUsed(documents db.WorkspaceDocumentRepository, workspaceId db.Id) (int64, error) {
	return documents.SumSizeByWorkspaceId(workspaceId)
}

// GetWorkspaceUsage returns all billing usage metrics for a workspace.
func (u *Usage) GetWorkspaceUsage(deps UsageSnapshotDeps, workspaceId db.Id) (*WorkspaceUsageMetrics, error) {
	members, err := u.MembersCount(deps.WorkspaceUserRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	documents, err := u.DocumentsCount(deps.WorkspaceDocumentRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	storageBytes, err := u.StorageUsed(deps.WorkspaceDocumentRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	aiTokens, err := u.AITokens(deps.UsageRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	aiCost, err := u.AICost(deps.UsageRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	return &WorkspaceUsageMetrics{
		WorkspaceMembers: members,
		DocumentsCount:   documents,
		StorageGB:        float64(storageBytes) / bytesPerGB,
		AITokens:         aiTokens,
		AICost:           aiCost,
	}, nil
}

// AITokens returns AI token usage for the current calendar month.
func (u *Usage) AITokens(usage db.UsageRepository, workspaceId db.Id) (int64, error) {
	start, _ := util.CurrentMonthPeriod()
	return usage.GetQuantityByPeriod(workspaceId, db.UsageTypeAITokens, start)
}

// AICost returns AI spend in USD for the current calendar month.
func (u *Usage) AICost(usage db.UsageRepository, workspaceId db.Id) (float64, error) {
	start, _ := util.CurrentMonthPeriod()
	cost, err := usage.GetQuantityByPeriod(workspaceId, db.UsageTypeAICost, start)
	if err != nil {
		return 0, err
	}

	return float64(cost) / ai.USDToNano, nil
}

// IncrementAIUsage adds AI token and cost usage for the current calendar month.
func (u *Usage) IncrementAIUsage(usage db.UsageRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, tokens, cost int64) error {
	if tokens == 0 && cost == 0 {
		return nil
	}

	start, end := util.CurrentMonthPeriod()

	if tokens != 0 {
		if err := usage.IncrementByPeriod(
			workspaceId,
			db.UsageTypeAITokens,
			start,
			end,
			tokens,
			db.UsageUnitTokens,
		); err != nil {
			return err
		}

		if err := subscriptions.ConsumeTokens(workspaceId, tokens); err != nil {
			return err
		}
	}

	if cost != 0 {
		if err := usage.IncrementByPeriod(
			workspaceId,
			db.UsageTypeAICost,
			start,
			end,
			cost,
			db.UsageUnitNanoUSD,
		); err != nil {
			return err
		}
	}

	return nil
}
