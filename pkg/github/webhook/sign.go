// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

// SignBody returns the SHA-1 HMAC for a webhook body.
func SignBody(secret string, body []byte) []byte {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(body)
	return mac.Sum(nil)
}

// SignBodyHex returns the legacy sha1= signature header value.
func SignBodyHex(secret string, body []byte) string {
	return "sha1=" + hex.EncodeToString(SignBody(secret, body))
}
