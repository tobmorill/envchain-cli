package inherit_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain-cli/internal/inherit"
	"github.com/user/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *inherit.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return inherit.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("child", "parent"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("child")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "parent" {
		t.Errorf("got %q, want %q", got, "parent")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	got, err := m.Get("missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSelfReferenceReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set("alpha", "alpha")
	if !errors.Is(err, inherit.ErrSelfReference) {
		t.Errorf("expected ErrSelfReference, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("child", "parent")
	if err := m.Delete("child"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, _ := m.Get("child")
	if got != "" {
		t.Errorf("expected empty after delete, got %q", got)
	}
}

func TestChainLinear(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("c", "b")
	_ = m.Set("b", "a")
	chain, err := m.Chain("c")
	if err != nil {
		t.Fatalf("Chain: %v", err)
	}
	want := []string{"b", "a"}
	if len(chain) != len(want) {
		t.Fatalf("chain length: got %d, want %d", len(chain), len(want))
	}
	for i, v := range want {
		if chain[i] != v {
			t.Errorf("chain[%d]: got %q, want %q", i, chain[i], v)
		}
	}
}

func TestChainCircularReturnsError(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("a", "b")
	_ = m.Set("b", "c")
	_ = m.Set("c", "a")
	_, err := m.Chain("a")
	if !errors.Is(err, inherit.ErrCircular) {
		t.Errorf("expected ErrCircular, got %v", err)
	}
}

func TestEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", "parent"); err == nil {
		t.Error("expected error for empty project")
	}
}
