// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package v1

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

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

func (c *Clause) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]yaml.Node
	err := value.Decode(&raw)
	if err != nil {
		return err
	}

	for key, node := range raw {
		switch key {
		case "files":
			err = node.Decode(&c.Files)
		case "max_files_changed":
			var n int
			err = node.Decode(&n)
			c.MaxFilesChanged = &n
		case "min_files_changed":
			var n int
			err = node.Decode(&n)
			c.MinFilesChanged = &n
		case "title":
			err = node.Decode(&c.Title)
		case "body":
			err = node.Decode(&c.Body)
		case "author_in":
			err = node.Decode(&c.AuthorIn)
		case "author_not_in":
			err = node.Decode(&c.AuthorNotIn)
		case "author_in_team":
			err = node.Decode(&c.AuthorInTeam)
		case "author_not_in_team":
			err = node.Decode(&c.AuthorNotInTeam)
		case "intention":
			err = node.Decode(&c.Intention)
		case "label":
			err = node.Decode(&c.Label)
		case "check":
			err = node.Decode(&c.Check)
		case "approvals":
			var a Approvals
			err = node.Decode(&a)
			c.Approvals = &a
		default:
			return fmt.Errorf("spec: unknown when key %q", key)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *AllowEntry) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]yaml.Node
	err := value.Decode(&raw)
	if err != nil {
		return err
	}

	for key, node := range raw {
		switch key {
		case "permission":
			err = node.Decode(&e.Permission)
		case "teams":
			err = node.Decode(&e.Teams)
		case "users":
			err = node.Decode(&e.Users)
		default:
			return fmt.Errorf("spec: unknown allow key %q", key)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BatchSize) UnmarshalYAML(value *yaml.Node) error {
	var n int
	err := value.Decode(&n)
	if err == nil {
		b.Value = &n
		return nil
	}

	var s struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	}
	err = value.Decode(&s)
	if err != nil {
		return err
	}

	b.Min = s.Min
	b.Max = s.Max

	return nil
}

func (b BatchSize) MarshalYAML() (any, error) {
	if b.Value != nil {
		return *b.Value, nil
	}

	return struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	}{Min: b.Min, Max: b.Max}, nil
}
