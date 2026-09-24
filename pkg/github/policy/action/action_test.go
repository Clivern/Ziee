// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package action

import (
	"context"
	"errors"
	"testing"

	"github.com/clivern/ziee/pkg/github/policy"

	"github.com/stretchr/testify/assert"
)

func TestUnitApply(t *testing.T) {
	repo := Repo{Owner: "clivern", Name: "ziee", Number: 12}
	client := &stubClient{}

	err := Apply(context.Background(), client, repo, Plan{Actions: []Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
		{Kind: policy.RemoveLabels, Labels: []string{"needs-triage"}},
		{Kind: policy.Assign, Users: []string{"maya"}},
		{Kind: policy.Unassign, Users: []string{"bot"}},
		{Kind: policy.Comment, Body: "Triaged this issue."},
		{Kind: policy.Close},
		{Kind: policy.Reopen},
		{Kind: policy.RequestReviewers, Users: []string{"maya"}},
		{Kind: policy.RequestReviewTeams, Users: []string{"sre"}},
		{Kind: policy.BlockAuthor, Users: []string{"spammer"}},
	}})

	assert.NoError(t, err)
	assert.Equal(t, repo, client.repo)
	assert.Equal(t, []Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
		{Kind: policy.RemoveLabels, Labels: []string{"needs-triage"}},
		{Kind: policy.Assign, Users: []string{"maya"}},
		{Kind: policy.Unassign, Users: []string{"bot"}},
		{Kind: policy.Comment, Body: "Triaged this issue."},
		{Kind: policy.Close},
		{Kind: policy.Reopen},
		{Kind: policy.RequestReviewers, Users: []string{"maya"}},
		{Kind: policy.RequestReviewTeams, Users: []string{"sre"}},
		{Kind: policy.BlockAuthor, Users: []string{"spammer"}},
	}, client.got)
}

func TestUnitApplyEmpty(t *testing.T) {
	assert.NoError(t, Apply(context.Background(), nil, Repo{}, Plan{}))
}

func TestUnitApplyStops(t *testing.T) {
	client := &stubClient{err: errors.New("github")}

	err := Apply(context.Background(), client, Repo{}, Plan{Actions: []Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
		{Kind: policy.Close},
	}})

	assert.EqualError(t, err, "github")
	assert.Equal(t, []Action{
		{Kind: policy.AddLabels, Labels: []string{"bug"}},
	}, client.got)
}

type stubClient struct {
	repo Repo
	got  []Action
	err  error
}

func (s *stubClient) record(repo Repo, item Action) error {
	s.repo = repo
	s.got = append(s.got, item)

	return s.err
}

func (s *stubClient) AddLabels(ctx context.Context, repo Repo, labels []string) error {
	return s.record(repo, Action{Kind: policy.AddLabels, Labels: labels})
}

func (s *stubClient) RemoveLabels(ctx context.Context, repo Repo, labels []string) error {
	return s.record(repo, Action{Kind: policy.RemoveLabels, Labels: labels})
}

func (s *stubClient) Assign(ctx context.Context, repo Repo, users []string) error {
	return s.record(repo, Action{Kind: policy.Assign, Users: users})
}

func (s *stubClient) Unassign(ctx context.Context, repo Repo, users []string) error {
	return s.record(repo, Action{Kind: policy.Unassign, Users: users})
}

func (s *stubClient) Comment(ctx context.Context, repo Repo, body string) error {
	return s.record(repo, Action{Kind: policy.Comment, Body: body})
}

func (s *stubClient) Close(ctx context.Context, repo Repo) error {
	return s.record(repo, Action{Kind: policy.Close})
}

func (s *stubClient) Reopen(ctx context.Context, repo Repo) error {
	return s.record(repo, Action{Kind: policy.Reopen})
}

func (s *stubClient) RequestReviewers(ctx context.Context, repo Repo, users []string) error {
	return s.record(repo, Action{Kind: policy.RequestReviewers, Users: users})
}

func (s *stubClient) RequestReviewTeams(ctx context.Context, repo Repo, teams []string) error {
	return s.record(repo, Action{Kind: policy.RequestReviewTeams, Users: teams})
}

func (s *stubClient) BlockAuthor(ctx context.Context, repo Repo, users []string) error {
	return s.record(repo, Action{Kind: policy.BlockAuthor, Users: users})
}
