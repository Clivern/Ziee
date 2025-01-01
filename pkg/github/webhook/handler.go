// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/actx0/ziee/pkg/github/event"
)

// Handler verifies and dispatches GitHub webhook payloads.
type Handler struct {
	Secret string
	Hooks  Hooks
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	delivery, err := ParseDelivery(r)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !delivery.VerifySignature(h.Secret) {
		log.Warn().Str("event", delivery.Event).Msg("github webhook signature verification failed")
		writeJSON(w, http.StatusUnauthorized, map[string]string{"status": "invalid signature"})
		return
	}

	log.Info().Str("event", delivery.Event).Str("delivery", delivery.ID).Msg("github webhook received")

	if err := h.dispatch(delivery.Event, delivery.Body); err != nil {
		log.Error().Err(err).Str("event", delivery.Event).Msg("github webhook handler failed")
		writeJSON(w, http.StatusInternalServerError, map[string]string{"status": "handler error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) dispatch(evt string, rawBody []byte) error {
	var actions Action
	var commands Commands

	switch evt {
	case "status":
		var payload event.Status
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.Status != nil {
			actions.RegisterStatusAction(wrapStatus(h.Hooks.Status))
		}
		actions.ExecuteStatusActions(payload)
	case "watch":
		var payload event.Watch
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.Watch != nil {
			actions.RegisterWatchAction(wrapWatch(h.Hooks.Watch))
		}
		actions.ExecuteWatchActions(payload)
	case "issues":
		var payload event.Issues
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.Issues != nil {
			actions.RegisterIssuesAction(wrapIssues(h.Hooks.Issues))
		}
		actions.ExecuteIssuesActions(payload)
		for name, fn := range h.Hooks.IssuesCommand {
			commands.RegisterIssuesAction(name, wrapIssuesCommand(fn))
		}
		commands.ExecuteIssuesActions(payload)
	case "push":
		var payload event.Push
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.Push != nil {
			actions.RegisterPushAction(wrapPush(h.Hooks.Push))
		}
		actions.ExecutePushActions(payload)
	case "issue_comment":
		var payload event.IssueComment
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.IssueComment != nil {
			actions.RegisterIssueCommentAction(wrapIssueComment(h.Hooks.IssueComment))
		}
		actions.ExecuteIssueCommentActions(payload)
		for name, fn := range h.Hooks.IssueCommentCommand {
			commands.RegisterIssueCommentAction(name, wrapIssueCommentCommand(fn))
		}
		commands.ExecuteIssueCommentActions(payload)
	case "create":
		var payload event.Create
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.Create != nil {
			actions.RegisterCreateAction(wrapCreate(h.Hooks.Create))
		}
		actions.ExecuteCreateActions(payload)
	case "label":
		var payload event.Label
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.Label != nil {
			actions.RegisterLabelAction(wrapLabel(h.Hooks.Label))
		}
		actions.ExecuteLabelActions(payload)
	case "delete":
		var payload event.Delete
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.Delete != nil {
			actions.RegisterDeleteAction(wrapDelete(h.Hooks.Delete))
		}
		actions.ExecuteDeleteActions(payload)
	case "milestone":
		var payload event.Milestone
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.Milestone != nil {
			actions.RegisterMilestoneAction(wrapMilestone(h.Hooks.Milestone))
		}
		actions.ExecuteMilestoneActions(payload)
	case "pull_request":
		var payload event.PullRequest
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.PullRequest != nil {
			actions.RegisterPullRequestAction(wrapPullRequest(h.Hooks.PullRequest))
		}
		actions.ExecutePullRequestActions(payload)
	case "pull_request_review":
		var payload event.PullRequestReview
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.PullRequestReview != nil {
			actions.RegisterPullRequestReviewAction(wrapPullRequestReview(h.Hooks.PullRequestReview))
		}
		actions.ExecutePullRequestReviewActions(payload)
	case "pull_request_review_comment":
		var payload event.PullRequestReviewComment
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return err
		}
		if h.Hooks.PullRequestReviewComment != nil {
			actions.RegisterPullRequestReviewCommentAction(wrapPullRequestReviewComment(h.Hooks.PullRequestReviewComment))
		}
		actions.ExecutePullRequestReviewCommentActions(payload)
	default:
		log.Info().Str("event", evt).Msg("unsupported github webhook event")
	}

	if h.Hooks.Raw != nil {
		var raw event.Raw
		raw.SetEvent(evt)
		raw.SetBody(string(rawBody))
		actions.RegisterRawAction(wrapRaw(h.Hooks.Raw))
		actions.ExecuteRawActions(raw)
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func wrapRaw(fn func(event.Raw) error) func(event.Raw) (bool, error) {
	return func(raw event.Raw) (bool, error) {
		err := fn(raw)
		return err == nil, err
	}
}

func wrapStatus(fn func(event.Status) error) func(event.Status) (bool, error) {
	return func(status event.Status) (bool, error) {
		err := fn(status)
		return err == nil, err
	}
}

func wrapWatch(fn func(event.Watch) error) func(event.Watch) (bool, error) {
	return func(watch event.Watch) (bool, error) {
		err := fn(watch)
		return err == nil, err
	}
}

func wrapIssues(fn func(event.Issues) error) func(event.Issues) (bool, error) {
	return func(issues event.Issues) (bool, error) {
		err := fn(issues)
		return err == nil, err
	}
}

func wrapIssueComment(fn func(event.IssueComment) error) func(event.IssueComment) (bool, error) {
	return func(issueComment event.IssueComment) (bool, error) {
		err := fn(issueComment)
		return err == nil, err
	}
}

func wrapPush(fn func(event.Push) error) func(event.Push) (bool, error) {
	return func(push event.Push) (bool, error) {
		err := fn(push)
		return err == nil, err
	}
}

func wrapCreate(fn func(event.Create) error) func(event.Create) (bool, error) {
	return func(create event.Create) (bool, error) {
		err := fn(create)
		return err == nil, err
	}
}

func wrapLabel(fn func(event.Label) error) func(event.Label) (bool, error) {
	return func(label event.Label) (bool, error) {
		err := fn(label)
		return err == nil, err
	}
}

func wrapDelete(fn func(event.Delete) error) func(event.Delete) (bool, error) {
	return func(del event.Delete) (bool, error) {
		err := fn(del)
		return err == nil, err
	}
}

func wrapMilestone(fn func(event.Milestone) error) func(event.Milestone) (bool, error) {
	return func(milestone event.Milestone) (bool, error) {
		err := fn(milestone)
		return err == nil, err
	}
}

func wrapPullRequest(fn func(event.PullRequest) error) func(event.PullRequest) (bool, error) {
	return func(pullRequest event.PullRequest) (bool, error) {
		err := fn(pullRequest)
		return err == nil, err
	}
}

func wrapPullRequestReview(fn func(event.PullRequestReview) error) func(event.PullRequestReview) (bool, error) {
	return func(pullRequestReview event.PullRequestReview) (bool, error) {
		err := fn(pullRequestReview)
		return err == nil, err
	}
}

func wrapPullRequestReviewComment(fn func(event.PullRequestReviewComment) error) func(event.PullRequestReviewComment) (bool, error) {
	return func(pullRequestReviewComment event.PullRequestReviewComment) (bool, error) {
		err := fn(pullRequestReviewComment)
		return err == nil, err
	}
}

func wrapIssuesCommand(fn func(event.Command, event.Issues) error) func(event.Command, event.Issues) (bool, error) {
	return func(command event.Command, issues event.Issues) (bool, error) {
		err := fn(command, issues)
		return err == nil, err
	}
}

func wrapIssueCommentCommand(fn func(event.Command, event.IssueComment) error) func(event.Command, event.IssueComment) (bool, error) {
	return func(command event.Command, issueComment event.IssueComment) (bool, error) {
		err := fn(command, issueComment)
		return err == nil, err
	}
}
