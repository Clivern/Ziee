// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientNewComment(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/repos/actx0/ziee/issues/42/comments", r.URL.Path)
		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))

		var body map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "LGTM", body["body"])

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":1,"body":"LGTM"}`))
	}))
	t.Cleanup(server.Close)

	client := New(Config{
		Token:      "token",
		Owner:      "actx0",
		Repository: "ziee",
	}, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	comment, err := client.NewComment(context.Background(), 42, "LGTM")
	require.NoError(t, err)
	assert.Equal(t, 1, comment.ID)
	assert.Equal(t, "LGTM", comment.Body)
}

func TestClientDeleteLabel(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/repos/actx0/ziee/labels/bug", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	client := New(Config{
		Token:      "token",
		Owner:      "actx0",
		Repository: "ziee",
	}, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	err := client.DeleteLabel(context.Background(), "bug")
	require.NoError(t, err)
}
