// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package conf

import "time"

const (
	// SlowRequestThreshold is the latency at or above which a request is counted as slow.
	SlowRequestThreshold = time.Second

	// DefaultSessionDuration is the auth session TTL without remember-me.
	DefaultSessionDuration = 7 * 24 * time.Hour

	// RememberMeSessionDuration is the auth session TTL with remember-me.
	RememberMeSessionDuration = 30 * 24 * time.Hour

	// InviteExpiry is how long a workspace invite stays valid.
	InviteExpiry = 7 * 24 * time.Hour

	// AsyncClaimInterval is how often the async manager polls for due tasks.
	AsyncClaimInterval = 10 * time.Second

	// AsyncTaskMaxRetries is how many times a failed task is retried.
	AsyncTaskMaxRetries = 5

	// AsyncCompletedTasksRetention is how far back completed tasks are kept before cleanup.
	// Negative duration: cutoff = now + retention.
	AsyncCompletedTasksRetention = -7 * 24 * time.Hour

	// MaxUploadBytes is the maximum multipart upload size.
	MaxUploadBytes = 2 * 1024 * 1024

	// DefaultSearchLimit is used when a knowledge search request omits limit.
	DefaultSearchLimit = 10
)
