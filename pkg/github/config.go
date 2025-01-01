// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"net/http"
	"time"
)

const defaultAPIURL = "https://api.github.com"

// Config holds credentials for a GitHub repository client.
type Config struct {
	Token      string
	Owner      string
	Repository string
}

// Client calls the GitHub REST API for a single repository.
type Client struct {
	token      string
	owner      string
	repository string
	baseURL    string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client used for API requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL overrides the GitHub API base URL (for tests or GitHub Enterprise).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// New returns a Client for the given repository.
func New(cfg Config, opts ...Option) *Client {
	c := &Client{
		token:      cfg.Token,
		owner:      cfg.Owner,
		repository: cfg.Repository,
		baseURL:    defaultAPIURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Owner returns the repository owner configured on the client.
func (c *Client) Owner() string {
	return c.owner
}

// Repository returns the repository name configured on the client.
func (c *Client) Repository() string {
	return c.repository
}
