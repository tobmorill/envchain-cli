// Package rotate provides passphrase rotation for envchain chains.
// It re-encrypts all stored environment variable sets under a new passphrase
// without exposing plaintext values beyond the duration of the operation.
package rotate

import (
	"fmt"

	"github.com/your-org/envchain-cli/internal/chain"
	"github.com/your-org/envchain-cli/internal/env"
	"github.com/your-org/envchain-cli/internal/store"
)

// Manager handles passphrase rotation for a given store.
type Manager struct {
	st *store.Store
}

// New returns a new rotation Manager backed by the provided store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

// Rotate re-encrypts the named chain from oldPassphrase to newPassphrase.
// It returns an error if the chain cannot be loaded or re-saved.
func (m *Manager) Rotate(name, oldPassphrase, newPassphrase string) error {
	if oldPassphrase == newPassphrase {
		return fmt.Errorf("rotate: new passphrase must differ from old passphrase")
	}

	cm := chain.New(m.st)

	entries, err := cm.Load(name, oldPassphrase)
	if err != nil {
		return fmt.Errorf("rotate: load chain %q: %w", name, err)
	}

	if err := cm.Save(name, newPassphrase, entries); err != nil {
		return fmt.Errorf("rotate: save chain %q with new passphrase: %w", name, err)
	}

	return nil
}

// RotateAll rotates every chain whose name is listed in names.
// It stops and returns on the first error, leaving already-rotated chains
// under the new passphrase and unprocessed chains under the old one.
func (m *Manager) RotateAll(names []string, oldPassphrase, newPassphrase string) error {
	for _, name := range names {
		if err := m.Rotate(name, oldPassphrase, newPassphrase); err != nil {
			return err
		}
	}
	return nil
}

// Preview loads the chain with oldPassphrase and returns the entries without
// modifying the store. Useful for dry-run validation before rotation.
func (m *Manager) Preview(name, oldPassphrase string) ([]env.Entry, error) {
	cm := chain.New(m.st)
	entries, err := cm.Load(name, oldPassphrase)
	if err != nil {
		return nil, fmt.Errorf("rotate: preview chain %q: %w", name, err)
	}
	return entries, nil
}
