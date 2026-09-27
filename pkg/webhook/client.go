// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/clivern/ziee/pkg/hmac"

	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

// Client posts JSON to a customer webhook and signs the body with HMAC-SHA256.
type Client struct {
	config Config
	http   *req.Client
}

// New returns a webhook client for a customer's endpoint.
func New(cfg Config) *Client {
	return &Client{
		config: cfg,
		http:   req.C(),
	}
}

// Send posts payload as JSON and sets X-Ziee-Signature-256 over those exact bytes.
func (c *Client) Send(ctx context.Context, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook body: %w", err)
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader(hmac.Header, hmac.Sign(c.config.Secret, body)).
		SetBody(body).
		Post(c.config.URL)
	if err != nil {
		return fmt.Errorf("webhook request: %w", err)
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("webhook request: status %d: %s", resp.StatusCode, resp.Bytes())
	}

	log.Info().
		Str("provider", "webhook").
		Msg("Webhook delivered")

	return nil
}
