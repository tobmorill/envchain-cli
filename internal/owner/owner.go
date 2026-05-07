// Package owner tracks the declared owner (team or individual) for each
// project's environment chain. Ownership metadata is stored in the local
// key-value store alongside other project attributes.
package owner

import (
	"errors"
	"fmt"
	"strings"

	"github.com/envchain/envchain-cli/internal/store"
)

// Record holds ownership information for a project.
type Record struct {
	Project string `json:"project"`
	Owner   string `json:"owner"`
}

// Manager persists and retrieves owner records.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(project string) string {
	return fmt.Sprintf("owner:%s", project)
}

// Set stores the owner for the given project.
// Both project and owner must be non-empty strings.
func (m *Manager) Set(project, owner string) error {
	project = strings.TrimSpace(project)
	owner = strings.TrimSpace(owner)
	if project == "" {
		return errors.New("owner: project name must not be empty")
	}
	if owner == "" {
		return errors.New("owner: owner must not be empty")
	}
	return m.st.Put(recordKey(project), []byte(owner))
}

// Get returns the owner record for the given project.
// It returns (zero Record, store.ErrNotFound) when no owner is set.
func (m *Manager) Get(project string) (Record, error) {
	project = strings.TrimSpace(project)
	if project == "" {
		return Record{}, errors.New("owner: project name must not be empty")
	}
	val, err := m.st.Get(recordKey(project))
	if err != nil {
		return Record{}, err
	}
	return Record{Project: project, Owner: string(val)}, nil
}

// Delete removes the owner record for the given project.
// Deleting a non-existent record is a no-op.
func (m *Manager) Delete(project string) error {
	project = strings.TrimSpace(project)
	if project == "" {
		return errors.New("owner: project name must not be empty")
	}
	return m.st.Delete(recordKey(project))
}
