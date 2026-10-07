// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, dir, path, content string) {
	t.Helper()

	full := filepath.Join(dir, path)
	require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
	require.NoError(t, os.WriteFile(full, []byte(content), 0o644))
}

// initRepo creates a repository in a temp dir with files committed to HEAD.
func initRepo(t *testing.T, files map[string]string) (string, *git.Worktree) {
	t.Helper()

	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)
	worktree, err := repo.Worktree()
	require.NoError(t, err)

	for path, content := range files {
		writeFile(t, dir, path, content)
	}
	require.NoError(t, worktree.AddGlob("."))
	_, err = worktree.Commit("init", &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
	})
	require.NoError(t, err)

	return dir, worktree
}

func changes(t *testing.T, dir string) []Change {
	t.Helper()

	changes, err := NewGit(dir).Changes()
	require.NoError(t, err)
	return changes
}

func TestUnitSandboxGit(t *testing.T) {
	t.Run("No changes", func(t *testing.T) {
		dir, _ := initRepo(t, map[string]string{"a.txt": "a\n"})
		assert.Empty(t, changes(t, dir))
	})

	t.Run("Updated file", func(t *testing.T) {
		dir, _ := initRepo(t, map[string]string{"main.go": "package main\n"})
		writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")

		assert.Equal(t, []Change{
			{Path: "main.go", Type: ChangeUpdated, Content: "package main\n\nfunc main() {}\n"},
		}, changes(t, dir))
	})

	t.Run("Deleted file", func(t *testing.T) {
		dir, _ := initRepo(t, map[string]string{"pkg/old.go": "package pkg\n", "keep.txt": "k\n"})
		require.NoError(t, os.Remove(filepath.Join(dir, "pkg/old.go")))

		assert.Equal(t, []Change{
			{Path: "pkg/old.go", Type: ChangeDeleted},
		}, changes(t, dir))
	})

	t.Run("Created file in a new directory", func(t *testing.T) {
		dir, _ := initRepo(t, map[string]string{"a.txt": "a\n"})
		writeFile(t, dir, "pkg/new/new.go", "package new\n")

		assert.Equal(t, []Change{
			{Path: "pkg/new/new.go", Type: ChangeCreated, Content: "package new\n"},
		}, changes(t, dir))
	})

	t.Run("Ignored file is skipped", func(t *testing.T) {
		dir, _ := initRepo(t, map[string]string{".gitignore": "*.log\n"})
		writeFile(t, dir, "debug.log", "ignored\n")

		assert.Empty(t, changes(t, dir))
	})

	t.Run("Binary file has no content", func(t *testing.T) {
		dir, _ := initRepo(t, map[string]string{"a.txt": "a\n"})
		writeFile(t, dir, "image.bin", "\x00\x01\x02")

		assert.Equal(t, []Change{
			{Path: "image.bin", Type: ChangeCreated, Binary: true},
		}, changes(t, dir))
	})

	t.Run("Symlink is not followed", func(t *testing.T) {
		dir, _ := initRepo(t, map[string]string{"a.txt": "a\n"})
		require.NoError(t, os.Symlink("/etc/passwd", filepath.Join(dir, "link")))

		assert.Equal(t, []Change{
			{Path: "link", Type: ChangeCreated, Content: "/etc/passwd", Symlink: true},
		}, changes(t, dir))
	})

	t.Run("Staged change", func(t *testing.T) {
		dir, worktree := initRepo(t, map[string]string{"staged.txt": "old\n"})
		writeFile(t, dir, "staged.txt", "new\n")
		_, err := worktree.Add("staged.txt")
		require.NoError(t, err)

		assert.Equal(t, []Change{
			{Path: "staged.txt", Type: ChangeUpdated, Content: "new\n"},
		}, changes(t, dir))
	})

	t.Run("Staged then removed from disk is skipped", func(t *testing.T) {
		dir, worktree := initRepo(t, map[string]string{"a.txt": "a\n"})
		writeFile(t, dir, "gone.txt", "x\n")
		_, err := worktree.Add("gone.txt")
		require.NoError(t, err)
		require.NoError(t, os.Remove(filepath.Join(dir, "gone.txt")))

		assert.Empty(t, changes(t, dir))
	})

	t.Run("Changes are sorted by path", func(t *testing.T) {
		dir, _ := initRepo(t, map[string]string{"b.txt": "b\n", "c.txt": "c\n"})
		writeFile(t, dir, "c.txt", "c2\n")
		writeFile(t, dir, "a.txt", "a\n")
		require.NoError(t, os.Remove(filepath.Join(dir, "b.txt")))

		assert.Equal(t, []Change{
			{Path: "a.txt", Type: ChangeCreated, Content: "a\n"},
			{Path: "b.txt", Type: ChangeDeleted},
			{Path: "c.txt", Type: ChangeUpdated, Content: "c2\n"},
		}, changes(t, dir))
	})

	t.Run("Not a repository", func(t *testing.T) {
		_, err := NewGit(t.TempDir()).Changes()
		assert.Error(t, err)
	})
}
