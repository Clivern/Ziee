// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package hmac

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	secret := "whsec_testsecret"
	body := []byte(`{"type":"memory.add"}`)

	signature := Sign(secret, body)

	assert.Equal(t, "sha256=71075f1b564f83edecda3b3be0d9faf890cd934a43dbf597579c852a33a9acde", signature)
	assert.True(t, Verify(secret, body, signature))
	assert.False(t, Verify("whsec_other", body, signature))
	assert.False(t, Verify(secret, []byte(`{"type":"memory.delete"}`), signature))
	assert.False(t, Verify(secret, body, "sha1=71075f1b564f83edecda3b3be0d9faf890cd934a43dbf597579c852a33a9acde"))
	assert.False(t, Verify(secret, body, "sha256=zz"))
}

func TestSetHeader(t *testing.T) {
	header := make(http.Header)
	body := []byte(`{"type":"memory.add"}`)

	SetHeader(header, "whsec_testsecret", body)

	assert.Equal(t, Sign("whsec_testsecret", body), header.Get(Header))
}
