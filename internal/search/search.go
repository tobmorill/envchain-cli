// Package search provides functionality for finding chains and environment
// variable keys across all stored projects.
package search

import (
	"strings"

	"github.com/user/envchain-cli/internal/chain"
)

// Result represents a single search match.
type Result struct {
	// Project is the name of the chain/project that contains the match.
	Project string
	// Key is the environment variable key that matched, if any.
	Key string
}

// Manager performs searches over a chain store.
type Manager struct {
	chains *chain.Manager
}

// New returns a Manager backed by the given chain manager.
func New(cm *chain.Manager) *Manager {
	return &Manager{chains: cm}
}

// FindProjects returns all project names that contain the given substring
// (case-insensitive). An empty query returns all project names.
func (m *Manager) FindProjects(query string, names []string) []Result {
	q := strings.ToLower(query)
	var results []Result
	for _, name := range names {
		if q == "" || strings.Contains(strings.ToLower(name), q) {
			results = append(results, Result{Project: name})
		}
	}
	return results
}

// FindKeys searches for environment variable keys matching the query across
// the provided project names. passphrase is used to decrypt each chain.
// Projects that cannot be decrypted are silently skipped.
func (m *Manager) FindKeys(query, passphrase string, names []string) ([]Result, error) {
	q := strings.ToLower(query)
	var results []Result
	for _, name := range names {
		entries, err := m.chains.Load(name, passphrase)
		if err != nil {
			// Skip chains we cannot decrypt (wrong passphrase, corrupt data).
			continue
		}
		for _, e := range entries {
			if q == "" || strings.Contains(strings.ToLower(e.Key), q) {
				results = append(results, Result{Project: name, Key: e.Key})
			}
		}
	}
	return results, nil
}
