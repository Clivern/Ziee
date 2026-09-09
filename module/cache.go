// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"time"

	"github.com/clivern/ziee/db"
)

// Cache is a kv-backed cache.
type Cache struct {
	kv db.KVRepository
}

// NewCache returns a kv-backed cache.
func NewCache(kv db.KVRepository) *Cache {
	return &Cache{kv: kv}
}

// Get returns a value for key, or empty if missing/expired.
func (c *Cache) Get(key string) (string, *time.Time, error) {
	item, err := c.kv.Get(key)
	if err != nil {
		return "", nil, err
	}
	if item == nil {
		return "", nil, nil
	}

	return item.Value, item.ExpiresAt, nil
}

// Set stores a value for key.
func (c *Cache) Set(key, value string, expiresAt *time.Time) error {
	return c.kv.Upsert(&db.KV{
		Key:       key,
		Value:     value,
		ExpiresAt: expiresAt,
	})
}

// DeleteExpired removes expired cache entries.
func (c *Cache) DeleteExpired() (int64, error) {
	return c.kv.DeleteExpired()
}
