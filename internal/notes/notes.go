// Package notes provides per-project free-text annotation storage.
// Notes are stored encrypted alongside chain data and support
// timestamped append operations as well as full replacement.
package notes

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
)

const currentVersion = 1

// Note holds a single annotation for a project.
type Note struct {
	Version   int       `json:"version"`
	Project   string    `json:"project"`
	Body      string    `json:"body"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Manager persists notes via an encrypted store.
type Manager struct {
	s *store.Store
}

// New returns a Manager backed by the given store.
func New(s *store.Store) *Manager {
	return &Manager{s: s}
}

func noteKey(project string) string {
	return fmt.Sprintf("notes:%s", project)
}

// Set replaces the note body for project, encrypting with passphrase.
func (m *Manager) Set(project, body, passphrase string) error {
	if project == "" {
		return fmt.Errorf("notes: project name must not be empty")
	}
	n := Note{
		Version:   currentVersion,
		Project:   project,
		Body:      body,
		UpdatedAt: time.Now().UTC(),
	}
	data, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("notes: marshal: %w", err)
	}
	return m.s.Put(noteKey(project), data, passphrase)
}

// Get retrieves the note for project, decrypting with passphrase.
// Returns ErrNotFound (from store) when no note exists.
func (m *Manager) Get(project, passphrase string) (Note, error) {
	data, err := m.s.Get(noteKey(project), passphrase)
	if err != nil {
		return Note{}, fmt.Errorf("notes: get %q: %w", project, err)
	}
	var n Note
	if err := json.Unmarshal(data, &n); err != nil {
		return Note{}, fmt.Errorf("notes: unmarshal: %w", err)
	}
	return n, nil
}

// Delete removes the note for project.
func (m *Manager) Delete(project string) error {
	return m.s.Delete(noteKey(project))
}
