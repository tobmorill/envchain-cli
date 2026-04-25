// Package cache wraps the ttl.Cache to provide a project-scoped
// passphrase cache keyed by project name.
package cache

import (
	"time"

	"github.com/yourorg/envchain-cli/internal/ttl"
)

// DefaultTTL is the default duration a cached passphrase remains valid.
const DefaultTTL = 15 * time.Minute

// Manager holds cached passphrases for named projects.
type Manager struct {
	cache *ttl.Cache
	ttl   time.Duration
}

// New creates a Manager with the given per-entry TTL.
// If d is zero, DefaultTTL is used.
func New(d time.Duration) *Manager {
	if d == 0 {
		d = DefaultTTL
	}
	return &Manager{cache: ttl.New(), ttl: d}
}

// Store saves the passphrase for project under the configured TTL.
func (m *Manager) Store(project, passphrase string) {
	m.cache.Set(project, passphrase, m.ttl)
}

// Retrieve returns the cached passphrase for project.
// Returns ttl.ErrNotFound or ttl.ErrExpired if unavailable.
func (m *Manager) Retrieve(project string) (string, error) {
	return m.cache.Get(project)
}

// Invalidate removes the cached passphrase for project.
func (m *Manager) Invalidate(project string) {
	m.cache.Delete(project)
}

// PurgeExpired removes all entries that have passed their TTL.
func (m *Manager) PurgeExpired() {
	m.cache.Purge()
}
