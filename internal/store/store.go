package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/envchain-cli/internal/crypto"
)

const storeFileName = "envchain.enc"

// EnvSet represents a named set of environment variables.
type EnvSet struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// Store manages encrypted storage of environment variable sets.
type Store struct {
	path string
}

// New creates a Store rooted at the given directory.
func New(dir string) *Store {
	return &Store{path: filepath.Join(dir, storeFileName)}
}

// load decrypts and deserialises all env sets from disk.
func (s *Store) load(passphrase string) (map[string]EnvSet, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]EnvSet{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("store: read file: %w", err)
	}

	plain, err := crypto.Decrypt(data, passphrase)
	if err != nil {
		return nil, fmt.Errorf("store: decrypt: %w", err)
	}

	var sets map[string]EnvSet
	if err := json.Unmarshal(plain, &sets); err != nil {
		return nil, fmt.Errorf("store: unmarshal: %w", err)
	}
	return sets, nil
}

// save serialises and encrypts all env sets to disk.
func (s *Store) save(sets map[string]EnvSet, passphrase string) error {
	plain, err := json.Marshal(sets)
	if err != nil {
		return fmt.Errorf("store: marshal: %w", err)
	}

	cipher, err := crypto.Encrypt(plain, passphrase)
	if err != nil {
		return fmt.Errorf("store: encrypt: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("store: mkdir: %w", err)
	}
	return os.WriteFile(s.path, cipher, 0o600)
}

// Put adds or replaces an env set.
func (s *Store) Put(set EnvSet, passphrase string) error {
	sets, err := s.load(passphrase)
	if err != nil {
		return err
	}
	sets[set.Name] = set
	return s.save(sets, passphrase)
}

// Get retrieves an env set by name.
func (s *Store) Get(name, passphrase string) (EnvSet, error) {
	sets, err := s.load(passphrase)
	if err != nil {
		return EnvSet{}, err
	}
	set, ok := sets[name]
	if !ok {
		return EnvSet{}, fmt.Errorf("store: env set %q not found", name)
	}
	return set, nil
}

// Delete removes an env set by name.
func (s *Store) Delete(name, passphrase string) error {
	sets, err := s.load(passphrase)
	if err != nil {
		return err
	}
	if _, ok := sets[name]; !ok {
		return fmt.Errorf("store: env set %q not found", name)
	}
	delete(sets, name)
	return s.save(sets, passphrase)
}

// List returns the names of all stored env sets.
func (s *Store) List(passphrase string) ([]string, error) {
	sets, err := s.load(passphrase)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(sets))
	for name := range sets {
		names = append(names, name)
	}
	return names, nil
}
