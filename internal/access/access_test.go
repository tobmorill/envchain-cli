package access_test

import (
	"testing"
	"time"

	"github.com/envchain/envchain-cli/internal/access"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *access.Manager {
	t.Helper()
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return access.New(st)
}

func TestTouchCreatesRecord(t *testing.T) {
	m := newTempManager(t)
	if err := m.Touch("myproject"); err != nil {
		t.Fatalf("Touch: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Count != 1 {
		t.Errorf("Count = %d, want 1", rec.Count)
	}
	if rec.Project != "myproject" {
		t.Errorf("Project = %q, want %q", rec.Project, "myproject")
	}
	if rec.FirstUsed.IsZero() {
		t.Error("FirstUsed should not be zero")
	}
}

func TestTouchIncrementsCount(t *testing.T) {
	m := newTempManager(t)
	for i := 0; i < 5; i++ {
		if err := m.Touch("proj"); err != nil {
			t.Fatalf("Touch #%d: %v", i, err)
		}
	}
	rec, err := m.Get("proj")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Count != 5 {
		t.Errorf("Count = %d, want 5", rec.Count)
	}
}

func TestTouchUpdatesLastUsed(t *testing.T) {
	m := newTempManager(t)
	_ = m.Touch("proj")
	before := time.Now().UTC()
	_ = m.Touch("proj")
	rec, _ := m.Get("proj")
	if rec.LastUsed.Before(before.Add(-time.Second)) {
		t.Errorf("LastUsed %v is too old", rec.LastUsed)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResetClearsRecord(t *testing.T) {
	m := newTempManager(t)
	_ = m.Touch("proj")
	if err := m.Reset("proj"); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	_, err := m.Get("proj")
	if err == nil {
		t.Fatal("expected error after reset, got nil")
	}
}

func TestTouchEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Touch(""); err == nil {
		t.Fatal("expected error for empty project, got nil")
	}
}
