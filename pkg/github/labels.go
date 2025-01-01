// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"context"
	"net/http"

	"github.com/actx0/ziee/pkg/github/response"
	"github.com/actx0/ziee/pkg/github/sender"
)

// CreateLabel creates a repository label.
func (c *Client) CreateLabel(ctx context.Context, name, color string) (response.Label, error) {
	var out response.Label
	err := c.call(ctx, callOptions{
		method:       http.MethodPost,
		path:         c.repoPath("/labels"),
		body:         sender.Label{Name: name, Color: color},
		expectedCode: http.StatusCreated,
	}, &out)
	return out, err
}

// UpdateLabel updates a repository label.
func (c *Client) UpdateLabel(ctx context.Context, currentName, name, color string) (response.Label, error) {
	var out response.Label
	err := c.call(ctx, callOptions{
		method:       http.MethodPatch,
		path:         c.repoPath("/labels/%s", currentName),
		body:         sender.Label{Name: name, Color: color},
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// DeleteLabel deletes a repository label.
func (c *Client) DeleteLabel(ctx context.Context, name string) error {
	return c.callNoContent(ctx, http.MethodDelete, c.repoPath("/labels/%s", name), http.StatusNoContent)
}

// GetRepositoryLabels lists repository labels.
func (c *Client) GetRepositoryLabels(ctx context.Context) ([]response.Label, error) {
	var out []response.Label
	err := c.call(ctx, callOptions{
		method:       http.MethodGet,
		path:         c.repoPath("/labels"),
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// GetRepositoryIssueLabels lists labels on an issue.
func (c *Client) GetRepositoryIssueLabels(ctx context.Context, issueID int) ([]response.Label, error) {
	var out []response.Label
	err := c.call(ctx, callOptions{
		method:       http.MethodGet,
		path:         c.repoPath("/issues/%d/labels", issueID),
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// GetLabel returns a single repository label.
func (c *Client) GetLabel(ctx context.Context, name string) (response.Label, error) {
	var out response.Label
	err := c.call(ctx, callOptions{
		method:       http.MethodGet,
		path:         c.repoPath("/labels/%s", name),
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// RemoveLabelFromIssue removes a label from an issue.
func (c *Client) RemoveLabelFromIssue(ctx context.Context, issueID int, labelName string) error {
	return c.callNoContent(
		ctx,
		http.MethodDelete,
		c.repoPath("/issues/%d/labels/%s", issueID, labelName),
		http.StatusNoContent,
	)
}

// RemoveAllLabelsFromIssue removes all labels from an issue.
func (c *Client) RemoveAllLabelsFromIssue(ctx context.Context, issueID int) error {
	return c.callNoContent(ctx, http.MethodDelete, c.repoPath("/issues/%d/labels", issueID), http.StatusNoContent)
}

// GetRepositoryMilestoneLabels lists labels for every issue in a milestone.
func (c *Client) GetRepositoryMilestoneLabels(ctx context.Context, milestoneID int) ([]response.Label, error) {
	var out []response.Label
	err := c.call(ctx, callOptions{
		method:       http.MethodGet,
		path:         c.repoPath("/milestones/%d/labels", milestoneID),
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// AddLabelsToIssue adds labels to an issue.
func (c *Client) AddLabelsToIssue(ctx context.Context, issueID int, labels []string) ([]response.Label, error) {
	var out []response.Label
	err := c.call(ctx, callOptions{
		method:       http.MethodPost,
		path:         c.repoPath("/issues/%d/labels", issueID),
		body:         labels,
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// ReplaceAllLabelsForIssue replaces all labels on an issue.
func (c *Client) ReplaceAllLabelsForIssue(ctx context.Context, issueID int, labels []string) ([]response.Label, error) {
	var out []response.Label
	err := c.call(ctx, callOptions{
		method:       http.MethodPut,
		path:         c.repoPath("/issues/%d/labels", issueID),
		body:         labels,
		expectedCode: http.StatusOK,
	}, &out)
	return out, err
}

// GetPullRequestLabels lists labels on a pull request.
func (c *Client) GetPullRequestLabels(ctx context.Context, pullRequestID int) ([]response.Label, error) {
	return c.GetRepositoryIssueLabels(ctx, pullRequestID)
}

// RemoveLabelFromPullRequest removes a label from a pull request.
func (c *Client) RemoveLabelFromPullRequest(ctx context.Context, pullRequestID int, labelName string) error {
	return c.RemoveLabelFromIssue(ctx, pullRequestID, labelName)
}

// RemoveAllLabelsFromPullRequest removes all labels from a pull request.
func (c *Client) RemoveAllLabelsFromPullRequest(ctx context.Context, pullRequestID int) error {
	return c.RemoveAllLabelsFromIssue(ctx, pullRequestID)
}

// AddLabelsToPullRequest adds labels to a pull request.
func (c *Client) AddLabelsToPullRequest(ctx context.Context, pullRequestID int, labels []string) ([]response.Label, error) {
	return c.AddLabelsToIssue(ctx, pullRequestID, labels)
}

// ReplaceAllLabelsForPullRequest replaces all labels on a pull request.
func (c *Client) ReplaceAllLabelsForPullRequest(ctx context.Context, pullRequestID int, labels []string) ([]response.Label, error) {
	return c.ReplaceAllLabelsForIssue(ctx, pullRequestID, labels)
}
