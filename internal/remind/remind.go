// Package remind provides scheduled reminders to rotate passphrases
// or review environment variable chains after a configurable interval.
package remind

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrNoReminder is returned when no reminder exists for a project.
var ErrNoReminder = errors.New("remind: no reminder set")

// Reminder holds the scheduled reminder metadata for a project.
type Reminder struct {
	Project   string        `json:"project"`
	Interval  time.Duration `json:"interval"`
	LastReset time.Time     `json:"last_reset"`
	Message   string        `json:"message,omitempty"`
}

// IsDue reports whether the reminder interval has elapsed since LastReset.
func (r *Reminder) IsDue() bool {
	return time.Since(r.LastReset) >= r.Interval
}

// Manager manages reminder records on disk.
type Manager struct {
	dir string
}

// New returns a Manager that stores reminders under dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) path(project string) string {
	return filepath.Join(m.dir, project+".json")
}

// Set creates or overwrites the reminder for project.
func (m *Manager) Set(r Reminder) error {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return err
	}
	r.LastReset = time.Now()
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return os.WriteFile(m.path(r.Project), data, 0o600)
}

// Get returns the reminder for project, or ErrNoReminder if none exists.
func (m *Manager) Get(project string) (Reminder, error) {
	data, err := os.ReadFile(m.path(project))
	if errors.Is(err, os.ErrNotExist) {
		return Reminder{}, ErrNoReminder
	}
	if err != nil {
		return Reminder{}, err
	}
	var r Reminder
	if err := json.Unmarshal(data, &r); err != nil {
		return Reminder{}, err
	}
	return r, nil
}

// Reset marks the reminder as acknowledged, restarting the interval.
func (m *Manager) Reset(project string) error {
	r, err := m.Get(project)
	if err != nil {
		return err
	}
	r.LastReset = time.Now()
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return os.WriteFile(m.path(project), data, 0o600)
}

// Delete removes the reminder for project.
func (m *Manager) Delete(project string) error {
	err := os.Remove(m.path(project))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNoReminder
	}
	return err
}
