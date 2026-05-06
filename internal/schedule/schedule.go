// Package schedule manages periodic refresh schedules for project chains.
// A schedule defines how often a chain's passphrase should be rotated or
// re-prompted, expressed as a duration stored per project.
package schedule

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrNotFound is returned when no schedule exists for a project.
var ErrNotFound = errors.New("schedule: not found")

// Record holds the schedule configuration for a single project.
type Record struct {
	Project  string        `json:"project"`
	Interval time.Duration `json:"interval_ns"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// IsDue reports whether the schedule interval has elapsed since UpdatedAt.
func (r Record) IsDue() bool {
	if r.Interval <= 0 {
		return false
	}
	return time.Since(r.UpdatedAt) >= r.Interval
}

// Manager persists and retrieves schedule records.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store directory.
func New(dir string) (*Manager, error) {
	st, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("schedule: open store: %w", err)
	}
	return &Manager{st: st}, nil
}

func recordKey(project string) string {
	return "schedule:" + project
}

// Set stores or overwrites the schedule for project.
func (m *Manager) Set(project string, interval time.Duration) error {
	if project == "" {
		return errors.New("schedule: project must not be empty")
	}
	if interval <= 0 {
		return errors.New("schedule: interval must be positive")
	}
	now := time.Now().UTC()
	rec := Record{
		Project:   project,
		Interval:  interval,
		CreatedAt: now,
		UpdatedAt: now,
	}
	// Preserve original CreatedAt if record already exists.
	if existing, err := m.Get(project); err == nil {
		rec.CreatedAt = existing.CreatedAt
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("schedule: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), b)
}

// Get retrieves the schedule record for project.
func (m *Manager) Get(project string) (Record, error) {
	b, err := m.st.Get(recordKey(project))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return Record{}, ErrNotFound
		}
		return Record{}, fmt.Errorf("schedule: get: %w", err)
	}
	var rec Record
	if err := json.Unmarshal(b, &rec); err != nil {
		return Record{}, fmt.Errorf("schedule: unmarshal: %w", err)
	}
	return rec, nil
}

// Touch updates the UpdatedAt timestamp, resetting the due window.
func (m *Manager) Touch(project string) error {
	rec, err := m.Get(project)
	if err != nil {
		return err
	}
	rec.UpdatedAt = time.Now().UTC()
	b, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("schedule: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), b)
}

// Delete removes the schedule for project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(recordKey(project))
}
