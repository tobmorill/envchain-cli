// Package immutable provides a mechanism to mark specific environment variable
// keys within a project as immutable (read-only), preventing accidental
// overwrites during merge, import, or manual edits.
package immutable

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrEmptyProject is returned when an empty project name is supplied.
var ErrEmptyProject = errors.New("immutable: project name must not be empty")

// ErrEmptyKey is returned when an empty key is supplied.
var ErrEmptyKey = errors.New("immutable: key must not be empty")

// Manager manages immutable key sets backed by an encrypted store.
type Manager struct {
	st *store.Store
}

// New creates a Manager using the provided store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("immutable:%s", project)
}

// Set replaces the immutable key set for project with keys.
// Duplicate keys are deduplicated and the slice is sorted for determinism.
func (m *Manager) Set(project string, keys []string, passphrase string) error {
	if project == "" {
		return ErrEmptyProject
	}
	seen := make(map[string]struct{}, len(keys))
	uniq := keys[:0:0]
	for _, k := range keys {
		if k == "" {
			return ErrEmptyKey
		}
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			uniq = append(uniq, k)
		}
	}
	sort.Strings(uniq)
	data, err := json.Marshal(uniq)
	if err != nil {
		return fmt.Errorf("immutable: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), data, passphrase)
}

// Get returns the sorted immutable key set for project.
func (m *Manager) Get(project, passphrase string) ([]string, error) {
	if project == "" {
		return nil, ErrEmptyProject
	}
	data, err := m.st.Get(recordKey(project), passphrase)
	if err != nil {
		return nil, err
	}
	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("immutable: unmarshal: %w", err)
	}
	return keys, nil
}

// IsImmutable reports whether key is marked immutable for project.
func (m *Manager) IsImmutable(project, key, passphrase string) (bool, error) {
	keys, err := m.Get(project, passphrase)
	if err != nil {
		return false, err
	}
	for _, k := range keys {
		if k == key {
			return true, nil
		}
	}
	return false, nil
}

// Delete removes the immutable key set for project.
func (m *Manager) Delete(project string) error {
	if project == "" {
		return ErrEmptyProject
	}
	return m.st.Delete(recordKey(project))
}
