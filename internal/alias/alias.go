// Package alias provides management of named aliases for environment chains,
// allowing users to reference chains by short memorable names.
package alias

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ErrNotFound is returned when an alias does not exist.
var ErrNotFound = errors.New("alias not found")

// ErrInvalidName is returned when an alias name contains invalid characters.
var ErrInvalidName = errors.New("invalid alias name")

var validAlias = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// Manager manages alias-to-chain mappings persisted on disk.
type Manager struct {
	path string
}

type aliasFile struct {
	Aliases map[string]string `json:"aliases"`
}

// New returns a Manager that stores aliases at the given file path.
func New(path string) *Manager {
	return &Manager{path: path}
}

// Set creates or updates an alias pointing to chainName.
func (m *Manager) Set(alias, chainName string) error {
	if !validAlias.MatchString(alias) {
		return fmt.Errorf("%w: %q", ErrInvalidName, alias)
	}
	af, err := m.load()
	if err != nil {
		return err
	}
	af.Aliases[strings.ToLower(alias)] = chainName
	return m.save(af)
}

// Get returns the chain name for the given alias.
func (m *Manager) Get(alias string) (string, error) {
	af, err := m.load()
	if err != nil {
		return "", err
	}
	chain, ok := af.Aliases[strings.ToLower(alias)]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrNotFound, alias)
	}
	return chain, nil
}

// Delete removes an alias. Returns ErrNotFound if it does not exist.
func (m *Manager) Delete(alias string) error {
	af, err := m.load()
	if err != nil {
		return err
	}
	key := strings.ToLower(alias)
	if _, ok := af.Aliases[key]; !ok {
		return fmt.Errorf("%w: %q", ErrNotFound, alias)
	}
	delete(af.Aliases, key)
	return m.save(af)
}

// List returns all alias->chain mappings.
func (m *Manager) List() (map[string]string, error) {
	af, err := m.load()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(af.Aliases))
	for k, v := range af.Aliases {
		out[k] = v
	}
	return out, nil
}

func (m *Manager) load() (aliasFile, error) {
	af := aliasFile{Aliases: make(map[string]string)}
	data, err := os.ReadFile(m.path)
	if errors.Is(err, os.ErrNotExist) {
		return af, nil
	}
	if err != nil {
		return af, err
	}
	if err := json.Unmarshal(data, &af); err != nil {
		return af, err
	}
	if af.Aliases == nil {
		af.Aliases = make(map[string]string)
	}
	return af, nil
}

func (m *Manager) save(af aliasFile) error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(af, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o600)
}
