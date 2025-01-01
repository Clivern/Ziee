// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package chunk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitNewAndSplit(t *testing.T) {
	t.Run("defaults empty splitter to recursive", func(t *testing.T) {
		c, err := New(Config{Size: 50, Overlap: 0})
		require.NoError(t, err)
		assert.Equal(t, Config{Size: 50, Overlap: 0}, c.Config())
	})

	t.Run("unsupported splitter", func(t *testing.T) {
		_, err := New(Config{Splitter: "nope", Size: 50})
		assert.Error(t, err)
	})

	t.Run("recursive split", func(t *testing.T) {
		c, err := New(Config{
			Splitter: SplitterRecursive,
			Size:     40,
			Overlap:  0,
			Recursive: RecursiveConfig{
				Separators: []string{"\n"},
			},
		})
		require.NoError(t, err)

		chunks, err := c.Split("alpha\nbeta\ngamma\ndelta\nepsilon")
		require.NoError(t, err)
		assert.NotEmpty(t, chunks)
		assert.Equal(t, "alpha\nbeta\ngamma\ndelta\nepsilon", strings.Join(chunks, "\n"))
	})

	t.Run("markdown split", func(t *testing.T) {
		c, err := New(Config{
			Splitter: SplitterMarkdown,
			Size:     80,
			Overlap:  0,
		})
		require.NoError(t, err)

		chunks, err := c.Split("# Title\n\nFirst paragraph.\n\n## Section\n\nSecond paragraph.")
		require.NoError(t, err)
		assert.NotEmpty(t, chunks)
	})
}
