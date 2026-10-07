// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"unicode/utf8"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ChangeType is how a file differs from HEAD.
type ChangeType string

const (
	// ChangeCreated is a file that is not in HEAD.
	ChangeCreated ChangeType = "created"
	// ChangeUpdated is a file in HEAD whose content or mode changed.
	ChangeUpdated ChangeType = "updated"
	// ChangeDeleted is a file in HEAD that is gone from the working tree.
	ChangeDeleted ChangeType = "deleted"
)

// Change is one file that differs between HEAD and the working tree.
type Change struct {
	Path string     `json:"path"`
	Type ChangeType `json:"type"`
	// Content is the new file content. It is empty for deleted and binary files.
	// For a symlink it is the link target.
	Content string `json:"content,omitempty"`
	Binary  bool   `json:"binary,omitempty"`
	Symlink bool   `json:"symlink,omitempty"`
}

// Git reads the changes made to a sandbox repository clone.
type Git struct {
	dir string
}

// NewGit returns a Git for the repository at dir.
func NewGit(dir string) *Git {
	return &Git{dir: dir}
}

// Changes lists every file that differs from HEAD, including staged and
// untracked files but not ignored ones, sorted by path.
func (g *Git) Changes() ([]Change, error) {
	repo, err := git.PlainOpen(g.dir)
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	changes := []Change{}
	for path, file := range status {
		if file.Staging == git.Unmodified && file.Worktree == git.Unmodified {
			continue
		}

		change, ok, err := g.change(tree, path)
		if err != nil {
			return nil, err
		}
		if ok {
			changes = append(changes, change)
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})

	return changes, nil
}

// change compares one path between HEAD and the working tree.
// It reports false when the path is in neither, such as a file added then removed.
func (g *Git) change(tree *object.Tree, path string) (Change, bool, error) {
	change := Change{Path: path}

	_, err := tree.FindEntry(path)
	inHead := err == nil
	if err != nil && !errors.Is(err, object.ErrEntryNotFound) && !errors.Is(err, object.ErrDirectoryNotFound) {
		return change, false, err
	}

	full := filepath.Join(g.dir, filepath.FromSlash(path))

	// Lstat so a symlink is never followed out of the repository.
	info, err := os.Lstat(full)
	onDisk := err == nil && !info.IsDir()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return change, false, err
	}

	switch {
	case inHead && !onDisk:
		change.Type = ChangeDeleted
		return change, true, nil
	case !inHead && !onDisk:
		return change, false, nil
	case inHead:
		change.Type = ChangeUpdated
	default:
		change.Type = ChangeCreated
	}

	if info.Mode()&fs.ModeSymlink != 0 {
		target, err := os.Readlink(full)
		if err != nil {
			return change, false, err
		}
		change.Symlink = true
		change.Content = target
		return change, true, nil
	}

	data, err := os.ReadFile(full)
	if err != nil {
		return change, false, err
	}

	if bytes.IndexByte(data, 0) >= 0 || !utf8.Valid(data) {
		change.Binary = true
		return change, true, nil
	}

	change.Content = string(data)
	return change, true, nil
}
