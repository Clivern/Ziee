// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package util

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitFValidate(t *testing.T) {
	t.Run("DefaultChunkingOptions", func(t *testing.T) {
		small := DefaultChunkingOptions(1_000, "notes.md")
		assert.Equal(t, "recursive", small.Strategy)
		assert.Equal(t, 2000, small.Size)
		assert.Equal(t, 400, small.Overlap)

		medium := DefaultChunkingOptions(100_000, "notes.md")
		assert.Equal(t, 4000, medium.Size)
		assert.Equal(t, 800, medium.Overlap)

		large := DefaultChunkingOptions(600_000, "notes.md")
		assert.Equal(t, 10000, large.Size)
		assert.Equal(t, 2000, large.Overlap)
	})

	t.Run("ParseLabels", func(t *testing.T) {
		empty, err := ParseLabels("")
		assert.NoError(t, err)
		assert.Empty(t, empty)

		jsonLabels, err := ParseLabels(`["env=prod","team"]`)
		assert.NoError(t, err)
		assert.Equal(t, []string{"env=prod", "team"}, jsonLabels)

		delimited, err := ParseLabels("env=prod, team;ops")
		assert.NoError(t, err)
		assert.Equal(t, []string{"env=prod", "team", "ops"}, delimited)

		_, err = ParseLabels("_secret")
		assert.Error(t, err)
	})

	t.Run("ParseUploadForm", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		part, err := writer.CreateFormFile("file", "readme.md")
		assert.NoError(t, err)
		_, err = part.Write([]byte("# hello"))
		assert.NoError(t, err)
		assert.NoError(t, writer.WriteField("title", "Readme"))
		assert.NoError(t, writer.WriteField("labels", `["docs"]`))
		assert.NoError(t, writer.Close())

		r := httptest.NewRequest(http.MethodPost, "/", &body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		form, err := ParseUploadForm(r)
		assert.NoError(t, err)
		assert.Equal(t, "Readme", form.Title)
		assert.Equal(t, "readme.md", form.Filename)
		assert.Equal(t, "text/markdown", form.ContentType)
		assert.Equal(t, "# hello", form.Content)
		assert.Equal(t, []string{"docs"}, form.Labels)
		assert.NotEmpty(t, form.Checksum)
		assert.Equal(t, int64(7), form.Size)

		bad := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-multipart"))
		_, err = ParseUploadForm(bad)
		assert.Error(t, err)
	})
}
