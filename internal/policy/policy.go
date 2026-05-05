// Package policy enforces per-project access policies, restricting which
// environment variable keys may be loaded for a given project or context.
package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// ErrNotFound is returned when no policy exists for a project.
var ErrNotFound = errors.New("policy: not found")

// Rule describes what keys are allowed or denied for a project.
type Rule struct {
	// AllowKeys is a list of exact key names that are permitted.
	// If empty, all keys are allowed (subject to DenyKeys).
	AllowKeys []string `json:"allow_keys,omitempty"`
	// DenyKeys is a list of exact key names that are always blocked.
	DenyKeys []string `json:"deny_keys,omitempty"`
	// AllowPattern is an optional regex; only matching keys are permitted.
	AllowPattern string `json:"allow_pattern,omitempty"`
}

// Manager persists and evaluates project policies.
type Manager struct {
	dir string
}

// New returns a Manager rooted at dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) path(project string) string {
	return filepath.Join(m.dir, project+".json")
}

// Set stores a Rule for the given project, creating the directory if needed.
func (m *Manager) Set(project string, r Rule) error {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return fmt.Errorf("policy: mkdir: %w", err)
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("policy: marshal: %w", err)
	}
	return os.WriteFile(m.path(project), b, 0o600)
}

// Get retrieves the Rule for the given project.
func (m *Manager) Get(project string) (Rule, error) {
	b, err := os.ReadFile(m.path(project))
	if errors.Is(err, os.ErrNotExist) {
		return Rule{}, ErrNotFound
	}
	if err != nil {
		return Rule{}, fmt.Errorf("policy: read: %w", err)
	}
	var r Rule
	if err := json.Unmarshal(b, &r); err != nil {
		return Rule{}, fmt.Errorf("policy: unmarshal: %w", err)
	}
	return r, nil
}

// Delete removes the policy for the given project.
func (m *Manager) Delete(project string) error {
	err := os.Remove(m.path(project))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	return err
}

// Allowed reports whether key is permitted under rule.
func Allowed(r Rule, key string) (bool, error) {
	for _, d := range r.DenyKeys {
		if d == key {
			return false, nil
		}
	}
	if r.AllowPattern != "" {
		re, err := regexp.Compile(r.AllowPattern)
		if err != nil {
			return false, fmt.Errorf("policy: invalid pattern %q: %w", r.AllowPattern, err)
		}
		if !re.MatchString(key) {
			return false, nil
		}
	}
	if len(r.AllowKeys) > 0 {
		for _, a := range r.AllowKeys {
			if a == key {
				return true, nil
			}
		}
		return false, nil
	}
	return true, nil
}
