package budget_test

import (
	"errors"
	"os"
	"testing"

	"github.com/user/envchain-cli/internal/budget"
)

func newTempManager(t *testing.T) *budget.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "budget-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return budget.New(dir)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", 4096); err != nil {
		t.Fatalf("Set: %v", err)
	}
	r, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if r.LimitBytes != 4096 {
		t.Errorf("LimitBytes = %d, want 4096", r.LimitBytes)
	}
	if r.Project != "myproject" {
		t.Errorf("Project = %q, want %q", r.Project, "myproject")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err == nil {
		t.Fatal("expected error for missing record")
	}
}

func TestCheckWithinLimit(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", 1000)
	if err := m.Check("proj", 500); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckExceedsLimit(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", 1000)
	err := m.Check("proj", 1500)
	if !errors.Is(err, budget.ErrExceeded) {
		t.Errorf("expected ErrExceeded, got %v", err)
	}
}

func TestCheckNoRecordAllows(t *testing.T) {
	m := newTempManager(t)
	if err := m.Check("unknown", 999999); err != nil {
		t.Errorf("expected nil for project with no budget record, got %v", err)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", 100); err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestSetNonPositiveLimitReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", 0); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", 512)
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDeleteNoop(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("nonexistent"); err != nil {
		t.Errorf("Delete of nonexistent should be noop, got %v", err)
	}
}
