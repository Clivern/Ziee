// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"
)

const jwtHeader = `{"alg":"RS256","typ":"JWT"}`

// ParsePrivateKey parses a PKCS#1 or PKCS#8 RSA private key PEM.
func ParsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, nil
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("github app parse key: %w", err)
	}

	return parsed.(*rsa.PrivateKey), nil
}

// JWT returns a GitHub App JWT signed with the private key.
func (a *App) JWT() (string, error) {
	now := time.Now()
	claims, err := json.Marshal(map[string]any{
		"iat": now.Add(-60 * time.Second).Unix(),
		"exp": now.Add(9 * time.Minute).Unix(),
		"iss": a.config.ClientID,
	})
	if err != nil {
		return "", fmt.Errorf("github app jwt claims: %w", err)
	}

	unsigned := base64.RawURLEncoding.EncodeToString([]byte(jwtHeader)) + "." +
		base64.RawURLEncoding.EncodeToString(claims)

	hashed := sha256.Sum256([]byte(unsigned))
	sig, err := rsa.SignPKCS1v15(rand.Reader, a.key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("github app jwt sign: %w", err)
	}

	return unsigned + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}
