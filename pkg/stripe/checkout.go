// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package stripe

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	stripesdk "github.com/stripe/stripe-go/v82"
)

// CheckoutOptions is input for creating a subscription Checkout Session.
type CheckoutOptions struct {
	WorkspaceId   string
	Plan          string
	CustomerId    string
	CustomerEmail string
	SuccessURL    string
	CancelURL     string
}

// CreateCheckoutSession starts a Stripe Checkout subscription flow.
func (c *Client) CreateCheckoutSession(ctx context.Context, opts CheckoutOptions) (*stripesdk.CheckoutSession, error) {
	priceId := c.config.PriceId(opts.Plan)
	if lo.IsEmpty(priceId) {
		return nil, ErrInvalidPlan
	}

	params := &stripesdk.CheckoutSessionCreateParams{
		Mode:       stripesdk.String(string(stripesdk.CheckoutSessionModeSubscription)),
		SuccessURL: stripesdk.String(opts.SuccessURL),
		CancelURL:  stripesdk.String(opts.CancelURL),
		LineItems: []*stripesdk.CheckoutSessionCreateLineItemParams{
			{
				Price:    stripesdk.String(priceId),
				Quantity: stripesdk.Int64(1),
			},
		},
		Metadata: map[string]string{
			"plan":          opts.Plan,
			"workspaceId":   opts.WorkspaceId,
			"customerEmail": opts.CustomerEmail,
		},
		SubscriptionData: &stripesdk.CheckoutSessionCreateSubscriptionDataParams{
			Metadata: map[string]string{
				"plan":          opts.Plan,
				"workspaceId":   opts.WorkspaceId,
				"customerEmail": opts.CustomerEmail,
			},
		},
	}

	if lo.IsNotEmpty(opts.CustomerId) {
		params.Customer = stripesdk.String(opts.CustomerId)
	} else if lo.IsNotEmpty(opts.CustomerEmail) {
		params.CustomerEmail = stripesdk.String(opts.CustomerEmail)
	}

	session, err := c.api.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create checkout session: %w", err)
	}

	return session, nil
}
