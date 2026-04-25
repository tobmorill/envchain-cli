package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/envchain-cli/internal/history"
)

func newTempManager(t *testing.T) *history.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return history.New(dir)
}

func TestRecordAndReadAll(t *testing.T) {
	m := newTempManager(t)
	if err := m.Record("myproject", "load"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := m.Record("myproject", "export"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := m.ReadAll("myproject")
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Action != "load" {
		t.Errorf("expected action 'load', got %q", entries[0].Action)
	}
	if entries[1].Action != "export" {
		t.Errorf("expected action 'export', got %q", entries[1].Action)
	}
}

func TestReadAllNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.ReadAll("ghost")
	if err != history.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRecordTimestamp(t *testing.T) {
	m := newTempManager(t)
	before := time.Now().UTC().Add(-time.Second)
	if err := m.Record("ts-project", "unlock"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	after := time.Now().UTC().Add(time.Second)
	entries, err := m.ReadAll("ts-project")
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if entries[0].AccessedAt.Before(before) || entries[0].AccessedAt.After(after) {
		t.Errorf("timestamp %v outside expected range", entries[0].AccessedAt)
	}
}

func TestClear(t *testing.T) {
	m := newTempManager(t)
	if err := m.Record("clearme", "load"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := m.Clear("clearme"); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	_, err := m.ReadAll("clearme")
	if err != history.ErrNotFound {
		t.Fatalf("expected ErrNotFound after clear, got %v", err)
	}
}

func TestClearNonExistent(t *testing.T) {
	m := newTempManager(t)
	if err := m.Clear("nope"); err != nil {
		t.Errorf("Clear of non-existent should not error, got %v", err)
	}
}
