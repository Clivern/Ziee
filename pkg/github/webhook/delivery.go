// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
	"strings"
)

// Delivery is a parsed GitHub webhook request.
type Delivery struct {
	Event           string
	ID              string
	Body            []byte
	SignatureSHA1   string
	SignatureSHA256 string
}

// ParseDelivery reads and parses a GitHub webhook HTTP request.
func ParseDelivery(r *http.Request) (Delivery, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return Delivery{}, err
	}

	return Delivery{
		Event:           r.Header.Get("X-GitHub-Event"),
		ID:              r.Header.Get("X-GitHub-Delivery"),
		Body:            body,
		SignatureSHA1:   r.Header.Get("X-Hub-Signature"),
		SignatureSHA256: r.Header.Get("X-Hub-Signature-256"),
	}, nil
}

// VerifySignature validates the webhook HMAC using SHA-256 or SHA-1.
func (d Delivery) VerifySignature(secret string) bool {
	if secret == "" {
		return false
	}

	if d.SignatureSHA256 != "" {
		return VerifySignature(secret, d.Body, "sha256", d.SignatureSHA256, sha256.New)
	}

	if d.SignatureSHA1 != "" {
		return VerifySignature(secret, d.Body, "sha1", d.SignatureSHA1, sha1.New)
	}

	return false
}

// VerifySignature validates the webhook HMAC using SHA-256 or SHA-1.
func VerifySignature(secret string, body []byte, prefix, signature string, newHash func() hash.Hash) bool {
	expectedPrefix := prefix + "="
	if !strings.HasPrefix(signature, expectedPrefix) {
		return false
	}

	expected, err := hex.DecodeString(signature[len(expectedPrefix):])
	if err != nil {
		return false
	}

	mac := hmac.New(newHash, []byte(secret))
	mac.Write(body)

	return hmac.Equal(mac.Sum(nil), expected)
}
