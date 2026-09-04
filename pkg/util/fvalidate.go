// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/clivern/ziee/conf"

	"github.com/samber/lo"
)

// ChunkingOptions holds server-side chunking settings stored on a document.
type ChunkingOptions struct {
	Strategy string `json:"chunkingStrategy,omitempty"`
	Size     int    `json:"chunkSize,omitempty"`
	Overlap  int    `json:"chunkOverlap,omitempty"`
}

// DefaultChunkingOptions returns chunk settings derived from document length
func DefaultChunkingOptions(charCount int64, _filename string) ChunkingOptions {
	var size int
	switch {
	case charCount > 500_000:
		size = 10000
	case charCount > 50_000:
		size = 4000
	default:
		size = 2000
	}

	return ChunkingOptions{
		Strategy: "recursive",
		Size:     size,
		Overlap:  size / 5,
	}
}

// UploadForm is the validated result of a document multipart upload.
type UploadForm struct {
	Title       string
	Filename    string
	ContentType string
	Checksum    string
	Size        int64
	CharCount   int64
	Labels      []string
	Content     string
}

// ParseUploadForm reads and validates a multipart upload request.
func ParseUploadForm(r *http.Request) (*UploadForm, error) {
	if err := r.ParseMultipartForm(conf.MaxUploadBytes); err != nil {
		return nil, fmt.Errorf("Invalid multipart form: %w", err)
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		return nil, fmt.Errorf("File is required")
	}
	defer file.Close()

	filename := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(filename))

	if !slices.Contains([]string{".txt", ".md"}, ext) {
		return nil, fmt.Errorf("Only .txt and .md files are allowed")
	}

	bytes, err := io.ReadAll(io.LimitReader(file, conf.MaxUploadBytes+1))
	if err != nil {
		return nil, fmt.Errorf("Failed to read uploaded file")
	}
	if int64(len(bytes)) > conf.MaxUploadBytes {
		return nil, fmt.Errorf("File must not exceed 2 MB")
	}
	if !utf8.Valid(bytes) {
		return nil, fmt.Errorf("File must be valid UTF-8 text")
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if lo.IsEmpty(title) {
		return nil, fmt.Errorf("Title is required")
	}
	if len(title) > 255 {
		return nil, fmt.Errorf("Title must not exceed 255 characters")
	}

	labels, err := ParseLabels(r.FormValue("labels"))
	if err != nil {
		return nil, err
	}

	contentType := "application/octet-stream"
	switch ext {
	case ".md":
		contentType = "text/markdown"
	case ".txt":
		contentType = "text/plain"
	}

	checksum := sha256.Sum256(bytes)

	return &UploadForm{
		Title:       title,
		Filename:    filename,
		ContentType: contentType,
		Checksum:    hex.EncodeToString(checksum[:]),
		Size:        int64(len(bytes)),
		CharCount:   int64(utf8.RuneCount(bytes)),
		Labels:      labels,
		Content:     string(bytes),
	}, nil
}

// ParseLabels parses labels from a JSON array or delimited string.
func ParseLabels(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if lo.IsEmpty(raw) {
		return []string{}, nil
	}

	var parts []string
	if err := json.Unmarshal([]byte(raw), &parts); err != nil {
		parts = strings.FieldsFunc(raw, func(r rune) bool {
			return r == ',' || r == ';' || r == '\n'
		})
	}

	seen := make(map[string]struct{}, len(parts))
	labels := make([]string, 0, len(parts))
	for _, part := range parts {
		label := strings.TrimSpace(part)
		if lo.IsEmpty(label) {
			continue
		}

		key := label
		if eq := strings.Index(label, "="); eq > 0 {
			key = strings.TrimSpace(label[:eq])
		}
		if strings.HasPrefix(key, "_") {
			return nil, fmt.Errorf("Labels must not use keys starting with _")
		}

		if _, ok := seen[label]; ok {
			continue
		}
		seen[label] = struct{}{}
		labels = append(labels, label)
	}

	return labels, nil
}
