// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

const ChatPostMessageURL = "https://slack.com/api/chat.postMessage"

// Client sends Slack notifications for a customer workspace.
type Client struct {
	config Config
}

// Message is a Slack notification payload.
type Message struct {
	Text      string
	Channel   string
	Username  string
	IconEmoji string
}

type payload struct {
	Text      string `json:"text"`
	Channel   string `json:"channel,omitempty"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

type apiResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// New returns a Slack client for a customer's workspace.
func New(cfg Config) *Client {
	return &Client{
		config: cfg,
	}
}

// Notify sends a message using the customer's webhook or bot token.
func (c *Client) Notify(ctx context.Context, text string) error {
	msg := Message{
		Text:    text,
		Channel: c.config.Channel,
	}

	if c.config.WebhookURL != "" {
		return c.PostWebhook(ctx, msg)
	}

	return c.PostMessage(ctx, msg)
}

// PostWebhook posts a message to the customer's Slack incoming webhook.
func (c *Client) PostWebhook(ctx context.Context, msg Message) error {
	body, err := c.PostJSON(ctx, c.config.WebhookURL, "", msg)
	if err != nil {
		return err
	}

	if string(body) != "ok" {
		return fmt.Errorf("slack webhook: %s", body)
	}

	log.Info().
		Str("provider", "slack").
		Str("channel", msg.Channel).
		Msg("Slack notification sent")

	return nil
}

// PostMessage posts a message with the customer's bot token via chat.postMessage.
func (c *Client) PostMessage(ctx context.Context, msg Message) error {
	if msg.Channel == "" {
		msg.Channel = c.config.Channel
	}

	body, err := c.PostJSON(ctx, ChatPostMessageURL, c.config.Token, msg)
	if err != nil {
		return err
	}

	var result apiResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("slack decode response: %w", err)
	}

	if !result.OK {
		return fmt.Errorf("slack chat.postMessage: %s", result.Error)
	}

	log.Info().
		Str("provider", "slack").
		Str("channel", msg.Channel).
		Msg("Slack notification sent")

	return nil
}

// PostJSON posts a JSON payload to a Slack endpoint.
func (c *Client) PostJSON(ctx context.Context, endpoint, token string, msg Message) ([]byte, error) {
	raw, err := json.Marshal(payload{
		Text:      msg.Text,
		Channel:   msg.Channel,
		Username:  msg.Username,
		IconEmoji: msg.IconEmoji,
	})
	if err != nil {
		return nil, fmt.Errorf("slack encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("slack build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("slack request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("slack read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("slack request: status %d: %s", resp.StatusCode, body)
	}

	return body, nil
}
