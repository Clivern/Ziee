// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeliveryVerifySignatureSHA256(t *testing.T) {
	secret := "test-secret"
	body := []byte(`{"action":"opened"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	delivery := Delivery{
		Body:            body,
		SignatureSHA256: signature,
	}

	assert.True(t, delivery.VerifySignature(secret))
	assert.False(t, delivery.VerifySignature("wrong-secret"))
}

func TestDeliveryVerifySignatureSHA1(t *testing.T) {
	secret := "test-secret"
	body := []byte(`{"action":"opened"}`)

	delivery := Delivery{
		Body:          body,
		SignatureSHA1: SignBodyHex(secret, body),
	}

	assert.True(t, delivery.VerifySignature(secret))
	assert.False(t, delivery.VerifySignature("wrong-secret"))
}
