// Package retention manages per-project data retention policies,
// allowing users to define how long environment chain history and
// snapshots should be kept before automatic pruning.
package retention

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/envchain-cli/internal/store"
)

// ErrNotFound is returned when no retention policy exists for a project.
var ErrNotFound = errors.New("retention: policy not found")

// Policy defines how long data should be retained for a project.
type Policy struct {
	Project     string        `json:"project"`
	MaxAge      time.Duration `json:"max_age"`
	MaxVersions int           `json:"max_versions"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// Manager handles persistence of retention policies.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("retention:%s", project)
}

// Set persists a retention policy for the given project.
func (m *Manager) Set(p Policy) error {
	if p.Project == "" {
		return errors.New("retention: project name must not be empty")
	}
	if p.MaxAge < 0 {
		return errors.New("retention: max_age must not be negative")
	}
	if p.MaxVersions < 0 {
		return errors.New("retention: max_versions must not be negative")
	}
	p.UpdatedAt = time.Now().UTC()
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("retention: marshal: %w", err)
	}
	return m.st.Put(recordKey(p.Project), data)
}

// Get retrieves the retention policy for the given project.
// Returns ErrNotFound if no policy has been set.
func (m *Manager) Get(project string) (Policy, error) {
	data, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return Policy{}, ErrNotFound
	}
	if err != nil {
		return Policy{}, fmt.Errorf("retention: get: %w", err)
	}
	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return Policy{}, fmt.Errorf("retention: unmarshal: %w", err)
	}
	return p, nil
}

// Delete removes the retention policy for the given project.
func (m *Manager) Delete(project string) error {
	err := m.st.Delete(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// ShouldPrune reports whether a record timestamped at t should be pruned
// according to the policy. If MaxAge is zero it is not considered.
func (p Policy) ShouldPrune(t time.Time) bool {
	if p.MaxAge > 0 && time.Since(t) > p.MaxAge {
		return true
	}
	return false
}
