// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clivern/ziee/pkg/hmac"

	"github.com/stretchr/testify/assert"
)

func TestUnitWebhookClient(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		client := New(Config{URL: "https://example.com/hooks/ziee", Secret: "whsec_test"})
		assert.Equal(t, "https://example.com/hooks/ziee", client.config.URL)
		assert.Equal(t, "whsec_test", client.config.Secret)
	})

	t.Run("Send", func(t *testing.T) {
		secret := "whsec_testsecret"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			assert.Contains(t, r.Header.Get("Content-Type"), "application/json")
			assert.True(t, hmac.Verify(secret, body, r.Header.Get(hmac.Header)))
			assert.Contains(t, string(body), `"event":"ziee.deploy"`)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := New(Config{URL: server.URL, Secret: secret})
		assert.NoError(t, client.Send(context.Background(), map[string]string{
			"event":   "ziee.deploy",
			"version": "abc123",
		}))
	})
}
