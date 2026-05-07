// Package signal tracks per-project notification signals that can be
// raised and acknowledged by the user or automation tooling.
package signal

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
)

// Level describes the severity of a signal.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Record holds a single signal entry for a project.
type Record struct {
	Project     string    `json:"project"`
	Message     string    `json:"message"`
	Level       Level     `json:"level"`
	RaisedAt    time.Time `json:"raised_at"`
	Acknowledged bool     `json:"acknowledged"`
}

// Manager persists and retrieves signal records.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store directory.
func New(dir string) (*Manager, error) {
	st, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("signal: open store: %w", err)
	}
	return &Manager{st: st}, nil
}

func recordKey(project string) string {
	return "signal:" + project
}

// Raise creates or replaces the signal for the given project.
func (m *Manager) Raise(project, message string, level Level) error {
	if project == "" {
		return errors.New("signal: project must not be empty")
	}
	if message == "" {
		return errors.New("signal: message must not be empty")
	}
	rec := Record{
		Project:  project,
		Message:  message,
		Level:    level,
		RaisedAt: time.Now().UTC(),
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("signal: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), data)
}

// Get retrieves the current signal for a project.
func (m *Manager) Get(project string) (Record, error) {
	data, err := m.st.Get(recordKey(project))
	if err != nil {
		return Record{}, fmt.Errorf("signal: get %q: %w", project, err)
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("signal: unmarshal: %w", err)
	}
	return rec, nil
}

// Acknowledge marks the signal for a project as acknowledged.
func (m *Manager) Acknowledge(project string) error {
	rec, err := m.Get(project)
	if err != nil {
		return err
	}
	rec.Acknowledged = true
	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("signal: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), data)
}

// Delete removes the signal record for a project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(recordKey(project))
}
