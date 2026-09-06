// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package spec

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	v1 "github.com/clivern/ziee/policy/spec/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitParseRoutesV1(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	data, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "..", ".ziee.yml"))
	require.NoError(t, err)

	file, err := Parse(data)
	require.NoError(t, err)
	assert.Equal(t, v1.Version, file.Version)
}

func TestUnitParseUnsupportedVersion(t *testing.T) {
	_, err := Parse([]byte("version: 2.0.0\n"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported version")

	_, err = Parse([]byte("version: 1.1.0\n"))
	require.Error(t, err)
}
