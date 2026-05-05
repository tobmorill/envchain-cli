// Package access tracks per-project access counts and last-used timestamps,
// enabling usage analytics and stale-chain detection.
package access

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
)

const keyPrefix = "access:"

// Record holds access metadata for a single project chain.
type Record struct {
	Project   string    `json:"project"`
	Count     int64     `json:"count"`
	FirstUsed time.Time `json:"first_used"`
	LastUsed  time.Time `json:"last_used"`
}

// Manager persists and retrieves access records.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by st.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

// Touch increments the access counter for project and updates LastUsed.
// If no record exists yet, one is created with FirstUsed set to now.
func (m *Manager) Touch(project string) error {
	if project == "" {
		return errors.New("access: project name must not be empty")
	}
	rec, err := m.Get(project)
	if err != nil {
		now := time.Now().UTC()
		rec = Record{Project: project, FirstUsed: now}
	}
	rec.Count++
	rec.LastUsed = time.Now().UTC()
	return m.save(rec)
}

// Get returns the access record for project.
// Returns an error wrapping store.ErrNotFound when no record exists.
func (m *Manager) Get(project string) (Record, error) {
	raw, err := m.st.Get(keyPrefix + project)
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(raw, &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

// Reset clears the access record for project.
func (m *Manager) Reset(project string) error {
	if project == "" {
		return errors.New("access: project name must not be empty")
	}
	return m.st.Delete(keyPrefix + project)
}

func (m *Manager) save(rec Record) error {
	raw, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return m.st.Put(keyPrefix+rec.Project, raw)
}
