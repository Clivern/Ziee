// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package action

import (
	"context"

	"github.com/clivern/ziee/pkg/github/policy"
)

// Apply runs a plan against GitHub (or a test double).
func Apply(ctx context.Context, client Client, repo Repo, plan Plan) error {
	for _, item := range plan.Actions {
		var err error

		switch item.Kind {
		case policy.AddLabels:
			err = client.AddLabels(ctx, repo, item.Labels)
		case policy.RemoveLabels:
			err = client.RemoveLabels(ctx, repo, item.Labels)
		case policy.Assign:
			err = client.Assign(ctx, repo, item.Users)
		case policy.Unassign:
			err = client.Unassign(ctx, repo, item.Users)
		case policy.Comment:
			err = client.Comment(ctx, repo, item.Body)
		case policy.Close:
			err = client.Close(ctx, repo)
		case policy.Reopen:
			err = client.Reopen(ctx, repo)
		case policy.RequestReviewers:
			err = client.RequestReviewers(ctx, repo, item.Users)
		case policy.RequestReviewTeams:
			err = client.RequestReviewTeams(ctx, repo, item.Users)
		case policy.BlockAuthor:
			err = client.BlockAuthor(ctx, repo, item.Users)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
