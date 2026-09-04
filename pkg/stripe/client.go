// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package stripe

import (
	"github.com/spf13/viper"
	stripesdk "github.com/stripe/stripe-go/v82"
)

// Client wraps the Stripe SDK using ziee billing configuration.
type Client struct {
	api    *stripesdk.Client
	config Config
}

// New returns a Stripe client loaded from app.billing config.
func New() (*Client, error) {
	if viper.GetString("app.billing.provider") != "stripe" {
		return nil, ErrBillingDisabled
	}

	return &Client{
		api: stripesdk.NewClient(viper.GetString("app.billing.secret_key")),
		config: Config{
			SecretKey:     viper.GetString("app.billing.secret_key"),
			WebhookSecret: viper.GetString("app.billing.webhook_secret"),
			ProductId:     viper.GetString("app.billing.product_id"),
			Prices: Prices{
				Hobby:   viper.GetString("app.billing.prices.hobby"),
				Pro:     viper.GetString("app.billing.prices.pro"),
				Growth:  viper.GetString("app.billing.prices.growth"),
				Starter: viper.GetString("app.billing.prices.starter"),
			},
			URLs: URLs{
				CheckoutSuccess: viper.GetString("app.billing.urls.checkout_success"),
				CheckoutCancel:  viper.GetString("app.billing.urls.checkout_cancel"),
				PortalReturn:    viper.GetString("app.billing.urls.portal_return"),
			},
		},
	}, nil
}

// Config returns the billing configuration used by this client.
func (c *Client) Config() Config {
	return c.config
}
