// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package action

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitApply(t *testing.T) {
	assert.NoError(t, Apply(context.Background(), nil, Repo{}, Plan{}))
}
