// Package tags provides functionality for tagging chains with
// arbitrary labels, enabling grouping and filtering of environment
// variable sets across projects.
package tags

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/envchain-cli/internal/store"
)

// ErrTagNotFound is returned when a requested tag does not exist.
var ErrTagNotFound = errors.New("tag not found")

// Manager handles tag associations for chains.
type Manager struct {
	st *store.Store
}

const tagPrefix = "tags:"

// New creates a Manager backed by the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

// Set replaces the full tag list for a given chain key.
func (m *Manager) Set(chainKey string, tags []string) error {
	if chainKey == "" {
		return errors.New("chain key must not be empty")
	}
	normalized := normalizeTags(tags)
	data, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("marshal tags: %w", err)
	}
	return m.st.Put(tagPrefix+chainKey, data)
}

// Get returns the tags associated with a chain key.
// Returns ErrTagNotFound if no tags have been set.
func (m *Manager) Get(chainKey string) ([]string, error) {
	data, err := m.st.Get(tagPrefix + chainKey)
	if errors.Is(err, store.ErrNotFound) {
		return nil, ErrTagNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("unmarshal tags: %w", err)
	}
	return tags, nil
}

// Delete removes all tags for the given chain key.
func (m *Manager) Delete(chainKey string) error {
	return m.st.Delete(tagPrefix + chainKey)
}

// FindByTag returns all chain keys that carry the given tag.
func (m *Manager) FindByTag(tag string) ([]string, error) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	keys, err := m.st.Keys(tagPrefix)
	if err != nil {
		return nil, fmt.Errorf("list tag keys: %w", err)
	}
	var matches []string
	for _, k := range keys {
		chainKey := strings.TrimPrefix(k, tagPrefix)
		tags, err := m.Get(chainKey)
		if err != nil {
			continue
		}
		for _, t := range tags {
			if t == tag {
				matches = append(matches, chainKey)
				break
			}
		}
	}
	sort.Strings(matches)
	return matches, nil
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		n := strings.ToLower(strings.TrimSpace(t))
		if n == "" {
			continue
		}
		if _, ok := seen[n]; !ok {
			seen[n] = struct{}{}
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out
}
