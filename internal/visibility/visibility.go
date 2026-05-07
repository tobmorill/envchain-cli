// Package visibility tracks per-project key visibility settings,
// allowing individual environment variable keys to be marked as
// hidden (redacted in output) or visible (shown in plain text).
package visibility

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/envchain/envchain-cli/internal/store"
)

// Level represents the visibility level of a key.
type Level string

const (
	LevelVisible Level = "visible"
	LevelHidden  Level = "hidden"
)

// Record holds visibility settings for a project's keys.
type Record struct {
	Project  string            `json:"project"`
	Settings map[string]Level  `json:"settings"`
}

// Manager persists and retrieves visibility records.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("visibility:%s", strings.ToLower(project))
}

// Set stores the visibility level for the given key within a project.
func (m *Manager) Set(project, key string, level Level) error {
	if project == "" {
		return errors.New("visibility: project must not be empty")
	}
	if key == "" {
		return errors.New("visibility: key must not be empty")
	}
	rec, err := m.getOrNew(project)
	if err != nil {
		return err
	}
	rec.Settings[strings.ToUpper(key)] = level
	return m.save(rec)
}

// Get returns the visibility level for the given key, defaulting to LevelVisible.
func (m *Manager) Get(project, key string) (Level, error) {
	rec, err := m.getOrNew(project)
	if err != nil {
		return LevelVisible, err
	}
	lvl, ok := rec.Settings[strings.ToUpper(key)]
	if !ok {
		return LevelVisible, nil
	}
	return lvl, nil
}

// GetAll returns all visibility settings for the given project.
func (m *Manager) GetAll(project string) (map[string]Level, error) {
	rec, err := m.getOrNew(project)
	if err != nil {
		return nil, err
	}
	out := make(map[string]Level, len(rec.Settings))
	for k, v := range rec.Settings {
		out[k] = v
	}
	return out, nil
}

// Delete removes the visibility setting for a specific key.
func (m *Manager) Delete(project, key string) error {
	rec, err := m.getOrNew(project)
	if err != nil {
		return err
	}
	delete(rec.Settings, strings.ToUpper(key))
	return m.save(rec)
}

func (m *Manager) getOrNew(project string) (*Record, error) {
	raw, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return &Record{Project: project, Settings: make(map[string]Level)}, nil
	}
	if err != nil {
		return nil, err
	}
	var rec Record
	if err := json.Unmarshal(raw, &rec); err != nil {
		return nil, fmt.Errorf("visibility: corrupt record: %w", err)
	}
	if rec.Settings == nil {
		rec.Settings = make(map[string]Level)
	}
	return &rec, nil
}

func (m *Manager) save(rec *Record) error {
	raw, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return m.st.Put(recordKey(rec.Project), raw)
}
