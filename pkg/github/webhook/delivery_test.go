// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func TestUnitInstallationEvent(t *testing.T) {
	var installation InstallationEvent
	json.Unmarshal([]byte(`{"action":"created","installation":{"id":99}}`), &installation)

	assert.Equal(t, "created", installation.Action)
	assert.Equal(t, int64(99), installation.Installation.ID)
}

func TestUnitInstallationRepositoriesEvent(t *testing.T) {
	var event InstallationRepositoriesEvent
	json.Unmarshal([]byte(`{"action":"added","installation":{"id":7},"repositories_added":[{"id":1,"name":"repo","full_name":"o/repo"}]}`), &event)

	assert.Equal(t, "added", event.Action)
	assert.Equal(t, int64(7), event.Installation.ID)
	assert.Len(t, event.RepositoriesAdded, 1)
	assert.Equal(t, "repo", event.RepositoriesAdded[0].Name)
}

func TestUnitIssueEvent(t *testing.T) {
	var event IssueEvent
	json.Unmarshal([]byte(`{"action":"opened","issue":{"number":3,"title":"bug","user":{"login":"maya"}},"repository":{"name":"ziee","owner":{"login":"acme","type":"Organization"}},"installation":{"id":9}}`), &event)

	assert.Equal(t, "opened", event.Action)
	assert.Equal(t, 3, event.Issue.Number)
	assert.Equal(t, "maya", event.Issue.User.Login)
	assert.Equal(t, "acme", event.Repository.Owner.Login)
	assert.Equal(t, "Organization", event.Repository.Owner.Type)
	assert.Equal(t, int64(9), event.Installation.ID)
}

func TestUnitIssueCommentEvent(t *testing.T) {
	var event IssueCommentEvent
	json.Unmarshal([]byte(`{"action":"created","issue":{"number":3},"comment":{"body":"@ziee label bug","user":{"login":"maya"}},"repository":{"name":"ziee","owner":{"login":"acme"}},"installation":{"id":9}}`), &event)

	assert.Equal(t, "created", event.Action)
	assert.Equal(t, 3, event.Issue.Number)
	assert.Equal(t, "@ziee label bug", event.Comment.Body)
	assert.Equal(t, int64(9), event.Installation.ID)
}

func TestUnitPushEventChanged(t *testing.T) {
	var event PushEvent
	json.Unmarshal([]byte(`{
		"ref":"refs/heads/main",
		"repository":{"name":"ziee","full_name":"acme/ziee","owner":{"login":"acme"},"default_branch":"main"},
		"installation":{"id":9},
		"commits":[{
			"added":["worker/labels.go"],
			"removed":[],
			"modified":["conf/conf.go",".ziee.yml"]
		}]
	}`), &event)

	assert.Equal(t, "refs/heads/main", event.Ref)
	assert.Equal(t, "main", event.Repository.DefaultBranch)
	assert.Equal(t, int64(9), event.Installation.ID)
	assert.True(t, event.Changed(".ziee.yml"))
	assert.True(t, event.Changed(".ziee.yml", ".ziee.yaml"))
	assert.True(t, event.Changed("worker/labels.go", "pkg/missing.go"))
	assert.False(t, event.Changed("README.md"))
}
