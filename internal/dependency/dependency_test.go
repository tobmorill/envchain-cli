package dependency_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/dependency"
)

func newTempManager(t *testing.T) *dependency.Manager {
	t.Helper()
	m, err := dependency.New(t.TempDir())
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}
	return m
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	deps := []string{"auth", "db"}
	if err := m.Set("api", deps); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("api")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got) != 2 || got[0] != "auth" || got[1] != "db" {
		t.Fatalf("unexpected deps: %v", got)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	got, err := m.Get("missing")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil slice, got: %v", got)
	}
}

func TestSetDeduplicates(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("api", []string{"auth", "auth", "db", "auth"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("api")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 unique deps, got %d: %v", len(got), got)
	}
}

func TestSelfDependencyReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set("api", []string{"auth", "api"})
	if err != dependency.ErrSelfDependency {
		t.Fatalf("expected ErrSelfDependency, got: %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("api", []string{"auth"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := m.Delete("api"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, err := m.Get("api")
	if err != nil {
		t.Fatalf("Get after Delete: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil after delete, got: %v", got)
	}
}

func TestEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", []string{"auth"}); err != dependency.ErrEmptyProject {
		t.Fatalf("expected ErrEmptyProject on Set, got: %v", err)
	}
	if _, err := m.Get(""); err != dependency.ErrEmptyProject {
		t.Fatalf("expected ErrEmptyProject on Get, got: %v", err)
	}
}
