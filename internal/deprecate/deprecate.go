// Package deprecate tracks deprecated keys within a project's environment chain.
// It allows marking specific keys as deprecated with an optional replacement hint,
// and querying which keys in a given set are deprecated.
package deprecate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrNotFound is returned when no deprecation record exists for a project.
var ErrNotFound = errors.New("deprecate: no record found")

// Entry holds deprecation metadata for a single key.
type Entry struct {
	Key         string `json:"key"`
	Replacement string `json:"replacement,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// Record is the full set of deprecated keys for a project.
type Record struct {
	Project string  `json:"project"`
	Entries []Entry `json:"entries"`
}

// Manager persists and retrieves deprecation records.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("deprecate:%s", strings.ToLower(project))
}

// Set overwrites the deprecation record for project.
func (m *Manager) Set(project string, entries []Entry) error {
	if project == "" {
		return errors.New("deprecate: project name must not be empty")
	}
	rec := Record{Project: project, Entries: entries}
	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("deprecate: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), data)
}

// Get returns the deprecation record for project.
func (m *Manager) Get(project string) (Record, error) {
	data, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, fmt.Errorf("deprecate: get: %w", err)
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("deprecate: unmarshal: %w", err)
	}
	return rec, nil
}

// Delete removes the deprecation record for project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(recordKey(project))
}

// Check returns the subset of keys (from the provided slice) that are
// deprecated according to the stored record for project.
func (m *Manager) Check(project string, keys []string) ([]Entry, error) {
	rec, err := m.Get(project)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[strings.ToLower(k)] = struct{}{}
	}
	var hits []Entry
	for _, e := range rec.Entries {
		if _, ok := keySet[strings.ToLower(e.Key)]; ok {
			hits = append(hits, e)
		}
	}
	return hits, nil
}
