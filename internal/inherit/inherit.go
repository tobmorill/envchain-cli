// Package inherit provides chain inheritance resolution, allowing one
// project's environment chain to extend another's entries.
package inherit

import (
	"errors"
	"fmt"

	"github.com/user/envchain-cli/internal/store"
)

// ErrCircular is returned when a circular inheritance chain is detected.
var ErrCircular = errors.New("inherit: circular dependency detected")

// ErrSelfReference is returned when a project attempts to inherit from itself.
var ErrSelfReference = errors.New("inherit: project cannot inherit from itself")

const inheritPrefix = "inherit:"

// Manager manages project inheritance relationships.
type Manager struct {
	st *store.Store
}

// New returns a new Manager backed by st.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

// Set records that project inherits from parent.
// Returns ErrSelfReference if project == parent.
func (m *Manager) Set(project, parent string) error {
	if project == "" || parent == "" {
		return errors.New("inherit: project and parent must not be empty")
	}
	if project == parent {
		return ErrSelfReference
	}
	return m.st.Put(inheritPrefix+project, []byte(parent))
}

// Get returns the direct parent of project, or ("", nil) if none is set.
func (m *Manager) Get(project string) (string, error) {
	v, err := m.st.Get(inheritPrefix + project)
	if errors.Is(err, store.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(v), nil
}

// Delete removes the inheritance relationship for project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(inheritPrefix + project)
}

// Chain returns the full ordered ancestry of project, starting with the
// immediate parent and ending at the root. Returns ErrCircular if a cycle
// is detected.
func (m *Manager) Chain(project string) ([]string, error) {
	seen := map[string]bool{project: true}
	var chain []string
	current := project
	for {
		parent, err := m.Get(current)
		if err != nil {
			return nil, fmt.Errorf("inherit: resolving %q: %w", current, err)
		}
		if parent == "" {
			break
		}
		if seen[parent] {
			return nil, ErrCircular
		}
		seen[parent] = true
		chain = append(chain, parent)
		current = parent
	}
	return chain, nil
}
