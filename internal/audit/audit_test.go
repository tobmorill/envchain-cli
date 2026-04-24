package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain-cli/internal/audit"
)

func newTempLogger(t *testing.T) *audit.Logger {
	t.Helper()
	dir := t.TempDir()
	l, err := audit.NewLogger(filepath.Join(dir, "audit", "log.jsonl"))
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	return l
}

func TestReadAllEmpty(t *testing.T) {
	l := newTempLogger(t)
	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestRecordAndReadAll(t *testing.T) {
	l := newTempLogger(t)

	if err := l.Record(audit.EventSave, "myproject", ""); err != nil {
		t.Fatalf("Record save: %v", err)
	}
	if err := l.Record(audit.EventLoad, "myproject", ""); err != nil {
		t.Fatalf("Record load: %v", err)
	}

	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Kind != audit.EventSave {
		t.Errorf("event[0].Kind = %q, want %q", events[0].Kind, audit.EventSave)
	}
	if events[1].Kind != audit.EventLoad {
		t.Errorf("event[1].Kind = %q, want %q", events[1].Kind, audit.EventLoad)
	}
	for _, e := range events {
		if e.Project != "myproject" {
			t.Errorf("event.Project = %q, want %q", e.Project, "myproject")
		}
		if e.Timestamp.IsZero() {
			t.Error("event.Timestamp is zero")
		}
	}
}

func TestRecordWithMessage(t *testing.T) {
	l := newTempLogger(t)
	if err := l.Record(audit.EventDelete, "proj", "removed by user"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if events[0].Message != "removed by user" {
		t.Errorf("Message = %q, want %q", events[0].Message, "removed by user")
	}
}

func TestNewLoggerCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "nested", "deep", "audit.jsonl")
	_, err := audit.NewLogger(logPath)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(logPath)); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}
