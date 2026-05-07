// Package readonly tracks which projects have been marked read-only,
// preventing accidental writes to sensitive environment variable chains.
package readonly

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrReadOnly is returned when a write is attempted on a read-only project.
var ErrReadOnly = errors.New("project is marked read-only")

// Record holds the read-only state for a project.
type Record struct {
	Project  string `json:"project"`
	ReadOnly bool   `json:"read_only"`
}

// Manager manages read-only flags for projects.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("readonly:%s", project)
}

// Set marks a project as read-only (enabled=true) or writable (enabled=false).
func (m *Manager) Set(project string, enabled bool) error {
	if project == "" {
		return errors.New("project name must not be empty")
	}
	rec := Record{Project: project, ReadOnly: enabled}
	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return m.st.Put(recordKey(project), data)
}

// Get returns the Record for the given project.
// If no record exists, it returns a Record with ReadOnly=false.
func (m *Manager) Get(project string) (Record, error) {
	data, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return Record{Project: project, ReadOnly: false}, nil
	}
	if err != nil {
		return Record{}, fmt.Errorf("get: %w", err)
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("unmarshal: %w", err)
	}
	return rec, nil
}

// IsReadOnly returns true if the project is currently marked read-only.
func (m *Manager) IsReadOnly(project string) (bool, error) {
	rec, err := m.Get(project)
	if err != nil {
		return false, err
	}
	return rec.ReadOnly, nil
}

// Delete removes the read-only record for the given project.
func (m *Manager) Delete(project string) error {
	if project == "" {
		return errors.New("project name must not be empty")
	}
	return m.st.Delete(recordKey(project))
}

// Guard returns ErrReadOnly if the project is marked read-only.
func (m *Manager) Guard(project string) error {
	ok, err := m.IsReadOnly(project)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%w: %s", ErrReadOnly, project)
	}
	return nil
}
