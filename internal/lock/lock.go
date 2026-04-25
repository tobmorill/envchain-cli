// Package lock provides session-level chain locking to prevent concurrent
// access and support automatic expiry of unlocked chains.
package lock

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrLocked is returned when a chain is already locked (session expired or manually locked).
var ErrLocked = errors.New("chain is locked")

// ErrNotLocked is returned when trying to unlock a chain that has no active session.
var ErrNotLocked = errors.New("chain is not unlocked")

// Session represents an active unlock session for a chain.
type Session struct {
	Chain     string    `json:"chain"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired reports whether the session has passed its expiry time.
func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Manager handles lock/unlock sessions stored on disk.
type Manager struct {
	dir string
}

// NewManager returns a Manager that stores session files under dir.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) sessionPath(chain string) string {
	return filepath.Join(m.dir, chain+".session.json")
}

// Unlock creates a session for chain that expires after ttl.
func (m *Manager) Unlock(chain string, ttl time.Duration) error {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return err
	}
	s := Session{
		Chain:     chain,
		ExpiresAt: time.Now().Add(ttl),
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(m.sessionPath(chain), data, 0o600)
}

// Lock removes the active session for chain.
func (m *Manager) Lock(chain string) error {
	path := m.sessionPath(chain)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return ErrNotLocked
	}
	return os.Remove(path)
}

// IsUnlocked reports whether chain has a valid, non-expired session.
func (m *Manager) IsUnlocked(chain string) (bool, error) {
	s, err := m.load(chain)
	if errors.Is(err, ErrLocked) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if s.IsExpired() {
		_ = os.Remove(m.sessionPath(chain))
		return false, nil
	}
	return true, nil
}

func (m *Manager) load(chain string) (Session, error) {
	data, err := os.ReadFile(m.sessionPath(chain))
	if errors.Is(err, os.ErrNotExist) {
		return Session{}, ErrLocked
	}
	if err != nil {
		return Session{}, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return Session{}, err
	}
	return s, nil
}
