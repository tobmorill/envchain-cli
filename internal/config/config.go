// Package config manages persistent CLI configuration for envchain.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	defaultFileName = "config.json"
	defaultDirName  = ".envchain"
)

// Config holds user-level configuration for the envchain CLI.
type Config struct {
	DefaultShell    string `json:"default_shell,omitempty"`
	StorePath       string `json:"store_path,omitempty"`
	PassphraseHint  string `json:"passphrase_hint,omitempty"`
}

// Manager handles reading and writing the config file.
type Manager struct {
	path string
}

// NewManager returns a Manager rooted at the given directory.
// If dir is empty, the user's home directory is used.
func NewManager(dir string) (*Manager, error) {
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(home, defaultDirName)
	}
	return &Manager{path: filepath.Join(dir, defaultFileName)}, nil
}

// Load reads the config from disk. If the file does not exist, a zero-value
// Config is returned without error.
func (m *Manager) Load() (Config, error) {
	data, err := os.ReadFile(m.path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Save writes cfg to disk, creating the directory if necessary.
func (m *Manager) Save(cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o600)
}

// Path returns the resolved path to the config file.
func (m *Manager) Path() string {
	return m.path
}
