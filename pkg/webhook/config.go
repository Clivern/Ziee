// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package webhook

// Config holds a customer's webhook endpoint and HMAC secret.
type Config struct {
	URL    string
	Secret string
}
