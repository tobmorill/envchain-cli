// Package expiry tracks per-project chain expiration dates and reports
// whether a chain has passed its configured deadline.
package expiry

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrNotFound is returned when no expiry record exists for a project.
var ErrNotFound = errors.New("expiry: record not found")

// Record holds the expiration metadata for a single project chain.
type Record struct {
	Project   string    `json:"project"`
	ExpiresAt time.Time `json:"expires_at"`
	Note      string    `json:"note,omitempty"`
}

// IsExpired reports whether the record's deadline has passed.
func (r Record) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// Manager persists expiry records using the underlying key-value store.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by a store rooted at dir.
func New(dir string) (*Manager, error) {
	st, err := store.New(filepath.Join(dir, "expiry.db"))
	if err != nil {
		return nil, fmt.Errorf("expiry: open store: %w", err)
	}
	return &Manager{st: st}, nil
}

func recordKey(project string) string {
	return "expiry:" + project
}

// Set stores or overwrites the expiry record for the given project.
func (m *Manager) Set(project string, expiresAt time.Time, note string) error {
	if project == "" {
		return errors.New("expiry: project name must not be empty")
	}
	rec := Record{Project: project, ExpiresAt: expiresAt, Note: note}
	b, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("expiry: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), b)
}

// Get retrieves the expiry record for the given project.
func (m *Manager) Get(project string) (Record, error) {
	b, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, fmt.Errorf("expiry: get: %w", err)
	}
	var rec Record
	if err := json.Unmarshal(b, &rec); err != nil {
		return Record{}, fmt.Errorf("expiry: unmarshal: %w", err)
	}
	return rec, nil
}

// Delete removes the expiry record for the given project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(recordKey(project))
}
