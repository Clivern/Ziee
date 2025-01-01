// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package storage

// DocumentKey returns the object path for a document id: "{uuid}.txt".
func DocumentKey(id string) string {
	return id + ".txt"
}
