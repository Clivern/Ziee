# pkg/github

GitHub webhook handling and REST helpers for Ziee. Migrated and polished from [Hamster](https://github.com/Clivern/Hamster).

## Layout

| Path | Purpose |
|------|---------|
| `pkg/github` | REST client, OAuth, shared errors |
| `pkg/github/webhook` | Webhook verification, dispatch, and HTTP handler |
| `pkg/github/event` | Typed GitHub webhook payloads |
| `pkg/github/sender` | Outbound request bodies |
| `pkg/github/response` | API response models |

## REST client

All API methods accept a `context.Context` and return structured `*github.APIError` values on non-success responses.

```go
client := github.New(github.Config{
    Token:      os.Getenv("GITHUB_TOKEN"),
    Owner:      "ziee",
    Repository: "ziee",
})

comment, err := client.NewComment(ctx, 42, "LGTM")
```

## OAuth

```go
oauth := github.NewOAuth(github.OAuthConfig{
    ClientID:     clientID,
    ClientSecret: clientSecret,
    RedirectURL:  redirectURL,
    Scopes:       []string{"repo"},
}, nil)

url, err := oauth.AuthorizeURL(state)
token, err := oauth.Exchange(ctx, code, state, expectedState)
```

## Webhook handler

Supports both `X-Hub-Signature-256` and legacy `X-Hub-Signature` headers.

```go
import (
    "github.com/actx0/ziee/pkg/github/event"
    "github.com/actx0/ziee/pkg/github/webhook"
)

handler := &webhook.Handler{
    Secret: os.Getenv("GITHUB_WEBHOOK_SECRET"),
    Hooks: webhook.Hooks{
        PullRequest: func(pr event.PullRequest) error {
            return nil
        },
    },
}

http.Handle("/webhooks/github", handler)
```

## Fixtures

Sample payloads live in `pkg/github/fixtures/` for tests and local development.
