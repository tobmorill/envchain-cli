// Package blame tracks which process or user last modified a project's
// environment chain, providing a lightweight audit trail at the chain level.
package blame

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/your-org/envchain-cli/internal/store"
)

// Record holds metadata about the last modification to a chain.
type Record struct {
	Project   string    `json:"project"`
	User      string    `json:"user"`
	Hostname  string    `json:"hostname"`
	ChangedAt time.Time `json:"changed_at"`
	Note      string    `json:"note,omitempty"`
}

// Manager persists blame records using the key-value store.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func blameKey(project string) string {
	return fmt.Sprintf("blame:%s", project)
}

// Touch records the current OS user and hostname as the last modifier of
// the given project chain. An optional human-readable note may be supplied.
func (m *Manager) Touch(project, note string) error {
	if project == "" {
		return fmt.Errorf("blame: project name must not be empty")
	}

	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME") // Windows fallback
	}
	if user == "" {
		user = "unknown"
	}

	hostname, _ := os.Hostname()

	rec := Record{
		Project:   project,
		User:      user,
		Hostname:  hostname,
		ChangedAt: time.Now().UTC(),
		Note:      note,
	}

	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("blame: marshal: %w", err)
	}

	return m.st.Put(blameKey(project), data)
}

// Get returns the blame record for the given project, or an error if none
// exists.
func (m *Manager) Get(project string) (Record, error) {
	data, err := m.st.Get(blameKey(project))
	if err != nil {
		return Record{}, fmt.Errorf("blame: get %q: %w", project, err)
	}

	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("blame: unmarshal: %w", err)
	}

	return rec, nil
}

// Delete removes the blame record for the given project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(blameKey(project))
}
