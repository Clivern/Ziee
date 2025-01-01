// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
}
