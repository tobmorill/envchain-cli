package owner_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/owner"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *owner.Manager {
	t.Helper()
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return owner.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", "alice"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Owner != "alice" {
		t.Errorf("expected owner %q, got %q", "alice", rec.Owner)
	}
	if rec.Project != "myproject" {
		t.Errorf("expected project %q, got %q", "myproject", rec.Project)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err == nil {
		t.Fatal("expected error for unknown project, got nil")
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", "alice"); err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestSetEmptyOwnerReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", ""); err == nil {
		t.Fatal("expected error for empty owner")
	}
}

func TestSetOverwritesPrevious(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("myproject", "alice")
	_ = m.Set("myproject", "bob")
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Owner != "bob" {
		t.Errorf("expected owner %q after overwrite, got %q", "bob", rec.Owner)
	}
}

func TestDeleteRemovesRecord(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("myproject", "alice")
	if err := m.Delete("myproject"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("myproject")
	if err == nil {
		t.Fatal("expected not-found after delete")
	}
}

func TestDeleteNoopOnMissing(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("nonexistent"); err != nil {
		t.Fatalf("Delete on missing project should be a no-op, got: %v", err)
	}
}

func TestSetTrimmsWhitespace(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("  myproject  ", "  team-a  "); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Owner != "team-a" {
		t.Errorf("expected trimmed owner %q, got %q", "team-a", rec.Owner)
	}
}
