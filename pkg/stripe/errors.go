// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package stripe

import (
	"errors"
)

var (
	ErrBillingDisabled      = errors.New("stripe billing is not configured")
	ErrInvalidPlan          = errors.New("invalid billing plan")
	ErrWebhookNotConfigured = errors.New("stripe webhook secret is not configured")
)
