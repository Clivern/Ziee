// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package spec

import (
	"fmt"

	v1 "github.com/clivern/ziee/policy/spec/v1"

	"gopkg.in/yaml.v3"
)

// Parse loads a `.ziee.yml` using the schema for its `version` field.
func Parse(data []byte) (*v1.File, error) {
	var head Header
	yaml.Unmarshal(data, &head)

	switch head.Version {
	case v1.Version:
		return v1.Parse(data)
	default:
		return nil, fmt.Errorf("policy: unsupported version %q", head.Version)
	}
}
