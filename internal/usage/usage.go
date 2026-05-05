// Package usage tracks per-project command invocation counts and last-used
// timestamps, giving operators a lightweight signal for which chains are
// actively referenced.
package usage

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
)

// Record holds aggregated usage data for a single project.
type Record struct {
	Project   string    `json:"project"`
	Count     int64     `json:"count"`
	FirstUsed time.Time `json:"first_used"`
	LastUsed  time.Time `json:"last_used"`
}

// Manager persists usage records via the key-value store.
type Manager struct {
	s *store.Store
}

// New returns a Manager backed by the given store.
func New(s *store.Store) *Manager {
	return &Manager{s: s}
}

func recordKey(project string) string {
	return fmt.Sprintf("usage:%s", project)
}

// Touch increments the invocation counter for project and updates timestamps.
// If no record exists it is created.
func (m *Manager) Touch(project string) error {
	if project == "" {
		return errors.New("usage: project name must not be empty")
	}

	rec, err := m.Get(project)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if rec == nil {
		rec = &Record{
			Project:   project,
			FirstUsed: now,
		}
	}
	rec.Count++
	rec.LastUsed = now

	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("usage: marshal: %w", err)
	}
	return m.s.Put(recordKey(project), data)
}

// Get returns the Record for project, or nil if none exists.
func (m *Manager) Get(project string) (*Record, error) {
	data, err := m.s.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("usage: get: %w", err)
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, fmt.Errorf("usage: unmarshal: %w", err)
	}
	return &rec, nil
}

// Reset removes the usage record for project.
func (m *Manager) Reset(project string) error {
	if project == "" {
		return errors.New("usage: project name must not be empty")
	}
	err := m.s.Delete(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return nil
	}
	return err
}
