// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package util

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash from a plain text password.
// It uses the default cost factor (bcrypt.DefaultCost = 10).
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a bcrypt hashed password with a plain text password.
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(password),
	)
	return err == nil
}

// RandomHash returns a 12-character cryptographically secure random hex string.
func RandomHash() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
