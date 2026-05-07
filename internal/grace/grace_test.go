package grace_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/envchain/envchain-cli/internal/grace"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *grace.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "grace-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return grace.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", 10*time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Project != "myproject" {
		t.Errorf("project = %q, want %q", rec.Project, "myproject")
	}
	if rec.Duration != 10*time.Minute {
		t.Errorf("duration = %s, want 10m", rec.Duration)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("missing")
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestIsActiveTrueWithinWindow(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, _ := m.Get("proj")
	if !rec.IsActive() {
		t.Error("expected grace period to be active")
	}
}

func TestIsActiveExpired(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", time.Millisecond); err != nil {
		t.Fatalf("Set: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	rec, _ := m.Get("proj")
	if rec.IsActive() {
		t.Error("expected grace period to have expired")
	}
}

func TestGuardBlocksActiveGrace(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", time.Hour)
	if err := m.Guard("proj"); !errors.Is(err, grace.ErrInGracePeriod) {
		t.Errorf("expected ErrInGracePeriod, got %v", err)
	}
}

func TestGuardPassesWhenNoRecord(t *testing.T) {
	m := newTempManager(t)
	if err := m.Guard("proj"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", time.Hour)
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := m.Guard("proj"); err != nil {
		t.Errorf("after delete Guard should pass, got %v", err)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", time.Minute); !errors.Is(err, grace.ErrEmptyProject) {
		t.Errorf("expected ErrEmptyProject, got %v", err)
	}
}

func TestSetNonPositiveDurationReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", 0); err == nil {
		t.Error("expected error for zero duration")
	}
}
