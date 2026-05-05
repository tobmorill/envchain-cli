// Package timeout provides per-project session timeout configuration,
// allowing chains to be automatically locked after a configurable idle period.
package timeout

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrNotFound is returned when no timeout rule exists for a project.
var ErrNotFound = errors.New("timeout: rule not found")

// Rule describes the idle-lock policy for a single project.
type Rule struct {
	Project  string        `json:"project"`
	Duration time.Duration `json:"duration"`
	Enabled  bool          `json:"enabled"`
}

// Manager persists timeout rules on disk.
type Manager struct {
	dir string
}

// New returns a Manager that stores rules under dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) path(project string) string {
	return filepath.Join(m.dir, project+".json")
}

// Set writes or replaces the timeout rule for project.
func (m *Manager) Set(r Rule) error {
	if r.Project == "" {
		return fmt.Errorf("timeout: project name must not be empty")
	}
	if r.Duration < 0 {
		return fmt.Errorf("timeout: duration must be non-negative")
	}
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return fmt.Errorf("timeout: mkdir: %w", err)
	}
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("timeout: marshal: %w", err)
	}
	return os.WriteFile(m.path(r.Project), data, 0o600)
}

// Get returns the timeout rule for project.
func (m *Manager) Get(project string) (Rule, error) {
	data, err := os.ReadFile(m.path(project))
	if errors.Is(err, os.ErrNotExist) {
		return Rule{}, ErrNotFound
	}
	if err != nil {
		return Rule{}, fmt.Errorf("timeout: read: %w", err)
	}
	var r Rule
	if err := json.Unmarshal(data, &r); err != nil {
		return Rule{}, fmt.Errorf("timeout: unmarshal: %w", err)
	}
	return r, nil
}

// Delete removes the timeout rule for project.
func (m *Manager) Delete(project string) error {
	err := os.Remove(m.path(project))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	return err
}

// IsDue reports whether the given lastActivity time has exceeded the project's
// configured idle duration. Returns false if no rule exists or the rule is
// disabled.
func (m *Manager) IsDue(project string, lastActivity time.Time) (bool, error) {
	r, err := m.Get(project)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !r.Enabled || r.Duration == 0 {
		return false, nil
	}
	return time.Since(lastActivity) >= r.Duration, nil
}
