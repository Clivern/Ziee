# pkg/github

GitHub OAuth and webhook verification for Ziee.

## Layout

| Path | Purpose |
|------|---------|
| `pkg/github` | OAuth helpers |
| `pkg/github/webhook` | Webhook parse and signature verification |

## OAuth

```go
oauth := github.NewOAuth(github.OAuthConfig{
    ClientID:     clientID,
    ClientSecret: clientSecret,
    RedirectURL:  redirectURL,
    Scopes:       []string{"read:user", "user:email"},
})

url := oauth.AuthorizeURL(state)
token, err := oauth.Exchange(ctx, code, state, expectedState)
user, err := oauth.User(ctx, token.AccessToken)
emails, err := oauth.Emails(ctx, token.AccessToken)
```

## Webhook verification

Supports both `X-Hub-Signature-256` and legacy `X-Hub-Signature` headers. Use from any framework route:

```go
import "github.com/actx0/ziee/pkg/github/webhook"

delivery, err := webhook.ParseDelivery(r)
if err != nil {
    // 400
}
if !delivery.VerifySignature(secret) {
    // 401
}
// delivery.Event, delivery.ID, delivery.Body
```
