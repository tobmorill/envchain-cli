// Package template provides functionality for rendering environment variable
// sets as named templates that can be reused across multiple projects.
package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/envchain/envchain-cli/internal/store"
)

var validTemplateName = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// ErrNotFound is returned when a named template does not exist.
var ErrNotFound = errors.New("template: not found")

// ErrInvalidName is returned when a template name contains invalid characters.
var ErrInvalidName = errors.New("template: invalid name")

// Template holds a named collection of environment variable keys.
type Template struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// Manager manages stored templates using the underlying key-value store.
type Manager struct {
	st *store.Store
}

// New creates a Manager backed by the provided store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func templateKey(name string) string {
	return fmt.Sprintf("template::%s", name)
}

// Save persists a template. Returns ErrInvalidName if the name is not valid.
func (m *Manager) Save(t Template) error {
	if !validTemplateName.MatchString(t.Name) {
		return ErrInvalidName
	}
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("template: marshal: %w", err)
	}
	return m.st.Put(templateKey(t.Name), data)
}

// Load retrieves a template by name. Returns ErrNotFound if absent.
func (m *Manager) Load(name string) (Template, error) {
	data, err := m.st.Get(templateKey(name))
	if errors.Is(err, store.ErrNotFound) {
		return Template{}, ErrNotFound
	}
	if err != nil {
		return Template{}, fmt.Errorf("template: get: %w", err)
	}
	var t Template
	if err := json.Unmarshal(data, &t); err != nil {
		return Template{}, fmt.Errorf("template: unmarshal: %w", err)
	}
	return t, nil
}

// Delete removes a template by name. Returns ErrNotFound if absent.
func (m *Manager) Delete(name string) error {
	err := m.st.Delete(templateKey(name))
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// IsValidName reports whether name is a valid template identifier.
func IsValidName(name string) bool {
	return validTemplateName.MatchString(name)
}
