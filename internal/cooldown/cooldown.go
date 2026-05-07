// Package cooldown tracks per-project cooldown periods, preventing
// repeated passphrase prompts or sensitive operations within a configurable
// window of time.
package cooldown

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/user/envchain-cli/internal/store"
)

// ErrEmptyProject is returned when an empty project name is provided.
var ErrEmptyProject = errors.New("cooldown: project name must not be empty")

// Record holds the cooldown state for a project.
type Record struct {
	Project   string        `json:"project"`
	StartedAt time.Time     `json:"started_at"`
	Duration  time.Duration `json:"duration"`
}

// IsActive reports whether the cooldown is still in effect.
func (r Record) IsActive() bool {
	return time.Now().Before(r.StartedAt.Add(r.Duration))
}

// ExpiresAt returns the time at which the cooldown expires.
func (r Record) ExpiresAt() time.Time {
	return r.StartedAt.Add(r.Duration)
}

// Manager manages cooldown records backed by a key-value store.
type Manager struct {
	s *store.Store
}

// New creates a new Manager using the given store.
func New(s *store.Store) *Manager {
	return &Manager{s: s}
}

func recordKey(project string) string {
	return fmt.Sprintf("cooldown:%s", project)
}

// Set records a cooldown for project lasting d. If a cooldown already exists
// it is overwritten.
func (m *Manager) Set(project string, d time.Duration) error {
	if project == "" {
		return ErrEmptyProject
	}
	rec := Record{
		Project:   project,
		StartedAt: time.Now().UTC(),
		Duration:  d,
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("cooldown: marshal: %w", err)
	}
	return m.s.Put(recordKey(project), b)
}

// Get returns the cooldown record for project. If no record exists,
// store.ErrNotFound is returned.
func (m *Manager) Get(project string) (Record, error) {
	if project == "" {
		return Record{}, ErrEmptyProject
	}
	b, err := m.s.Get(recordKey(project))
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(b, &rec); err != nil {
		return Record{}, fmt.Errorf("cooldown: unmarshal: %w", err)
	}
	return rec, nil
}

// Delete removes the cooldown record for project. It is not an error if no
// record exists.
func (m *Manager) Delete(project string) error {
	if project == "" {
		return ErrEmptyProject
	}
	return m.s.Delete(recordKey(project))
}

// IsActive reports whether an active cooldown exists for project.
func (m *Manager) IsActive(project string) (bool, error) {
	rec, err := m.Get(project)
	if errors.Is(err, store.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return rec.IsActive(), nil
}
