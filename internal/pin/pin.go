// Package pin provides per-project pinned variable management.
// A pinned variable is a key that is always injected into the shell
// environment regardless of which chain is currently loaded.
package pin

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrNotFound is returned when no pin file exists for a project.
var ErrNotFound = errors.New("pin: no pinned keys found")

// Manager persists pinned keys for projects.
type Manager struct {
	dir string
}

type pinFile struct {
	Keys []string `json:"keys"`
}

// New creates a Manager that stores pin files under dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) path(project string) string {
	return filepath.Join(m.dir, strings.ToLower(project)+".pin.json")
}

// Set replaces the pinned key list for project.
// Keys are deduplicated and normalised to upper-case.
func (m *Manager) Set(project string, keys []string) error {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return err
	}
	seen := make(map[string]struct{}, len(keys))
	uniq := keys[:0:0]
	for _, k := range keys {
		norm := strings.ToUpper(strings.TrimSpace(k))
		if norm == "" {
			continue
		}
		if _, ok := seen[norm]; ok {
			continue
		}
		seen[norm] = struct{}{}
		uniq = append(uniq, norm)
	}
	data, err := json.MarshalIndent(pinFile{Keys: uniq}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path(project), data, 0o600)
}

// Get returns the pinned keys for project.
func (m *Manager) Get(project string) ([]string, error) {
	data, err := os.ReadFile(m.path(project))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var pf pinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	return pf.Keys, nil
}

// Delete removes the pin file for project.
func (m *Manager) Delete(project string) error {
	err := os.Remove(m.path(project))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
