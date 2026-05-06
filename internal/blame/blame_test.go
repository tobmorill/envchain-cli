package blame_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/envchain-cli/internal/blame"
	"github.com/your-org/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *blame.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "blame.db")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return blame.New(st)
}

func TestTouchAndGet(t *testing.T) {
	m := newTempManager(t)

	if err := m.Touch("myproject", "initial setup"); err != nil {
		t.Fatalf("Touch: %v", err)
	}

	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if rec.Project != "myproject" {
		t.Errorf("project = %q, want %q", rec.Project, "myproject")
	}
	if rec.Note != "initial setup" {
		t.Errorf("note = %q, want %q", rec.Note, "initial setup")
	}
	if rec.ChangedAt.IsZero() {
		t.Error("ChangedAt should not be zero")
	}
	if rec.User == "" {
		t.Error("User should not be empty")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)

	_, err := m.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing record, got nil")
	}
}

func TestTouchEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)

	if err := m.Touch("", ""); err == nil {
		t.Fatal("expected error for empty project, got nil")
	}
}

func TestTouchUpdatesTimestamp(t *testing.T) {
	m := newTempManager(t)

	if err := m.Touch("proj", ""); err != nil {
		t.Fatalf("first Touch: %v", err)
	}

	rec1, _ := m.Get("proj")
	time.Sleep(2 * time.Millisecond)

	if err := m.Touch("proj", ""); err != nil {
		t.Fatalf("second Touch: %v", err)
	}

	rec2, _ := m.Get("proj")
	if !rec2.ChangedAt.After(rec1.ChangedAt) {
		t.Error("second touch should have a later timestamp")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)

	_ = m.Touch("proj", "")
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := m.Get("proj")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestTouchUsesEnvUser(t *testing.T) {
	m := newTempManager(t)

	t.Setenv("USER", "testuser")
	_ = m.Touch("proj", "")

	rec, err := m.Get("proj")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.User != "testuser" {
		t.Errorf("user = %q, want %q", rec.User, "testuser")
	}
	_ = os.Getenv("USER") // silence unused import
}
