// Package namespace provides grouping of projects under named namespaces,
// allowing teams to organise related chains without name collisions.
package namespace

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/envchain/envchain-cli/internal/store"
)

var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ErrInvalidName is returned when a namespace name contains illegal characters.
var ErrInvalidName = errors.New("namespace: name must match [a-zA-Z0-9_-]+")

// ErrNotFound is returned when no namespace record exists for the given name.
var ErrNotFound = errors.New("namespace: not found")

// Record holds the list of project names belonging to a namespace.
type Record struct {
	Projects []string `json:"projects"`
}

// Manager persists namespace records via the key-value store.
type Manager struct {
	st *store.Store
}

// New returns a Manager backed by st.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

func recordKey(name string) string {
	return fmt.Sprintf("namespace:%s", strings.ToLower(name))
}

// Set overwrites the project list for the given namespace.
func (m *Manager) Set(name string, projects []string) error {
	if !validName.MatchString(name) {
		return ErrInvalidName
	}
	norm := make([]string, 0, len(projects))
	seen := make(map[string]struct{})
	for _, p := range projects {
		key := strings.ToLower(p)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		norm = append(norm, p)
	}
	sort.Strings(norm)
	rec := Record{Projects: norm}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return m.st.Put(recordKey(name), data)
}

// Get returns the Record for name, or ErrNotFound.
func (m *Manager) Get(name string) (Record, error) {
	data, err := m.st.Get(recordKey(name))
	if errors.Is(err, store.ErrNotFound) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

// Delete removes the namespace record for name.
func (m *Manager) Delete(name string) error {
	err := m.st.Delete(recordKey(name))
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
