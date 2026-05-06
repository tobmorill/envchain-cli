// Package checksum provides integrity verification for environment variable chains.
// It computes and stores a deterministic hash of a chain's entries so that
// out-of-band modifications can be detected before a chain is loaded.
package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/envchain/envchain-cli/internal/env"
	"github.com/envchain/envchain-cli/internal/store"
)

// ErrMismatch is returned when the stored checksum does not match the computed one.
var ErrMismatch = errors.New("checksum: integrity check failed")

const keyPrefix = "checksum:"

// Manager stores and verifies checksums for named chains.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

// Save computes the SHA-256 checksum of entries and persists it under project.
func (m *Manager) Save(project string, entries []env.Entry) error {
	if project == "" {
		return errors.New("checksum: project name must not be empty")
	}
	sum, err := compute(entries)
	if err != nil {
		return fmt.Errorf("checksum: compute: %w", err)
	}
	return m.st.Put(keyPrefix+project, []byte(sum))
}

// Verify recomputes the checksum of entries and compares it to the stored value.
// Returns ErrMismatch when the values differ and ErrNotFound when no checksum
// has been saved yet.
func (m *Manager) Verify(project string, entries []env.Entry) error {
	if project == "" {
		return errors.New("checksum: project name must not be empty")
	}
	stored, err := m.st.Get(keyPrefix + project)
	if err != nil {
		return err
	}
	got, err := compute(entries)
	if err != nil {
		return fmt.Errorf("checksum: compute: %w", err)
	}
	if string(stored) != got {
		return ErrMismatch
	}
	return nil
}

// Delete removes the stored checksum for project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(keyPrefix + project)
}

// compute returns a stable hex-encoded SHA-256 digest of entries.
// Entries are sorted by key before hashing to ensure determinism.
func compute(entries []env.Entry) (string, error) {
	sorted := make([]env.Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	b, err := json.Marshal(sorted)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
