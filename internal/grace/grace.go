// Package grace manages per-project grace periods — a window of time after
// a chain is modified during which destructive operations (delete, rotate)
// are blocked. This provides a safety net against accidental data loss.
package grace

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrInGracePeriod is returned when an operation is blocked by an active grace period.
var ErrInGracePeriod = errors.New("grace: operation blocked by active grace period")

// ErrEmptyProject is returned when the project name is empty.
var ErrEmptyProject = errors.New("grace: project name must not be empty")

// Record holds the grace period configuration for a project.
type Record struct {
	Project   string        `json:"project"`
	Duration  time.Duration `json:"duration"`
	Activated time.Time     `json:"activated"`
}

// IsActive reports whether the grace period is still in effect.
func (r Record) IsActive() bool {
	return time.Now().Before(r.Activated.Add(r.Duration))
}

// ExpiresAt returns the time at which the grace period ends.
func (r Record) ExpiresAt() time.Time {
	return r.Activated.Add(r.Duration)
}

// Manager persists and retrieves grace period records.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("grace:%s", project)
}

// Set stores a grace period for the given project, starting now.
func (m *Manager) Set(project string, d time.Duration) error {
	if project == "" {
		return ErrEmptyProject
	}
	if d <= 0 {
		return fmt.Errorf("grace: duration must be positive, got %s", d)
	}
	rec := Record{
		Project:   project,
		Duration:  d,
		Activated: time.Now().UTC(),
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("grace: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), data)
}

// Get retrieves the grace period record for the given project.
func (m *Manager) Get(project string) (Record, error) {
	if project == "" {
		return Record{}, ErrEmptyProject
	}
	data, err := m.st.Get(recordKey(project))
	if err != nil {
		return Record{}, fmt.Errorf("grace: %w", err)
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("grace: unmarshal: %w", err)
	}
	return rec, nil
}

// Guard returns ErrInGracePeriod if an active grace period exists for project.
func (m *Manager) Guard(project string) error {
	rec, err := m.Get(project)
	if errors.Is(err, store.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if rec.IsActive() {
		return fmt.Errorf("%w: expires at %s", ErrInGracePeriod, rec.ExpiresAt().Format(time.RFC3339))
	}
	return nil
}

// Delete removes the grace period record for the given project.
func (m *Manager) Delete(project string) error {
	if project == "" {
		return ErrEmptyProject
	}
	return m.st.Delete(recordKey(project))
}
