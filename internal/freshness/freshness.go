// Package freshness tracks how recently a project's environment chain was
// loaded or refreshed, and exposes a simple staleness check.
package freshness

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrNotFound is returned when no freshness record exists for a project.
var ErrNotFound = errors.New("freshness: record not found")

// Record holds the last-touched timestamp for a project chain.
type Record struct {
	Project   string    `json:"project"`
	TouchedAt time.Time `json:"touched_at"`
}

// IsStale reports whether the record is older than the given threshold.
func (r Record) IsStale(threshold time.Duration) bool {
	return time.Since(r.TouchedAt) > threshold
}

// Manager persists and retrieves freshness records.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the provided store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("freshness::%s", project)
}

// Touch records the current time as the last-touched timestamp for project.
func (m *Manager) Touch(project string) error {
	if project == "" {
		return errors.New("freshness: project name must not be empty")
	}
	rec := Record{
		Project:   project,
		TouchedAt: time.Now().UTC(),
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("freshness: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), data)
}

// Get returns the freshness record for project.
func (m *Manager) Get(project string) (Record, error) {
	data, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, fmt.Errorf("freshness: get: %w", err)
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("freshness: unmarshal: %w", err)
	}
	return rec, nil
}

// Delete removes the freshness record for project.
func (m *Manager) Delete(project string) error {
	err := m.st.Delete(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return nil
	}
	return err
}
