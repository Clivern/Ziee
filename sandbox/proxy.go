// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/clivern/ziee/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	ChatCompletionsPath = "/api/v1/chat/completions"
	MessagesPath        = "/api/v1/messages"
	ResponsesPath       = "/api/v1/responses"
	OpenRouterAPIHost   = "openrouter.ai"
)

// Proxy forwards sandbox chat requests to OpenRouter.
type Proxy struct {
	token string
}

// NewProxy returns a proxy that authenticates to OpenRouter with the sandbox API key.
func NewProxy() *Proxy {
	return &Proxy{
		token: viper.GetString("app.sandbox.api_key"),
	}
}

// DedupePath collapses consecutive duplicate path segments.
// Query strings are preserved.
//
//	/api/v1/v1/messages?beta=true     → /api/v1/messages?beta=true
//	/api/api/v1/v1/messages?beta=true → /api/v1/messages?beta=true
func DedupePath(raw string) string {
	path, query, found := strings.Cut(raw, "?")
	if found {
		query = "?" + query
	}

	parts := strings.Split(path, "/")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(out) > 0 && out[len(out)-1] == part && part != "" {
			continue
		}
		out = append(out, part)
	}

	return strings.Join(out, "/") + query
}

// IsChatEndpoint reports whether raw is an OpenRouter chat endpoint path.
func IsChatEndpoint(raw string) bool {
	path, _, _ := strings.Cut(DedupePath(raw), "?")
	path = strings.TrimRight(path, "/")

	switch path {
	case ChatCompletionsPath, MessagesPath, ResponsesPath:
		return true
	default:
		return false
	}
}

// ServeHTTP proxies /api/v1/sandbox/{Id}/{openrouter_path} to OpenRouter.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "Id")

	if IsBlocked(id) {
		log.Info().
			Str("id", id).
			Msg("Blocked sandbox id")

		util.WriteJSON(w, http.StatusForbidden, map[string]any{
			"errorMessage": "access forbidden",
		})
		return
	}

	rest := chi.URLParam(r, "*")
	openPath := DedupePath("/" + strings.TrimPrefix(rest, "/"))
	if !IsChatEndpoint(openPath) {
		util.WriteJSON(w, http.StatusForbidden, map[string]any{
			"errorMessage": "access forbidden",
		})
		return
	}

	pathOnly, _, _ := strings.Cut(openPath, "?")
	token := p.token
	rawQuery := r.URL.RawQuery

	proxy := &httputil.ReverseProxy{
		FlushInterval: -1,
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = OpenRouterAPIHost
			req.URL.Path = pathOnly
			req.URL.RawPath = ""
			req.URL.RawQuery = rawQuery
			req.Host = OpenRouterAPIHost
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			req.Header.Del("X-Api-Key")
			// Let Transport decompress gzip so ModifyResponse sees plain JSON.
			req.Header.Del("Accept-Encoding")
		},
		ModifyResponse: func(resp *http.Response) error {
			return CaptureUsage(resp, id, pathOnly)
		},
	}

	proxy.ServeHTTP(w, r)
}
