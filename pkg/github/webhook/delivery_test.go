// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeliveryVerifySignature(t *testing.T) {
	secret := "test-secret"
	body := []byte(`{"action":"opened"}`)

	t.Run("SHA256", func(t *testing.T) {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

		delivery := Delivery{
			Body:            body,
			SignatureSHA256: signature,
		}

		assert.True(t, delivery.VerifySignature(secret))
		assert.False(t, delivery.VerifySignature("wrong-secret"))
	})

	t.Run("SHA1", func(t *testing.T) {
		delivery := Delivery{
			Body:          body,
			SignatureSHA1: SignBodyHex(secret, body),
		}

		assert.True(t, delivery.VerifySignature(secret))
		assert.False(t, delivery.VerifySignature("wrong-secret"))
	})

	t.Run("empty secret or signature", func(t *testing.T) {
		assert.False(t, Delivery{Body: body, SignatureSHA256: "sha256=abc"}.VerifySignature(""))
		assert.False(t, Delivery{Body: body}.VerifySignature(secret))
		assert.False(t, VerifySignature(secret, body, "sha256", "bad", sha256.New))
		assert.False(t, VerifySignature(secret, body, "sha256", "sha256=zz", sha256.New))
	})
}

func TestUnitParseDelivery(t *testing.T) {
	body := []byte(`{"zen":"keep it simple"}`)
	r := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	r.Header.Set("X-GitHub-Event", "ping")
	r.Header.Set("X-GitHub-Delivery", "delivery-1")
	r.Header.Set("X-Hub-Signature-256", "sha256=abc")
	r.Header.Set("X-Hub-Signature", "sha1=def")

	delivery, err := ParseDelivery(r)
	assert.NoError(t, err)
	assert.Equal(t, "ping", delivery.Event)
	assert.Equal(t, "delivery-1", delivery.ID)
	assert.Equal(t, body, delivery.Body)
	assert.Equal(t, "sha256=abc", delivery.SignatureSHA256)
	assert.Equal(t, "sha1=def", delivery.SignatureSHA1)
}

func TestUnitSignBody(t *testing.T) {
	body := []byte("payload")
	assert.Equal(t, SignBodyHex("secret", body), "sha1="+hex.EncodeToString(SignBody("secret", body)))
}

func TestUnitDispatch(t *testing.T) {
	ctx := context.Background()

	t.Run("installation", func(t *testing.T) {
		var received Delivery
		var installation InstallationEvent
		Received.On(func(_ context.Context, d Delivery) { received = d })
		Installation.On(func(_ context.Context, e InstallationEvent) { installation = e })

		d := Delivery{
			Event: "installation",
			ID:    "1",
			Body:  []byte(`{"action":"created","installation":{"id":99}}`),
		}
		d.Dispatch(ctx)

		assert.Equal(t, "installation", received.Event)
		assert.Equal(t, "created", installation.Action)
		assert.Equal(t, int64(99), installation.Installation.ID)
	})

	t.Run("installation_repositories", func(t *testing.T) {
		var event InstallationRepositoriesEvent
		InstallationRepositories.On(func(_ context.Context, e InstallationRepositoriesEvent) { event = e })

		d := Delivery{
			Event: "installation_repositories",
			Body:  []byte(`{"action":"added","installation":{"id":7},"repositories_added":[{"id":1,"name":"repo","full_name":"o/repo"}]}`),
		}
		d.Dispatch(ctx)

		assert.Equal(t, "added", event.Action)
		assert.Equal(t, int64(7), event.Installation.ID)
		assert.Len(t, event.RepositoriesAdded, 1)
		assert.Equal(t, "repo", event.RepositoriesAdded[0].Name)
	})
}
