// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUnitPasswordHashing(t *testing.T) {
	t.Run("HashPassword", func(t *testing.T) {
		passwords := []string{
			"mySecurePassword123",
			"",
			"thisIsAVeryLongPasswordWithCharacters1234567890!@#$%^&*()_+-=[]",
			"p@ssw0rd!#$%^&*()",
		}

		for _, password := range passwords {
			t.Run(password, func(t *testing.T) {
				hashed, err := HashPassword(password)
				assert.NoError(t, err)
				assert.NotEmpty(t, hashed)
				assert.NotEqual(t, password, hashed)
			})
		}

		hash1, err1 := HashPassword("samePassword")
		hash2, err2 := HashPassword("samePassword")
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("ComparePassword", func(t *testing.T) {
		cases := []struct {
			name     string
			password string
			compare  string
			hash     string
			expected bool
		}{
			{
				name:     "matching password",
				password: "mySecurePassword123",
				compare:  "mySecurePassword123",
				expected: true,
			},
			{
				name:     "non-matching password",
				password: "correctPassword",
				compare:  "wrongPassword",
				expected: false,
			},
			{
				name:     "empty password",
				password: "somePassword",
				compare:  "",
				expected: false,
			},
			{
				name:     "case sensitivity",
				password: "Password123",
				compare:  "password123",
				expected: false,
			},
			{
				name:     "special characters",
				password: "p@ssw0rd!#$%",
				compare:  "p@ssw0rd!#$%",
				expected: true,
			},
			{
				name:     "invalid hash",
				hash:     "notAValidBcryptHash",
				compare:  "somePassword",
				expected: false,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				hash := tc.hash
				if lo.IsEmpty(hash) {
					var err error
					hash, err = HashPassword(tc.password)
					assert.NoError(t, err)
				}
				assert.Equal(t, tc.expected, ComparePassword(hash, tc.compare))
			})
		}
	})

	t.Run("Complete password workflow", func(t *testing.T) {
		passwords := []string{
			"simplePassword",
			"Complex123!@#",
			"",
			"with spaces in password",
			"unicode密码🔐",
		}

		for _, password := range passwords {
			t.Run("Password: "+password, func(t *testing.T) {
				hashed, err := HashPassword(password)
				assert.NoError(t, err)
				assert.True(t, ComparePassword(hashed, password))
				assert.False(t, ComparePassword(hashed, password+"wrong"))
			})
		}
	})

	t.Run("RandomHash", func(t *testing.T) {
		hash, err := RandomHash()
		assert.NoError(t, err)
		assert.Len(t, hash, 12)

		hash2, err := RandomHash()
		assert.NoError(t, err)
		assert.NotEqual(t, hash, hash2)
	})
}
