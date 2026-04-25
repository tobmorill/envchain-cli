// Package snapshot provides point-in-time captures of environment variable
// chains, allowing users to restore previous states after edits.
package snapshot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/store"
)

const snapshotPrefix = "snapshot:"

// Snapshot holds a named point-in-time copy of a chain's raw entries.
type Snapshot struct {
	Project   string    `json:"project"`
	CreatedAt time.Time `json:"created_at"`
	Label     string    `json:"label"`
	Entries   []string  `json:"entries"`
}

// Manager handles saving and restoring snapshots.
type Manager struct {
	st *store.Store
	cm *chain.Manager
}

// New returns a Manager backed by the given store and chain manager.
func New(st *store.Store, cm *chain.Manager) *Manager {
	return &Manager{st: st, cm: cm}
}

func snapshotKey(project, label string) string {
	return fmt.Sprintf("%s%s:%s", snapshotPrefix, project, label)
}

// Save captures the current entries of a chain under the given label.
func (m *Manager) Save(project, label, passphrase string) error {
	entries, err := m.cm.Load(project, passphrase)
	if err != nil {
		return fmt.Errorf("snapshot: load chain: %w", err)
	}

	lines := make([]string, len(entries))
	for i, e := range entries {
		lines[i] = e.String()
	}

	snap := Snapshot{
		Project:   project,
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   lines,
	}

	data, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := m.st.Put(snapshotKey(project, label), data); err != nil {
		return fmt.Errorf("snapshot: store: %w", err)
	}
	return nil
}

// Get retrieves a previously saved snapshot.
func (m *Manager) Get(project, label string) (*Snapshot, error) {
	data, err := m.st.Get(snapshotKey(project, label))
	if err != nil {
		return nil, fmt.Errorf("snapshot: get: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

// Delete removes a snapshot by label.
func (m *Manager) Delete(project, label string) error {
	if err := m.st.Delete(snapshotKey(project, label)); err != nil {
		return fmt.Errorf("snapshot: delete: %w", err)
	}
	return nil
}
