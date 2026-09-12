// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testApp(t *testing.T, handler http.HandlerFunc) (*App, *httptest.Server) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return &App{
		config: Config{ClientID: "Iv1.test", EncryptionKey: "test-secret"},
		key:    key,
		apiURL: server.URL,
		cache:  &memCache{items: map[string]memItem{}},
	}, server
}

type memItem struct {
	value     string
	expiresAt *time.Time
}

type memCache struct {
	items map[string]memItem
}

func (c *memCache) Get(key string) (string, *time.Time, error) {
	item, ok := c.items[key]
	if !ok {
		return "", nil, nil
	}
	return item.value, item.expiresAt, nil
}

func (c *memCache) Set(key, value string, expiresAt *time.Time) error {
	c.items[key] = memItem{value: value, expiresAt: expiresAt}
	return nil
}

func (c *memCache) DeleteExpired() (int64, error) {
	now := time.Now().UTC()
	var n int64
	for key, item := range c.items {
		if item.expiresAt != nil && !item.expiresAt.After(now) {
			delete(c.items, key)
			n++
		}
	}
	return n, nil
}

func withInstallToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/access_tokens") && r.Method == http.MethodPost {
			_ = json.NewEncoder(w).Encode(InstallationToken{Token: "ghs_test", ExpiresAt: "2099-01-01T00:00:00Z"})
			return
		}
		next(w, r)
	}
}

func TestUnitAppHTTP(t *testing.T) {
	ctx := context.Background()

	t.Run("GetInstallation", func(t *testing.T) {
		client, _ := testApp(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/app/installations/42", r.URL.Path)
			assert.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "))
			_ = json.NewEncoder(w).Encode(Installation{
				ID:      42,
				Account: InstallationAccount{Login: "acme"},
			})
		})

		inst, err := client.GetInstallation(ctx, 42)
		assert.NoError(t, err)
		assert.Equal(t, int64(42), inst.ID)
		assert.Equal(t, "acme", inst.Account.Login)
	})

	t.Run("GetInstallationToken", func(t *testing.T) {
		client, _ := testApp(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/app/installations/7/access_tokens", r.URL.Path)
			_ = json.NewEncoder(w).Encode(InstallationToken{Token: "ghs_abc", ExpiresAt: "2099-01-01T00:00:00Z"})
		})

		token, err := client.GetInstallationToken(ctx, 7)
		assert.NoError(t, err)
		assert.Equal(t, "ghs_abc", token.Token)
	})

	t.Run("ListRepositories", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/installation/repositories", r.URL.Path)
			_ = json.NewEncoder(w).Encode(RepositoriesResponse{
				TotalCount:   1,
				Repositories: []Repository{{ID: 1, Name: "ziee", FullName: "acme/ziee"}},
			})
		}))

		repos, err := client.ListRepositories(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, repos, 1)
		assert.Equal(t, "ziee", repos[0].Name)
	})

	t.Run("Issues", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/ziee/issues/3":
				_ = json.NewEncoder(w).Encode(Issue{Number: 3, Title: "Bug", State: "open"})
			case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/ziee/issues":
				_ = json.NewEncoder(w).Encode(Issue{Number: 4, Title: "Setup"})
			case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/repos/acme/ziee/issues"):
				_ = json.NewEncoder(w).Encode([]Issue{{Number: 3, Title: "Bug"}})
			case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/ziee/issues/3/comments":
				_ = json.NewEncoder(w).Encode(Comment{ID: 9, Body: "noted"})
			case r.Method == http.MethodPatch && r.URL.Path == "/repos/acme/ziee/issues/3":
				w.WriteHeader(http.StatusOK)
			case strings.HasSuffix(r.URL.Path, "/assignees"):
				w.WriteHeader(http.StatusOK)
			default:
				http.NotFound(w, r)
			}
		}))

		issue, err := client.GetIssue(ctx, 1, "acme", "ziee", 3)
		assert.NoError(t, err)
		assert.Equal(t, "Bug", issue.Title)

		issues, err := client.ListIssues(ctx, 1, "acme", "ziee")
		assert.NoError(t, err)
		assert.Len(t, issues, 1)

		comment, err := client.CreateComment(ctx, 1, "acme", "ziee", 3, "noted")
		assert.NoError(t, err)
		assert.Equal(t, "noted", comment.Body)

		assert.NoError(t, client.CloseIssue(ctx, 1, "acme", "ziee", 3))
		assert.NoError(t, client.ReopenIssue(ctx, 1, "acme", "ziee", 3))
		assert.NoError(t, client.AddAssignees(ctx, 1, "acme", "ziee", 3, []string{"octocat"}))
		assert.NoError(t, client.RemoveAssignees(ctx, 1, "acme", "ziee", 3, []string{"octocat"}))

		created, err := client.CreateIssue(ctx, 1, "acme", "ziee", "Setup", "body")
		assert.NoError(t, err)
		assert.Equal(t, 4, created.Number)
	})

	t.Run("FirstIssueAndPullRequest", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/search/issues", r.URL.Path)
			_ = json.NewEncoder(w).Encode(map[string]int{"total_count": 1})
		}))

		firstIssue, err := client.IsFirstIssue(ctx, 1, "acme", "ziee", "maya")
		assert.NoError(t, err)
		assert.True(t, firstIssue)

		firstPR, err := client.IsFirstPullRequest(ctx, 1, "acme", "ziee", "maya")
		assert.NoError(t, err)
		assert.True(t, firstPR)
	})

	t.Run("NotFirstIssueOrPullRequest", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]int{"total_count": 2})
		}))

		firstIssue, err := client.IsFirstIssue(ctx, 1, "acme", "ziee", "maya")
		assert.NoError(t, err)
		assert.False(t, firstIssue)

		firstPR, err := client.IsFirstPullRequest(ctx, 1, "acme", "ziee", "maya")
		assert.NoError(t, err)
		assert.False(t, firstPR)
	})

	t.Run("PullRequests", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/ziee/pulls/2":
				_ = json.NewEncoder(w).Encode(PullRequest{Number: 2, Title: "Feat", State: "open"})
			case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/repos/acme/ziee/pulls"):
				_ = json.NewEncoder(w).Encode([]PullRequest{{Number: 2, Title: "Feat"}})
			default:
				http.NotFound(w, r)
			}
		}))

		pr, err := client.GetPullRequest(ctx, 1, "acme", "ziee", 2)
		assert.NoError(t, err)
		assert.Equal(t, "Feat", pr.Title)

		prs, err := client.ListPullRequests(ctx, 1, "acme", "ziee")
		assert.NoError(t, err)
		assert.Len(t, prs, 1)
	})

	t.Run("FileExists", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/repos/acme/ziee/contents/.ziee.yml", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"path":".ziee.yml"}`))
		}))

		exists, err := client.FileExists(ctx, 1, "acme", "ziee", ".ziee.yml")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("GetFile", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/repos/acme/ziee/contents/.ziee.yml", r.URL.Path)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"content":  base64.StdEncoding.EncodeToString([]byte("version: 1.0.0\n")),
				"encoding": "base64",
			})
		}))

		data, err := client.GetFile(ctx, 1, "acme", "ziee", ".ziee.yml")
		assert.NoError(t, err)
		assert.Equal(t, "version: 1.0.0\n", string(data))
	})

	t.Run("FileExists missing", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))

		exists, err := client.FileExists(ctx, 1, "acme", "ziee", ".ziee.yml")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("CreatePullRequest", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/ziee":
				_ = json.NewEncoder(w).Encode(Repository{DefaultBranch: "main"})
			case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/ziee/git/ref/heads/main":
				_ = json.NewEncoder(w).Encode(map[string]any{
					"object": map[string]string{"sha": "abc123"},
				})
			case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/ziee/git/refs":
				w.WriteHeader(http.StatusCreated)
			case r.Method == http.MethodPut && r.URL.Path == "/repos/acme/ziee/contents/.ziee.yml":
				w.WriteHeader(http.StatusCreated)
			case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/ziee/pulls":
				_ = json.NewEncoder(w).Encode(PullRequest{Number: 8, Title: "Add .ziee.yml"})
			default:
				http.NotFound(w, r)
			}
		}))

		pr, err := client.CreatePullRequest(ctx, 1, "acme", "ziee", NewPullRequest{
			Branch:  "ziee/init",
			Path:    ".ziee.yml",
			Content: "version: 1.0.0\n",
			Title:   "Add .ziee.yml",
			Body:    "starter",
		})
		assert.NoError(t, err)
		assert.Equal(t, 8, pr.Number)
	})

	t.Run("Labels", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/repos/acme/ziee/labels"):
				_ = json.NewEncoder(w).Encode([]Label{{Name: "bug", Color: "f00"}})
			case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/labels"):
				w.WriteHeader(http.StatusOK)
			case r.Method == http.MethodDelete:
				w.WriteHeader(http.StatusOK)
			default:
				http.NotFound(w, r)
			}
		}))

		labels, err := client.ListLabels(ctx, 1, "acme", "ziee")
		assert.NoError(t, err)
		assert.Equal(t, "bug", labels[0].Name)

		assert.NoError(t, client.AddLabels(ctx, 1, "acme", "ziee", 3, []string{"bug"}))
		assert.NoError(t, client.RemoveLabels(ctx, 1, "acme", "ziee", 3, []string{"bug"}))
		assert.NoError(t, client.CreateLabel(ctx, 1, "acme", "ziee", Label{
			Name:        "lang/go",
			Color:       "#00ADD8",
			Description: "Go files",
		}))
	})

	t.Run("CheckRun", func(t *testing.T) {
		client, _ := testApp(t, withInstallToken(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/ziee/check-runs":
				_ = json.NewEncoder(w).Encode(CheckRun{ID: 42, Name: "Ziee configs sync", Status: "in_progress"})
			case r.Method == http.MethodPatch && r.URL.Path == "/repos/acme/ziee/check-runs/42":
				w.WriteHeader(http.StatusOK)
			default:
				http.NotFound(w, r)
			}
		}))

		run, err := client.CreateCheckRun(ctx, 1, "acme", "ziee", "Ziee configs sync", "abc123")
		assert.NoError(t, err)
		assert.Equal(t, int64(42), run.ID)

		assert.NoError(t, client.CompleteCheckRun(ctx, 1, "acme", "ziee", 42, "success", "Config synced", "Created 3 labels."))
	})
}
