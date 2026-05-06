// Package priority manages per-project key priority levels, allowing
// operators to mark environment variables as critical, normal, or low
// importance for use in ordering, display, and policy enforcement.
package priority

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/envchain/envchain-cli/internal/store"
)

// Level represents the importance of an environment variable.
type Level int

const (
	Low    Level = -1
	Normal Level = 0
	High   Level = 1
)

// ErrInvalidLevel is returned when an unrecognised priority level is provided.
var ErrInvalidLevel = errors.New("priority: invalid level")

// ParseLevel converts a string ("low", "normal", "high") to a Level.
func ParseLevel(s string) (Level, error) {
	switch s {
	case "low":
		return Low, nil
	case "normal":
		return Normal, nil
	case "high":
		return High, nil
	default:
		return Normal, fmt.Errorf("%w: %q", ErrInvalidLevel, s)
	}
}

// String returns the canonical string representation of a Level.
func (l Level) String() string {
	switch l {
	case Low:
		return "low"
	case High:
		return "high"
	default:
		return "normal"
	}
}

// Record holds the priority levels for all keys in a project.
type Record struct {
	Levels map[string]Level `json:"levels"`
}

// Manager persists priority records using the shared store.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return "priority:" + project
}

// Set persists the priority level for a key within a project.
func (m *Manager) Set(project, key string, level Level) error {
	if project == "" {
		return errors.New("priority: project must not be empty")
	}
	rec, err := m.get(project)
	if err != nil {
		return err
	}
	rec.Levels[key] = level
	return m.save(project, rec)
}

// Get returns the priority level for a key. Normal is returned when the key
// has no explicit priority set.
func (m *Manager) Get(project, key string) (Level, error) {
	rec, err := m.get(project)
	if err != nil {
		return Normal, err
	}
	level, ok := rec.Levels[key]
	if !ok {
		return Normal, nil
	}
	return level, nil
}

// Delete removes the priority entry for a key within a project.
func (m *Manager) Delete(project, key string) error {
	rec, err := m.get(project)
	if err != nil {
		return err
	}
	delete(rec.Levels, key)
	return m.save(project, rec)
}

// GetAll returns all key→level mappings for a project.
func (m *Manager) GetAll(project string) (map[string]Level, error) {
	rec, err := m.get(project)
	if err != nil {
		return nil, err
	}
	out := make(map[string]Level, len(rec.Levels))
	for k, v := range rec.Levels {
		out[k] = v
	}
	return out, nil
}

func (m *Manager) get(project string) (Record, error) {
	raw, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return Record{Levels: make(map[string]Level)}, nil
	}
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(raw, &rec); err != nil {
		return Record{}, fmt.Errorf("priority: corrupt record: %w", err)
	}
	if rec.Levels == nil {
		rec.Levels = make(map[string]Level)
	}
	return rec, nil
}

func (m *Manager) save(project string, rec Record) error {
	raw, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return m.st.Put(recordKey(project), raw)
}
