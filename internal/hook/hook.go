// Package hook manages shell hooks that run before/after envchain loads
// a chain, allowing users to execute custom scripts at defined lifecycle points.
package hook

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Phase represents when a hook fires relative to chain loading.
type Phase string

const (
	PhasePre  Phase = "pre"
	PhasePost Phase = "post"
)

// Hook holds a single shell command attached to a project and phase.
type Hook struct {
	Project string `json:"project"`
	Phase   Phase  `json:"phase"`
	Command string `json:"command"`
}

// Manager stores and retrieves hooks for projects.
type Manager struct {
	dir string
}

// New returns a Manager rooted at dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) hookPath(project string, phase Phase) string {
	return filepath.Join(m.dir, fmt.Sprintf("%s.%s.json", project, phase))
}

// Set persists a hook for the given project and phase, overwriting any
// existing hook.
func (m *Manager) Set(h Hook) error {
	if h.Project == "" {
		return errors.New("hook: project name must not be empty")
	}
	if h.Phase != PhasePre && h.Phase != PhasePost {
		return fmt.Errorf("hook: unknown phase %q", h.Phase)
	}
	if h.Command == "" {
		return errors.New("hook: command must not be empty")
	}
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return os.WriteFile(m.hookPath(h.Project, h.Phase), data, 0o600)
}

// Get retrieves the hook for the given project and phase.
// Returns (Hook{}, false, nil) when no hook is registered.
func (m *Manager) Get(project string, phase Phase) (Hook, bool, error) {
	data, err := os.ReadFile(m.hookPath(project, phase))
	if errors.Is(err, os.ErrNotExist) {
		return Hook{}, false, nil
	}
	if err != nil {
		return Hook{}, false, err
	}
	var h Hook
	if err := json.Unmarshal(data, &h); err != nil {
		return Hook{}, false, err
	}
	return h, true, nil
}

// Delete removes the hook for the given project and phase.
// Returns nil if no hook exists.
func (m *Manager) Delete(project string, phase Phase) error {
	err := os.Remove(m.hookPath(project, phase))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// List returns all hooks stored for the given project.
func (m *Manager) List(project string) ([]Hook, error) {
	var hooks []Hook
	for _, phase := range []Phase{PhasePre, PhasePost} {
		h, ok, err := m.Get(project, phase)
		if err != nil {
			return nil, err
		}
		if ok {
			hooks = append(hooks, h)
		}
	}
	return hooks, nil
}
