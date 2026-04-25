// Package history tracks per-project passphrase usage history,
// recording timestamps of chain access events for auditing and
// session-awareness purposes.
package history

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrNotFound is returned when no history exists for a project.
var ErrNotFound = errors.New("history: no history found for project")

// Entry represents a single access record for a project chain.
type Entry struct {
	Project   string    `json:"project"`
	AccessedAt time.Time `json:"accessed_at"`
	Action    string    `json:"action"`
}

// Manager handles reading and writing history records.
type Manager struct {
	dir string
}

// New creates a Manager that stores history files under dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) historyPath(project string) string {
	return filepath.Join(m.dir, project+".json")
}

// Record appends an access entry for the given project and action.
func (m *Manager) Record(project, action string) error {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return err
	}
	entries, _ := m.ReadAll(project) // ignore not-found
	entries = append(entries, Entry{
		Project:    project,
		AccessedAt: time.Now().UTC(),
		Action:     action,
	})
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	return os.WriteFile(m.historyPath(project), data, 0o600)
}

// ReadAll returns all recorded entries for the given project.
func (m *Manager) ReadAll(project string) ([]Entry, error) {
	data, err := os.ReadFile(m.historyPath(project))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// Clear removes all history for the given project.
func (m *Manager) Clear(project string) error {
	err := os.Remove(m.historyPath(project))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
