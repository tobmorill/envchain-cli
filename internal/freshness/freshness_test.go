package freshness_test

import (
	"os"
	"testing"
	"time"

	"github.com/envchain/envchain-cli/internal/freshness"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *freshness.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "freshness-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return freshness.New(st)
}

func TestTouchAndGet(t *testing.T) {
	m := newTempManager(t)
	before := time.Now().UTC().Add(-time.Second)
	if err := m.Touch("myproject"); err != nil {
		t.Fatalf("Touch: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Project != "myproject" {
		t.Errorf("project = %q, want %q", rec.Project, "myproject")
	}
	if rec.TouchedAt.Before(before) {
		t.Errorf("TouchedAt %v is before test start %v", rec.TouchedAt, before)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err != freshness.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestIsStale(t *testing.T) {
	m := newTempManager(t)
	if err := m.Touch("proj"); err != nil {
		t.Fatalf("Touch: %v", err)
	}
	rec, _ := m.Get("proj")
	if rec.IsStale(time.Hour) {
		t.Error("expected record to be fresh within an hour threshold")
	}
	if !rec.IsStale(time.Nanosecond) {
		t.Error("expected record to be stale with nanosecond threshold")
	}
}

func TestTouchUpdatesTimestamp(t *testing.T) {
	m := newTempManager(t)
	if err := m.Touch("proj"); err != nil {
		t.Fatalf("first Touch: %v", err)
	}
	rec1, _ := m.Get("proj")
	time.Sleep(2 * time.Millisecond)
	if err := m.Touch("proj"); err != nil {
		t.Fatalf("second Touch: %v", err)
	}
	rec2, _ := m.Get("proj")
	if !rec2.TouchedAt.After(rec1.TouchedAt) {
		t.Errorf("expected second touch %v to be after first %v", rec2.TouchedAt, rec1.TouchedAt)
	}
}

func TestDeleteRemovesRecord(t *testing.T) {
	m := newTempManager(t)
	if err := m.Touch("proj"); err != nil {
		t.Fatalf("Touch: %v", err)
	}
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if err != freshness.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNoop(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("nonexistent"); err != nil {
		t.Errorf("Delete of nonexistent should not error, got %v", err)
	}
}

func TestTouchEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Touch(""); err == nil {
		t.Error("expected error for empty project name")
	}
}
