// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ReactionPlusOne  = "+1"
	ReactionMinusOne = "-1"
	ReactionLaugh    = "laugh"
	ReactionConfused = "confused"
	ReactionHeart    = "heart"
	ReactionHooray   = "hooray"
	ReactionRocket   = "rocket"
	ReactionEyes     = "eyes"
)

// Reaction is a GitHub reaction on a comment or issue.
type Reaction struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

// CreateIssueCommentReaction adds a reaction to an issue or pull request comment.
func (a *App) CreateIssueCommentReaction(ctx context.Context, installationID int64, owner, repo string, commentID int64, content string) (*Reaction, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var reaction Reaction
	path := fmt.Sprintf("%s/repos/%s/%s/issues/comments/%d/reactions", a.apiURL, owner, repo, commentID)
	err = Call(ctx, http.MethodPost, path, token.Token, GetHeaders(), map[string]string{"content": content}, &reaction)
	if err != nil {
		return nil, err
	}

	return &reaction, nil
}
