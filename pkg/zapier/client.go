// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package zapier

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

// Client posts events to a Zapier catch hook.
type Client struct {
	config Config
	http   *req.Client
}

// New returns a Zapier client for a customer's catch hook.
func New(cfg Config) *Client {
	return &Client{
		config: cfg,
		http:   req.C(),
	}
}

// Trigger sends a JSON payload to the catch hook. Zapier turns each field into a trigger field.
func (c *Client) Trigger(ctx context.Context, payload any) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post(c.config.WebhookURL)
	if err != nil {
		return fmt.Errorf("zapier request: %w", err)
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("zapier request: status %d: %s", resp.StatusCode, resp.Bytes())
	}

	log.Info().
		Str("provider", "zapier").
		Msg("Zapier hook triggered")

	return nil
}
