// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package conf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnitConf(t *testing.T) {
	t.Run("constants", func(t *testing.T) {
		assert.Equal(t, time.Second, SlowRequestThreshold)
		assert.Equal(t, 7*24*time.Hour, DefaultSessionDuration)
		assert.Equal(t, 30*24*time.Hour, RememberMeSessionDuration)
		assert.Equal(t, 7*24*time.Hour, InviteExpiry)
		assert.Equal(t, int64(2*1024*1024), int64(MaxUploadBytes))
		assert.Equal(t, 10, DefaultSearchLimit)
	})

	t.Run("Edition", func(t *testing.T) {
		prev := BuiltEdition
		t.Cleanup(func() { BuiltEdition = prev })

		t.Setenv("ZIEE_EDITION", "")
		BuiltEdition = EditionOSS
		assert.Equal(t, EditionOSS, Edition())
		assert.False(t, IsSaaS())

		BuiltEdition = EditionSaaS
		assert.Equal(t, EditionSaaS, Edition())
		assert.True(t, IsSaaS())

		t.Setenv("ZIEE_EDITION", EditionOSS)
		assert.Equal(t, EditionOSS, Edition())
		assert.False(t, IsSaaS())

		t.Setenv("ZIEE_EDITION", EditionSaaS)
		assert.Equal(t, EditionSaaS, Edition())
		assert.True(t, IsSaaS())

		t.Setenv("ZIEE_EDITION", "unknown")
		BuiltEdition = EditionOSS
		assert.Equal(t, EditionOSS, Edition())
	})
}
