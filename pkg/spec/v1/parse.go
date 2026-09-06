// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package v1

import "gopkg.in/yaml.v3"

const Version = "1.0.0"

// Parse unmarshals a version 1.0.0 `.ziee.yml` document.
func Parse(data []byte) (*File, error) {
	var file File
	err := yaml.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
