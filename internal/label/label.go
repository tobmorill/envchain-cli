// Package label provides per-project free-form label management.
// Labels are arbitrary key-value string pairs attached to a project,
// useful for organising chains by team, environment, or any other
// dimension that does not warrant a dedicated data structure.
package label

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/envchain/envchain-cli/internal/store"
)

// ErrNotFound is returned when no labels exist for the requested project.
var ErrNotFound = errors.New("label: no labels found for project")

// Labels is a map of label key to label value.
type Labels map[string]string

// Manager persists and retrieves labels using a backing store.
type Manager struct {
	st *store.Store
}

// New creates a Manager backed by the store at dir.
func New(dir string) (*Manager, error) {
	st, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("label: open store: %w", err)
	}
	return &Manager{st: st}, nil
}

func labelKey(project string) string {
	return "label:" + project
}

// Set replaces all labels for project with the provided map.
// An empty map is valid and will clear all labels.
func (m *Manager) Set(project string, labels Labels) error {
	if project == "" {
		return errors.New("label: project name must not be empty")
	}
	data, err := json.Marshal(labels)
	if err != nil {
		return fmt.Errorf("label: marshal: %w", err)
	}
	return m.st.Put(labelKey(project), data)
}

// Get returns the labels for project. Returns ErrNotFound if none have
// been set.
func (m *Manager) Get(project string) (Labels, error) {
	data, err := m.st.Get(labelKey(project))
	if errors.Is(err, store.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("label: get: %w", err)
	}
	var labels Labels
	if err := json.Unmarshal(data, &labels); err != nil {
		return nil, fmt.Errorf("label: unmarshal: %w", err)
	}
	return labels, nil
}

// Delete removes all labels for project. It is not an error if no labels
// exist.
func (m *Manager) Delete(project string) error {
	if err := m.st.Delete(labelKey(project)); err != nil && !errors.Is(err, store.ErrNotFound) {
		return fmt.Errorf("label: delete: %w", err)
	}
	return nil
}
