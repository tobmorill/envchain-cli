// Package embargo provides time-window based access restrictions for project
// environment chains. An embargo prevents a chain from being loaded outside
// of a configured daily time window.
package embargo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/user/envchain-cli/internal/store"
)

// ErrEmbargoActive is returned when access is attempted outside the allowed window.
var ErrEmbargoActive = errors.New("embargo: access denied outside permitted time window")

// Window defines the daily start and end time (UTC) during which access is permitted.
type Window struct {
	StartHour int `json:"start_hour"`
	StartMin  int `json:"start_min"`
	EndHour   int `json:"end_hour"`
	EndMin    int `json:"end_min"`
}

// Manager manages embargo windows backed by a key-value store.
type Manager struct {
	st *store.Store
}

func recordKey(project string) string {
	return "embargo:" + project
}

// New returns a Manager using the provided store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

// Set persists an embargo window for the given project.
func (m *Manager) Set(project string, w Window) error {
	if project == "" {
		return errors.New("embargo: project name must not be empty")
	}
	if w.StartHour < 0 || w.StartHour > 23 || w.EndHour < 0 || w.EndHour > 23 {
		return errors.New("embargo: hour must be between 0 and 23")
	}
	if w.StartMin < 0 || w.StartMin > 59 || w.EndMin < 0 || w.EndMin > 59 {
		return errors.New("embargo: minute must be between 0 and 59")
	}
	b, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("embargo: marshal: %w", err)
	}
	return m.st.Put(recordKey(project), b)
}

// Get retrieves the embargo window for the given project.
func (m *Manager) Get(project string) (Window, error) {
	b, err := m.st.Get(recordKey(project))
	if err != nil {
		return Window{}, fmt.Errorf("embargo: get: %w", err)
	}
	var w Window
	if err := json.Unmarshal(b, &w); err != nil {
		return Window{}, fmt.Errorf("embargo: unmarshal: %w", err)
	}
	return w, nil
}

// Delete removes the embargo window for the given project.
func (m *Manager) Delete(project string) error {
	return m.st.Delete(recordKey(project))
}

// Check returns ErrEmbargoActive if the current UTC time falls outside the
// window stored for project. If no window is set, access is always permitted.
func (m *Manager) Check(project string) error {
	w, err := m.Get(project)
	if err != nil {
		// No window configured — access permitted.
		return nil
	}
	now := time.Now().UTC()
	current := now.Hour()*60 + now.Minute()
	start := w.StartHour*60 + w.StartMin
	end := w.EndHour*60 + w.EndMin
	if start <= end {
		if current >= start && current < end {
			return nil
		}
	} else {
		// Window spans midnight.
		if current >= start || current < end {
			return nil
		}
	}
	return ErrEmbargoActive
}
