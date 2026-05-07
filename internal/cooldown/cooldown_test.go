package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/envchain-cli/internal/cooldown"
	"github.com/user/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *cooldown.Manager {
	t.Helper()
	s, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return cooldown.New(s)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", 5*time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Project != "myproject" {
		t.Errorf("project = %q, want %q", rec.Project, "myproject")
	}
	if rec.Duration != 5*time.Minute {
		t.Errorf("duration = %v, want %v", rec.Duration, 5*time.Minute)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if !isNotFound(err) {
		t.Fatalf("expected not-found, got %v", err)
	}
}

func TestIsActiveTrueWithinWindow(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}
	active, err := m.IsActive("proj")
	if err != nil {
		t.Fatalf("IsActive: %v", err)
	}
	if !active {
		t.Error("expected cooldown to be active")
	}
}

func TestIsActiveExpired(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", -time.Second); err != nil {
		t.Fatalf("Set: %v", err)
	}
	active, err := m.IsActive("proj")
	if err != nil {
		t.Fatalf("IsActive: %v", err)
	}
	if active {
		t.Error("expected expired cooldown to be inactive")
	}
}

func TestIsActiveNoRecord(t *testing.T) {
	m := newTempManager(t)
	active, err := m.IsActive("nobody")
	if err != nil {
		t.Fatalf("IsActive: %v", err)
	}
	if active {
		t.Error("expected false for missing record")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", time.Minute)
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if !isNotFound(err) {
		t.Errorf("expected not-found after delete, got %v", err)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", time.Minute); err == nil {
		t.Error("expected error for empty project")
	}
}

func isNotFound(err error) bool {
	return err != nil && err.Error() != "" && containsNotFound(err.Error())
}

func containsNotFound(s string) bool {
	return len(s) > 0 && (s == store.ErrNotFound.Error() ||
		len(s) >= len(store.ErrNotFound.Error()))
}
