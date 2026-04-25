// Package ttl provides time-to-live management for cached passphrases
// and unlocked chain sessions.
package ttl

import (
	"errors"
	"sync"
	"time"
)

// ErrExpired is returned when a cached entry has passed its TTL.
var ErrExpired = errors.New("ttl: entry has expired")

// ErrNotFound is returned when no entry exists for the given key.
var ErrNotFound = errors.New("ttl: entry not found")

// entry holds a value and its expiry time.
type entry struct {
	value     string
	expiresAt time.Time
}

// Cache is a thread-safe in-memory store with per-entry expiry.
type Cache struct {
	mu      sync.Mutex
	entries map[string]entry
}

// New creates an empty Cache.
func New() *Cache {
	return &Cache{entries: make(map[string]entry)}
}

// Set stores value under key and marks it to expire after ttl.
func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves the value for key. Returns ErrNotFound if the key
// was never set, or ErrExpired if the entry has passed its TTL.
func (c *Cache) Get(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return "", ErrNotFound
	}
	if time.Now().After(e.expiresAt) {
		delete(c.entries, key)
		return "", ErrExpired
	}
	return e.value, nil
}

// Delete removes the entry for key. It is a no-op if the key does not exist.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Purge removes all expired entries from the cache.
func (c *Cache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, e := range c.entries {
		if now.After(e.expiresAt) {
			delete(c.entries, k)
		}
	}
}
