// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type OpenRouterUsage struct {
	PromptTokens     int64    `json:"prompt_tokens"`
	CompletionTokens int64    `json:"completion_tokens"`
	TotalTokens      int64    `json:"total_tokens"`
	InputTokens      int64    `json:"input_tokens"`
	OutputTokens     int64    `json:"output_tokens"`
	Cost             *float64 `json:"cost"`
}

type UsageEnvelope struct {
	Usage *OpenRouterUsage `json:"usage"`
}

type SSEUsageReader struct {
	io.ReadCloser
	Id      string
	Path    string
	pending []byte
}

func CaptureUsage(resp *http.Response, id, path string) error {
	if resp.StatusCode >= http.StatusBadRequest {
		return nil
	}

	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "text/event-stream") {
		resp.Body = &SSEUsageReader{
			ReadCloser: resp.Body,
			Id:         id,
			Path:       path,
		}
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		resp.Body = io.NopCloser(bytes.NewReader(nil))
		return err
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(bytes.NewReader(body))
		if err == nil {
			plain, err := io.ReadAll(gr)
			gr.Close()
			if err == nil {
				body = plain
				resp.Header.Del("Content-Encoding")
			}
		}
	}

	LogUsage(body, id, path)

	resp.Body = io.NopCloser(bytes.NewReader(body))
	resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))
	resp.ContentLength = int64(len(body))

	return nil
}

func LogUsage(raw []byte, id, path string) bool {
	var env UsageEnvelope
	err := json.Unmarshal(raw, &env)
	if err != nil || env.Usage == nil {
		return false
	}

	prompt, total, cost := env.Usage.Values()
	if prompt == 0 && total == 0 && cost == 0 {
		return false
	}

	log.Info().
		Str("id", id).
		Str("path", path).
		Int64("promptTokens", prompt).
		Int64("totalTokens", total).
		Float64("costUsd", cost).
		Msg("Sandbox OpenRouter usage")

	return true
}

func (u *OpenRouterUsage) Values() (prompt, total int64, cost float64) {
	prompt = u.PromptTokens
	if prompt == 0 {
		prompt = u.InputTokens
	}

	total = u.TotalTokens
	if total == 0 {
		total = u.PromptTokens + u.CompletionTokens
	}
	if total == 0 {
		total = u.InputTokens + u.OutputTokens
	}

	if u.Cost != nil {
		cost = *u.Cost
	}

	return prompt, total, cost
}

func (s *SSEUsageReader) Read(p []byte) (int, error) {
	n, err := s.ReadCloser.Read(p)
	if n > 0 {
		s.pending = append(s.pending, p[:n]...)
		s.drain()
	}

	return n, err
}

func (s *SSEUsageReader) drain() {
	for {
		i := bytes.IndexByte(s.pending, '\n')
		if i < 0 {
			return
		}

		line := bytes.TrimSpace(s.pending[:i])
		s.pending = s.pending[i+1:]
		if !bytes.HasPrefix(line, []byte("data:")) {
			continue
		}

		payload := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data:")))
		if bytes.Equal(payload, []byte("[DONE]")) {
			continue
		}

		LogUsage(payload, s.Id, s.Path)
	}
}
