// Package budget tracks cumulative value-byte consumption per project and
// enforces a configurable ceiling, complementing the per-chain quota checks.
package budget

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrExceeded is returned when a project would exceed its byte budget.
var ErrExceeded = errors.New("budget: byte limit exceeded")

// Record holds the tracked byte consumption for a single project.
type Record struct {
	Project string `json:"project"`
	UsedBytes int    `json:"used_bytes"`
	LimitBytes int   `json:"limit_bytes"`
}

// Manager persists budget records on disk.
type Manager struct {
	dir string
}

// New returns a Manager that stores records under dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) recordPath(project string) string {
	return filepath.Join(m.dir, project+".json")
}

// Set stores or replaces the budget record for project.
func (m *Manager) Set(project string, limitBytes int) error {
	if project == "" {
		return errors.New("budget: project name must not be empty")
	}
	if limitBytes <= 0 {
		return errors.New("budget: limit must be positive")
	}
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return err
	}
	r := Record{Project: project, LimitBytes: limitBytes}
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return os.WriteFile(m.recordPath(project), data, 0o600)
}

// Get returns the budget record for project.
func (m *Manager) Get(project string) (Record, error) {
	data, err := os.ReadFile(m.recordPath(project))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Record{}, fmt.Errorf("budget: no record for %q", project)
		}
		return Record{}, err
	}
	var r Record
	if err := json.Unmarshal(data, &r); err != nil {
		return Record{}, err
	}
	return r, nil
}

// Check returns ErrExceeded if usedBytes would exceed the stored limit.
func (m *Manager) Check(project string, usedBytes int) error {
	r, err := m.Get(project)
	if err != nil {
		// No record means no limit configured; allow.
		return nil
	}
	if usedBytes > r.LimitBytes {
		return fmt.Errorf("%w: project %q uses %d bytes, limit is %d",
			ErrExceeded, project, usedBytes, r.LimitBytes)
	}
	return nil
}

// Delete removes the budget record for project.
func (m *Manager) Delete(project string) error {
	err := os.Remove(m.recordPath(project))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
