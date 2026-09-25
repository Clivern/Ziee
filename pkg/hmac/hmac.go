// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package hmac

import (
	cryptohmac "crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

const (
	// Header is the signature header on outgoing webhook deliveries.
	Header = "X-Ziee-Signature-256"
	prefix = "sha256="
)

// Sign returns the HMAC-SHA256 signature for an outgoing webhook body.
func Sign(secret string, body []byte) string {
	mac := cryptohmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	return prefix + hex.EncodeToString(mac.Sum(nil))
}

// SetHeader signs body and sets Header on an outgoing request.
func SetHeader(header http.Header, secret string, body []byte) {
	header.Set(Header, Sign(secret, body))
}

// Verify reports whether signature matches the HMAC-SHA256 of body.
func Verify(secret string, body []byte, signature string) bool {
	if !strings.HasPrefix(signature, prefix) {
		return false
	}

	expected, err := hex.DecodeString(signature[len(prefix):])
	if err != nil {
		return false
	}

	mac := cryptohmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	return cryptohmac.Equal(mac.Sum(nil), expected)
}
