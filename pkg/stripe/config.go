// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package stripe

const (
	PlanHobby   = "hobby"
	PlanPro     = "pro"
	PlanGrowth  = "growth"
	PlanStarter = "starter"
)

// Prices maps actx0 plan names to Stripe Price IDs.
type Prices struct {
	Hobby   string
	Pro     string
	Growth  string
	Starter string
}

// URLs holds Stripe redirect URLs (relative or absolute).
type URLs struct {
	CheckoutSuccess string
	CheckoutCancel  string
	PortalReturn    string
}

// Config holds Stripe client configuration.
type Config struct {
	SecretKey     string
	WebhookSecret string
	ProductId     string
	Prices        Prices
	URLs          URLs
}

// PriceId returns the Stripe Price Id for a plan name.
func (c Config) PriceId(plan string) string {
	switch plan {
	case PlanHobby:
		return c.Prices.Hobby
	case PlanPro:
		return c.Prices.Pro
	case PlanGrowth:
		return c.Prices.Growth
	case PlanStarter:
		return c.Prices.Starter
	default:
		return ""
	}
}

// PlanForPriceId returns the actx0 plan name for a Stripe Price Id.
func (c Config) PlanForPriceId(priceId string) string {
	switch priceId {
	case c.Prices.Hobby:
		return PlanHobby
	case c.Prices.Pro:
		return PlanPro
	case c.Prices.Growth:
		return PlanGrowth
	case c.Prices.Starter:
		return PlanStarter
	default:
		return ""
	}
}
