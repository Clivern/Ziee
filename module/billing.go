// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/stripe"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	stripesdk "github.com/stripe/stripe-go/v82"
)

const (
	BillingProviderStripe = "stripe"
	BillingStatusActive   = "active"
	BillingStatusCanceled = "canceled"
)

var (
	ErrBillingPortalUnavailable    = errors.New("stripe billing portal is unavailable for this workspace")
	ErrBillingFreePlan             = errors.New("hobby is the default free plan and does not use checkout")
	ErrBillingSubscriptionNotFound = errors.New("workspace subscription not found")
)

// Billing coordinates local subscription state with Stripe Checkout and Portal.
type Billing struct {
	WorkspaceRepository    db.WorkspaceRepository
	SubscriptionRepository db.SubscriptionRepository
	Usage                  Usage
}

// BillingStatusResponse is the billing state exposed to the UI.
type BillingStatusResponse struct {
	Id                 db.Id   `json:"id"`
	WorkspaceId        db.Id   `json:"workspaceId"`
	Plan               string  `json:"plan"`
	Status             string  `json:"status"`
	Provider           *string `json:"provider"`
	ProviderCustomerId string  `json:"providerCustomerId"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

// BillingCheckoutRequest is the body for creating a Stripe Checkout session.
type BillingCheckoutRequest struct {
	Plan string `json:"plan" validate:"required" label:"Plan"`
}

// BillingSessionResponse contains a Stripe redirect URL.
type BillingSessionResponse struct {
	URL string `json:"url"`
}

// BillingWebhookResponse is returned after a Stripe webhook is handled.
type BillingWebhookResponse struct {
	Received bool `json:"received"`
}

// PlanUsageLimits are the consumption caps for a billing plan.
type PlanUsageLimits struct {
	APICalls         int64   `json:"apiCalls"`
	WorkspaceMembers int64   `json:"workspaceMembers"`
	DocumentsCount   int64   `json:"documentsCount"`
	PromptsCount     int64   `json:"promptsCount"`
	StorageGB        float64 `json:"storageGB"`
	HistoryRecords   int64   `json:"historyRecords"`
	MemoryRecords    int64   `json:"memoryRecords"`
	AITokens         int64   `json:"aiTokens"`
}

// BillingUsageResponse is workspace usage and limits for the billing page.
type BillingUsageResponse struct {
	Plan        string                `json:"plan"`
	PeriodReset string                `json:"periodReset"`
	Used        WorkspaceUsageMetrics `json:"used"`
	Limits      PlanUsageLimits       `json:"limits"`
}

var planUsageLimits = map[string]PlanUsageLimits{
	stripe.PlanHobby: {
		APICalls:         10_000,
		WorkspaceMembers: 3,
		DocumentsCount:   100,
		PromptsCount:     50,
		StorageGB:        5,
		HistoryRecords:   1_000,
		MemoryRecords:    500,
		AITokens:         500_000,
	},
	stripe.PlanStarter: {
		APICalls:         50_000,
		WorkspaceMembers: 10,
		DocumentsCount:   500,
		PromptsCount:     250,
		StorageGB:        25,
		HistoryRecords:   10_000,
		MemoryRecords:    5_000,
		AITokens:         2_000_000,
	},
	stripe.PlanGrowth: {
		APICalls:         200_000,
		WorkspaceMembers: 25,
		DocumentsCount:   2_000,
		PromptsCount:     1_000,
		StorageGB:        100,
		HistoryRecords:   50_000,
		MemoryRecords:    25_000,
		AITokens:         10_000_000,
	},
	stripe.PlanPro: {
		APICalls:         1_000_000,
		WorkspaceMembers: 100,
		DocumentsCount:   10_000,
		PromptsCount:     5_000,
		StorageGB:        500,
		HistoryRecords:   500_000,
		MemoryRecords:    250_000,
		AITokens:         50_000_000,
	},
}

// NewBilling creates a billing module.
func NewBilling(workspaces db.WorkspaceRepository, subscriptions db.SubscriptionRepository, usage Usage) *Billing {
	return &Billing{
		WorkspaceRepository:    workspaces,
		SubscriptionRepository: subscriptions,
		Usage:                  usage,
	}
}

// GetBillingStatus returns billing status for a workspace.
func (b *Billing) GetBillingStatus(workspaceId db.Id) (*BillingStatusResponse, error) {
	workspace, err := b.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}

	err = b.UpdateHobbyPlanPeriods(subscription)
	if err != nil {
		return nil, err
	}

	return &BillingStatusResponse{
		Id:                 subscription.Id,
		WorkspaceId:        subscription.WorkspaceId,
		Plan:               subscription.Plan,
		Status:             subscription.Status,
		Provider:           subscription.Provider,
		ProviderCustomerId: lo.FromPtr(subscription.ProviderCustomerId),
		CreatedAt:          subscription.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:          subscription.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetBillingUsage returns usage and plan limits for a workspace.
func (b *Billing) GetBillingUsage(workspaceId db.Id, deps UsageSnapshotDeps) (*BillingUsageResponse, error) {
	workspace, err := b.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}

	err = b.UpdateHobbyPlanPeriods(subscription)
	if err != nil {
		return nil, err
	}

	used, err := b.Usage.GetWorkspaceUsage(deps, workspaceId)
	if err != nil {
		return nil, err
	}

	limits, ok := planUsageLimits[subscription.Plan]
	if !ok {
		limits = planUsageLimits[stripe.PlanHobby]
	}

	_, periodEnd := util.CurrentMonthPeriod()

	return &BillingUsageResponse{
		Plan:        subscription.Plan,
		PeriodReset: periodEnd.UTC().Format(time.RFC3339),
		Used:        *used,
		Limits:      limits,
	}, nil
}

// CreateCheckoutSession starts Stripe Checkout for a paid plan.
func (b *Billing) CreateCheckoutSession(ctx context.Context, workspaceId db.Id, user *db.User, req *BillingCheckoutRequest, successURL string, cancelURL string) (*BillingSessionResponse, error) {
	workspace, err := b.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	if req.Plan == stripe.PlanHobby {
		return nil, ErrBillingFreePlan
	}

	client, err := stripe.New()
	if err != nil {
		return nil, err
	}
	if lo.IsEmpty(client.Config().PriceId(req.Plan)) {
		return nil, stripe.ErrInvalidPlan
	}

	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}

	session, err := client.CreateCheckoutSession(ctx, stripe.CheckoutOptions{
		WorkspaceId:   workspaceId.String(),
		Plan:          req.Plan,
		CustomerId:    lo.FromPtr(subscription.ProviderCustomerId),
		CustomerEmail: user.Email,
		SuccessURL:    successURL,
		CancelURL:     cancelURL,
	})
	if err != nil {
		return nil, err
	}

	return &BillingSessionResponse{URL: session.URL}, nil
}

// CreatePortalSession opens Stripe's customer portal.
func (b *Billing) CreatePortalSession(ctx context.Context, workspaceId db.Id, returnURL string) (*BillingSessionResponse, error) {
	workspace, err := b.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	client, err := stripe.New()
	if err != nil {
		return nil, err
	}
	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}
	customerId := lo.FromPtr(subscription.ProviderCustomerId)
	if lo.IsEmpty(customerId) {
		return nil, ErrBillingPortalUnavailable
	}

	session, err := client.CreatePortalSession(ctx, stripe.PortalOptions{
		CustomerId: customerId,
		ReturnURL:  returnURL,
	})
	if err != nil {
		return nil, err
	}

	return &BillingSessionResponse{URL: session.URL}, nil
}

// HandleWebhook updates local billing state from Stripe webhook events.
func (b *Billing) HandleWebhook(payload []byte, signature string) error {
	client, err := stripe.New()
	if err != nil {
		return err
	}

	event, err := client.ConstructWebhookEvent(payload, signature)
	if err != nil {
		return err
	}

	switch event.Type {
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
		var subscription stripesdk.Subscription
		err = json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return fmt.Errorf("decode subscription: %w", err)
		}
		return b.SyncStripeSubscription(client, &subscription, string(event.Type))
	default:
		return nil
	}
}

// SyncStripeSubscription syncs a Stripe subscription with the local database.
func (b *Billing) SyncStripeSubscription(client *stripe.Client, payload *stripesdk.Subscription, eventType string) error {
	if payload == nil {
		return nil
	}

	workspaceId := db.Id(payload.Metadata["workspaceId"])
	plan := payload.Metadata["plan"]

	if lo.IsEmpty(workspaceId) {
		return nil
	}

	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return err
	}
	if subscription == nil {
		log.Warn().
			Str("workspaceId", workspaceId.String()).
			Str("subscriptionId", payload.ID).
			Msg("Stripe subscription sync ignored: workspace subscription not found")
		return nil
	}

	status := string(payload.Status)
	if eventType == "customer.subscription.deleted" {
		status = BillingStatusCanceled
	}

	subscription.Provider = new(BillingProviderStripe)
	if payload.Customer != nil && lo.IsNotEmpty(payload.Customer.ID) {
		subscription.ProviderCustomerId = new(payload.Customer.ID)
	}
	if lo.IsNotEmpty(payload.ID) && status != BillingStatusCanceled {
		subscription.ProviderSubscriptionId = new(payload.ID)
	}

	if payload.Items != nil && len(payload.Items.Data) > 0 {
		item := payload.Items.Data[0]
		if item.CurrentPeriodStart > 0 {
			subscription.CurrentPeriodStart = new(time.Unix(item.CurrentPeriodStart, 0).UTC())
		}
		if item.CurrentPeriodEnd > 0 {
			subscription.CurrentPeriodEnd = new(time.Unix(item.CurrentPeriodEnd, 0).UTC())
		}
		if lo.IsEmpty(plan) && item.Price != nil {
			plan = client.Config().PlanForPriceId(item.Price.ID)
		}
	}

	if status == BillingStatusActive || status == BillingStatusCanceled {
		if lo.IsNotEmpty(plan) {
			subscription.Plan = plan
		}
		subscription.Status = status
	}

	err = b.SubscriptionRepository.Update(subscription)
	if err != nil {
		return err
	}

	log.Info().
		Str("workspaceId", subscription.WorkspaceId.String()).
		Str("subscriptionId", payload.ID).
		Str("plan", subscription.Plan).
		Str("status", subscription.Status).
		Msg("Stripe subscription synced")

	return nil
}

// UpdateHobbyPlanPeriods updates the hobby plan periods for a subscription.
func (b *Billing) UpdateHobbyPlanPeriods(subscription *db.Subscription) error {
	if subscription == nil || subscription.Plan != stripe.PlanHobby {
		return nil
	}

	now := time.Now().UTC()
	if subscription.CurrentPeriodEnd != nil && !subscription.CurrentPeriodEnd.Before(now) {
		return nil
	}

	periodStart := subscription.CurrentPeriodEnd.UTC()
	periodEnd := periodStart.AddDate(0, 1, 0)

	for !periodEnd.After(now) {
		periodStart = periodEnd
		periodEnd = periodStart.AddDate(0, 1, 0)
	}

	subscription.CurrentPeriodStart = new(periodStart)
	subscription.CurrentPeriodEnd = new(periodEnd)

	err := b.SubscriptionRepository.Update(subscription)

	if err != nil {
		return err
	}
	return nil
}

// GetWorkspaceSubscription gets a workspace subscription from the database.
func (b *Billing) GetWorkspaceSubscription(workspaceId db.Id) (*db.Subscription, error) {
	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}
	return subscription, nil
}
