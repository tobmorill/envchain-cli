package signal_test

import (
	"errors"
	"testing"

	"github.com/envchain/envchain-cli/internal/signal"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *signal.Manager {
	t.Helper()
	dir := t.TempDir()
	m, err := signal.New(dir)
	if err != nil {
		t.Fatalf("signal.New: %v", err)
	}
	return m
}

func TestRaiseAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Raise("myproject", "disk usage high", signal.LevelWarn); err != nil {
		t.Fatalf("Raise: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Project != "myproject" {
		t.Errorf("project: got %q, want %q", rec.Project, "myproject")
	}
	if rec.Message != "disk usage high" {
		t.Errorf("message: got %q", rec.Message)
	}
	if rec.Level != signal.LevelWarn {
		t.Errorf("level: got %q, want warn", rec.Level)
	}
	if rec.Acknowledged {
		t.Error("expected acknowledged=false on new signal")
	}
	if rec.RaisedAt.IsZero() {
		t.Error("expected RaisedAt to be set")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err == nil {
		t.Fatal("expected error for missing signal")
	}
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAcknowledge(t *testing.T) {
	m := newTempManager(t)
	_ = m.Raise("proj", "something broke", signal.LevelError)
	if err := m.Acknowledge("proj"); err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}
	rec, _ := m.Get("proj")
	if !rec.Acknowledged {
		t.Error("expected acknowledged=true after Acknowledge")
	}
}

func TestRaiseOverwrites(t *testing.T) {
	m := newTempManager(t)
	_ = m.Raise("proj", "first", signal.LevelInfo)
	_ = m.Raise("proj", "second", signal.LevelError)
	rec, _ := m.Get("proj")
	if rec.Message != "second" {
		t.Errorf("expected second message, got %q", rec.Message)
	}
	if rec.Level != signal.LevelError {
		t.Errorf("expected error level, got %q", rec.Level)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Raise("proj", "msg", signal.LevelInfo)
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestRaiseEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Raise("", "msg", signal.LevelInfo); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestRaiseEmptyMessageReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Raise("proj", "", signal.LevelInfo); err == nil {
		t.Error("expected error for empty message")
	}
}
