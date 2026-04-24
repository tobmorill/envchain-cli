// Package audit provides a simple append-only audit log for envchain
// operations such as chain creation, deletion, and access.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventLoad   EventKind = "load"
	EventSave   EventKind = "save"
	EventDelete EventKind = "delete"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      EventKind `json:"kind"`
	Project   string    `json:"project"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit events to a JSONL file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the file at path.
// The parent directory is created if it does not exist.
func NewLogger(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("audit: create directory: %w", err)
	}
	return &Logger{path: path}, nil
}

// Record appends an event to the log file.
func (l *Logger) Record(kind EventKind, project, message string) error {
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	event := Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Project:   project,
		Message:   message,
	}
	enc := json.NewEncoder(f)
	if err := enc.Encode(event); err != nil {
		return fmt.Errorf("audit: write event: %w", err)
	}
	return nil
}

// ReadAll reads and returns all events from the log file.
// Returns an empty slice if the file does not exist.
func (l *Logger) ReadAll() ([]Event, error) {
	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return []Event{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	var events []Event
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Event
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode event: %w", err)
		}
		events = append(events, e)
	}
	return events, nil
}
