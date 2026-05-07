// Package shield provides per-project key protection rules,
// preventing accidental overwrite or deletion of critical environment variables.
package shield

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrShielded is returned when an operation is blocked by a shield rule.
var ErrShielded = errors.New("key is shielded")

// ErrEmptyProject is returned when a blank project name is provided.
var ErrEmptyProject = errors.New("project name must not be empty")

type record struct {
	Keys []string `json:"keys"`
}

// Manager manages shield rules backed by a key-value store.
type Manager struct {
	st *store.Store
}

// New returns a Manager using the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return "shield:" + project
}

// Set replaces the shielded key list for project with keys.
func (m *Manager) Set(project string, keys []string) error {
	if strings.TrimSpace(project) == "" {
		return ErrEmptyProject
	}
	norm := dedup(keys)
	b, err := json.Marshal(record{Keys: norm})
	if err != nil {
		return fmt.Errorf("shield: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), b)
}

// Get returns the shielded keys for project.
// If no rule exists, an empty slice and nil error are returned.
func (m *Manager) Get(project string) ([]string, error) {
	b, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("shield: get: %w", err)
	}
	var r record
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, fmt.Errorf("shield: unmarshal: %w", err)
	}
	return r.Keys, nil
}

// Guard returns ErrShielded if key is in the shielded set for project.
func (m *Manager) Guard(project, key string) error {
	keys, err := m.Get(project)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if strings.EqualFold(k, key) {
			return fmt.Errorf("%w: %s", ErrShielded, key)
		}
	}
	return nil
}

// Delete removes the shield rule for project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(recordKey(project))
}

func dedup(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, v := range in {
		upper := strings.ToUpper(strings.TrimSpace(v))
		if upper == "" {
			continue
		}
		if _, ok := seen[upper]; ok {
			continue
		}
		seen[upper] = struct{}{}
		out = append(out, upper)
	}
	return out
}
