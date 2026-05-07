package namespace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain-cli/internal/namespace"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *namespace.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "ns.db")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return namespace.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	projects := []string{"alpha", "beta", "gamma"}
	if err := m.Set("team-a", projects); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("team-a")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(rec.Projects) != 3 {
		t.Fatalf("expected 3 projects, got %d", len(rec.Projects))
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("missing")
	if err != namespace.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSetDeduplicatesProjects(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("dup", []string{"proj", "PROJ", "proj"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("dup")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(rec.Projects) != 1 {
		t.Fatalf("expected 1 project after dedup, got %d", len(rec.Projects))
	}
}

func TestInvalidNameReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("bad name!", []string{"p"}); err != namespace.ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("ns1", []string{"x"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := m.Delete("ns1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := m.Get("ns1"); err != namespace.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("ghost"); err != namespace.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSetOverwrites(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("ns", []string{"old"})
	_ = m.Set("ns", []string{"new1", "new2"})
	rec, err := m.Get("ns")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(rec.Projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(rec.Projects))
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
