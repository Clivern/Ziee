// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"context"
	"net/http"

	"github.com/actx0/ziee/pkg/github/response"
	"github.com/actx0/ziee/pkg/github/sender"
)

// NewComment creates an issue comment.
func (c *Client) NewComment(ctx context.Context, issueID int, body string) (response.CreatedComment, error) {
	var out response.CreatedComment
	err := c.call(ctx, callOptions{
		method:       http.MethodPost,
		path:         c.repoPath("/issues/%d/comments", issueID),
		body:         sender.Comment{Body: body},
		expectedCode: http.StatusCreated,
	}, &out)
	return out, err
}
