// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/clivern/ziee/db"
	"github.com/clivern/ziee/locale"
	"github.com/clivern/ziee/middleware"
	"github.com/clivern/ziee/module"
	"github.com/clivern/ziee/pkg/github/webhook"
	"github.com/clivern/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

// GitHubWebhookAction receives GitHub App webhook deliveries.
func GitHubWebhookAction(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	delivery, err := webhook.ParseDelivery(r)
	if err != nil {
		log.Warn().Err(err).Msg("GitHub webhook rejected: failed to parse delivery")
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_webhook_payload"),
		})
		return
	}

	secret := viper.GetString("app.oauth.github.webhook_secret")
	if !delivery.VerifySignature(secret) {
		log.Warn().
			Str("event", delivery.Event).
			Str("deliveryId", delivery.ID).
			Msg("GitHub webhook rejected: invalid signature")
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "invalid_github_webhook"),
		})
		return
	}

	// TODO: Remove this after testing
	err = os.MkdirAll("events", 0o755)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GitHub webhook events dir")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "invalid_webhook_payload"),
		})
		return
	}

	path := filepath.Join(
		"events",
		fmt.Sprintf("gh-%s-%s.json", delivery.Event, delivery.ID),
	)

	var pretty bytes.Buffer
	err = json.Indent(&pretty, delivery.Body, "", "  ")
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to format GitHub webhook")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "invalid_webhook_payload"),
		})
		return
	}

	pretty.WriteByte('\n')
	err = os.WriteFile(path, pretty.Bytes(), 0o644)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to dump GitHub webhook")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "invalid_webhook_payload"),
		})
		return
	}
	// TODO: Remove this after testing

	im := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewWorkspaceGitHubRepoRepository(db.GetDB()),
	)

	err = im.HandleWebhook(delivery.Event, delivery.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("event", delivery.Event).
			Str("deliveryId", delivery.ID).
			Msg("Failed to persist GitHub installation")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "invalid_webhook_payload"),
		})
		return
	}

	log.Info().
		Str("event", delivery.Event).
		Str("deliveryId", delivery.ID).
		Str("path", path).
		Int("payloadBytes", len(delivery.Body)).
		Msg("GitHub webhook received")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"received": true,
	})
}

// ListGitHubInstallationsAction lists pending GitHub App installations for the signed-in GitHub user.
func ListGitHubInstallationsAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	im := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewWorkspaceGitHubRepoRepository(db.GetDB()),
	)

	installations, err := im.ListPending(lo.FromPtr(user.ProviderUserId))
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to list GitHub installations")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_list_github_installations"),
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"installations": installations,
	})
}

// AttachGitHubInstallationAction attaches a pending GitHub App installation to a workspace.
func AttachGitHubInstallationAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	id := chi.URLParam(r, "installationId")
	if lo.IsEmpty(id) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "github_installation_not_found"),
		})
		return
	}

	var req module.AttachInstallationRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	im := module.NewInstallation(
		db.NewGitHubInstallationRepository(db.GetDB()),
		db.NewWorkspaceGitHubRepoRepository(db.GetDB()),
	)

	err = im.Attach(db.Id(id), db.Id(req.WorkspaceId), lo.FromPtr(user.ProviderUserId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrInstallationNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "github_installation_not_found"),
			})
		default:
			log.Error().
				Err(err).
				Str("installationId", id).
				Str("workspaceId", req.WorkspaceId).
				Str("userId", user.Id.String()).
				Msg("Failed to attach GitHub installation")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_attach_github_installation"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"attached": true,
	})
}
