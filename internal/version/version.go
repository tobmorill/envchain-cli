// Package version tracks per-project chain version numbers, allowing
// consumers to detect when a chain has been modified since it was last read.
package version

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrNotFound is returned when no version record exists for a project.
var ErrNotFound = errors.New("version: record not found")

// Record holds the monotonic version counter for a single project chain.
type Record struct {
	Project string `json:"project"`
	Version uint64 `json:"version"`
}

// Manager persists version records using the underlying store.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("version::%s", project)
}

// Get returns the current version record for project.
// Returns ErrNotFound if no record has been written yet.
func (m *Manager) Get(project string) (Record, error) {
	raw, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, fmt.Errorf("version get: %w", err)
	}
	var r Record
	if err := json.Unmarshal(raw, &r); err != nil {
		return Record{}, fmt.Errorf("version decode: %w", err)
	}
	return r, nil
}

// Bump increments the version counter for project by 1, creating it at 1 if
// it does not yet exist. The updated Record is returned.
func (m *Manager) Bump(project string) (Record, error) {
	r, err := m.Get(project)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Record{}, err
	}
	r.Project = project
	r.Version++
	raw, err := json.Marshal(r)
	if err != nil {
		return Record{}, fmt.Errorf("version encode: %w", err)
	}
	if err := m.st.Put(recordKey(project), raw); err != nil {
		return Record{}, fmt.Errorf("version put: %w", err)
	}
	return r, nil
}

// Reset removes the version record for project.
func (m *Manager) Reset(project string) error {
	if err := m.st.Delete(recordKey(project)); err != nil && !errors.Is(err, store.ErrNotFound) {
		return fmt.Errorf("version reset: %w", err)
	}
	return nil
}
