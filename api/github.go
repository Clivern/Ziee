// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/ziee/locale"
	"github.com/clivern/ziee/pkg/github/webhook"
	"github.com/clivern/ziee/pkg/util"

	"github.com/rs/zerolog/log"
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

	log.Info().
		Str("event", delivery.Event).
		Str("deliveryId", delivery.ID).
		Int("payloadBytes", len(delivery.Body)).
		Msg("GitHub webhook received")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"received": true,
	})
}
