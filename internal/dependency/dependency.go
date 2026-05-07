// Package dependency tracks inter-project dependency relationships,
// allowing a project to declare that it depends on one or more other projects.
package dependency

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrSelfDependency is returned when a project lists itself as a dependency.
var ErrSelfDependency = errors.New("dependency: project cannot depend on itself")

// ErrEmptyProject is returned when an empty project name is supplied.
var ErrEmptyProject = errors.New("dependency: project name must not be empty")

// Manager persists and retrieves dependency lists.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store directory.
func New(dir string) (*Manager, error) {
	st, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("dependency: open store: %w", err)
	}
	return &Manager{st: st}, nil
}

func recordKey(project string) string {
	return "dep:" + project
}

// Set replaces the dependency list for project with deps.
// Duplicates are removed and order is preserved (first occurrence wins).
func (m *Manager) Set(project string, deps []string) error {
	if project == "" {
		return ErrEmptyProject
	}
	seen := make(map[string]struct{}, len(deps))
	unique := make([]string, 0, len(deps))
	for _, d := range deps {
		if d == project {
			return ErrSelfDependency
		}
		if _, ok := seen[d]; ok {
			continue
		}
		seen[d] = struct{}{}
		unique = append(unique, d)
	}
	b, err := json.Marshal(unique)
	if err != nil {
		return fmt.Errorf("dependency: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), b)
}

// Get returns the dependency list for project.
// If no record exists, a nil slice and nil error are returned.
func (m *Manager) Get(project string) ([]string, error) {
	if project == "" {
		return nil, ErrEmptyProject
	}
	b, err := m.st.Get(recordKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dependency: get: %w", err)
	}
	var deps []string
	if err := json.Unmarshal(b, &deps); err != nil {
		return nil, fmt.Errorf("dependency: unmarshal: %w", err)
	}
	return deps, nil
}

// Delete removes the dependency record for project.
func (m *Manager) Delete(project string) error {
	if project == "" {
		return ErrEmptyProject
	}
	return m.st.Delete(recordKey(project))
}
