// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/actx0/ziee/db"
	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/middleware"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/stripe"
	"github.com/actx0/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

// GetBillingStatusAction returns the current workspace billing state.
func GetBillingStatusAction(w http.ResponseWriter, r *http.Request) {
	workspaceId := chi.URLParam(r, "workspaceId")
	log.Info().
		Str("workspaceId", workspaceId).
		Msg("Getting billing status")

	bm := module.NewBilling(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		module.Usage{},
	)

	status, err := bm.GetBillingStatus(db.Id(workspaceId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrBillingSubscriptionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_subscription_not_found"),
			})
		default:
			log.Error().Err(err).Msg("Failed to get billing status")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_billing_status"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, status)
}

// GetBillingUsageAction returns workspace usage and plan limits for billing.
func GetBillingUsageAction(w http.ResponseWriter, r *http.Request) {
	workspaceId := chi.URLParam(r, "workspaceId")
	log.Info().
		Str("workspaceId", workspaceId).
		Msg("Getting billing usage")

	bm := module.NewBilling(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		module.Usage{},
	)

	usage, err := bm.GetBillingUsage(db.Id(workspaceId), module.UsageSnapshotDeps{
		WorkspaceUserRepository:     db.NewWorkspaceUserRepository(db.GetDB()),
		WorkspaceDocumentRepository: db.NewWorkspaceDocumentRepository(db.GetDB()),
		UsageRepository:             db.NewUsageRepository(db.GetDB()),
		SubscriptionRepository:      db.NewSubscriptionRepository(db.GetDB()),
	})
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrBillingSubscriptionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_subscription_not_found"),
			})
		default:
			log.Error().Err(err).Msg("Failed to get billing usage")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_billing_usage"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, usage)
}

// CreateBillingCheckoutAction creates a Stripe Checkout session for a paid plan.
func CreateBillingCheckoutAction(w http.ResponseWriter, r *http.Request) {
	var req module.BillingCheckoutRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	workspaceId := chi.URLParam(r, "workspaceId")
	log.Info().
		Str("workspaceId", workspaceId).
		Str("userId", user.Id.String()).
		Str("plan", req.Plan).
		Msg("New billing checkout request")

	bm := module.NewBilling(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		module.Usage{},
	)
	session, err := bm.CreateCheckoutSession(
		r.Context(),
		db.Id(workspaceId),
		user,
		&req,
		fmt.Sprintf("%s/billing?checkout=success", viper.GetString("app.url")),
		fmt.Sprintf("%s/billing?checkout=cancel", viper.GetString("app.url")),
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrBillingFreePlan):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "billing_free_plan_default"),
			})
		case errors.Is(err, module.ErrBillingSubscriptionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_subscription_not_found"),
			})
		case errors.Is(err, stripe.ErrInvalidPlan):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_billing_plan"),
			})
		case errors.Is(err, stripe.ErrBillingDisabled):
			util.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
				"errorMessage": locale.TR(r, "stripe_billing_not_configured"),
			})
		default:
			log.Error().Err(err).Msg("Failed to create billing checkout session")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_billing_checkout_session"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, session)
}

// CreateBillingPortalAction creates a Stripe Billing Portal session.
func CreateBillingPortalAction(w http.ResponseWriter, r *http.Request) {
	workspaceId := chi.URLParam(r, "workspaceId")
	log.Info().
		Str("workspaceId", workspaceId).
		Msg("New billing portal request")

	bm := module.NewBilling(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		module.Usage{},
	)
	session, err := bm.CreatePortalSession(
		r.Context(),
		db.Id(workspaceId),
		fmt.Sprintf("%s/billing", viper.GetString("app.url")),
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrBillingPortalUnavailable):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "stripe_billing_available_after_upgrade"),
			})
		case errors.Is(err, module.ErrBillingSubscriptionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_subscription_not_found"),
			})
		case errors.Is(err, stripe.ErrBillingDisabled):
			util.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
				"errorMessage": locale.TR(r, "stripe_billing_not_configured"),
			})
		default:
			log.Error().Err(err).Msg("Failed to create billing portal session")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_billing_portal_session"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, session)
}

// StripeWebhookAction receives Stripe billing webhook events.
func StripeWebhookAction(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		log.Warn().Err(err).Msg("Stripe webhook rejected: failed to read payload")
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_webhook_payload"),
		})
		return
	}

	signature := r.Header.Get("Stripe-Signature")
	log.Info().
		Int("payloadBytes", len(payload)).
		Bool("hasSignature", lo.IsNotEmpty(signature)).
		Msg("Stripe webhook received")

	bm := module.NewBilling(
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewSubscriptionRepository(db.GetDB()),
		module.Usage{},
	)

	err = bm.HandleWebhook(payload, signature)
	if err != nil {
		switch {
		case errors.Is(err, stripe.ErrBillingDisabled), errors.Is(err, stripe.ErrWebhookNotConfigured):
			log.Warn().Err(err).Msg("Stripe webhook rejected: billing not configured")
			util.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
				"errorMessage": locale.TR(r, "stripe_billing_not_configured"),
			})
		default:
			log.Error().Err(err).Msg("Failed to handle Stripe webhook")
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_stripe_webhook"),
			})
		}
		return
	}

	log.Info().Msg("Stripe webhook handled")

	util.WriteJSON(w, http.StatusOK, module.BillingWebhookResponse{
		Received: true,
	})
}
