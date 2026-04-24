// Package chain manages named environment variable chains (sets)
// associated with a project directory.
package chain

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/user/envchain-cli/internal/env"
	"github.com/user/envchain-cli/internal/store"
)

// ErrChainNotFound is returned when a named chain does not exist.
var ErrChainNotFound = errors.New("chain not found")

// Manager handles CRUD operations for named env chains backed by a Store.
type Manager struct {
	store *store.Store
}

// New creates a Manager using the given store.
func New(s *store.Store) *Manager {
	return &Manager{store: s}
}

// chainKey builds the store key for a project + chain name.
func chainKey(project, name string) string {
	project = filepath.ToSlash(strings.TrimSuffix(project, "/"))
	return fmt.Sprintf("%s::%s", project, name)
}

// Save encrypts and persists a slice of env entries under the given chain.
func (m *Manager) Save(project, name, passphrase string, entries []env.Entry) error {
	if name == "" {
		return errors.New("chain name must not be empty")
	}
	lines := env.ToLines(entries)
	key := chainKey(project, name)
	return m.store.Put(key, passphrase, []byte(lines))
}

// Load decrypts and returns the env entries for the given chain.
func (m *Manager) Load(project, name, passphrase string) ([]env.Entry, error) {
	if name == "" {
		return nil, errors.New("chain name must not be empty")
	}
	key := chainKey(project, name)
	data, err := m.store.Get(key, passphrase)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("%w: %s/%s", ErrChainNotFound, project, name)
		}
		return nil, err
	}
	entries, err := env.ParseAll(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing chain data: %w", err)
	}
	return entries, nil
}

// Delete removes a chain from the store.
func (m *Manager) Delete(project, name string) error {
	if name == "" {
		return errors.New("chain name must not be empty")
	}
	key := chainKey(project, name)
	return m.store.Delete(key)
}
