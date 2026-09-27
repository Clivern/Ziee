// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package zapier

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitZapierClient(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		client := New(Config{WebhookURL: "https://hooks.zapier.com/hooks/catch/1/abc"})
		assert.Equal(t, "https://hooks.zapier.com/hooks/catch/1/abc", client.config.WebhookURL)
	})

	t.Run("Trigger", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			assert.Contains(t, r.Header.Get("Content-Type"), "application/json")
			assert.Contains(t, string(body), `"event":"ziee.deploy"`)
			assert.Contains(t, string(body), `"version":"abc123"`)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		}))
		defer server.Close()

		client := New(Config{WebhookURL: server.URL})
		assert.NoError(t, client.Trigger(context.Background(), map[string]string{
			"event":   "ziee.deploy",
			"version": "abc123",
		}))
	})
}
