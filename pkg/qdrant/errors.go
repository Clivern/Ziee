// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package qdrant

import (
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotConfigured = errors.New("qdrant is not configured")
	ErrEmptyVector   = errors.New("vector must not be empty")
	ErrEmptyPoints   = errors.New("at least one point is required")
)

// IsIndexExists reports whether an error indicates an index already exists.
func IsIndexExists(err error) bool {
	if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
		return true
	}

	return strings.Contains(strings.ToLower(err.Error()), "already exists")
}
