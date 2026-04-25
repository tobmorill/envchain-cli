package tags_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/internal/store"
	"github.com/envchain-cli/internal/tags"
)

func newTempManager(t *testing.T) *tags.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "store")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return tags.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	err := m.Set("proj:default", []string{"production", "backend"})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("proj:default")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(got))
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("nonexistent")
	if err == nil {
		t.Fatal("expected ErrTagNotFound, got nil")
	}
	if err != tags.ErrTagNotFound {
		t.Fatalf("expected ErrTagNotFound, got %v", err)
	}
}

func TestNormalizesAndDeduplicates(t *testing.T) {
	m := newTempManager(t)
	err := m.Set("k", []string{"Beta", "alpha", "ALPHA", "", "  beta  "})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("k")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 unique tags, got %v", got)
	}
	if got[0] != "alpha" || got[1] != "beta" {
		t.Fatalf("unexpected tags: %v", got)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("k", []string{"x"})
	if err := m.Delete("k"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("k")
	if err != tags.ErrTagNotFound {
		t.Fatalf("expected ErrTagNotFound after delete, got %v", err)
	}
}

func TestFindByTag(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj-a:default", []string{"production", "backend"})
	_ = m.Set("proj-b:default", []string{"staging", "backend"})
	_ = m.Set("proj-c:default", []string{"production"})

	matches, err := m.FindByTag("backend")
	if err != nil {
		t.Fatalf("FindByTag: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %v", matches)
	}

	matches, err = m.FindByTag("PRODUCTION")
	if err != nil {
		t.Fatalf("FindByTag case-insensitive: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 production matches, got %v", matches)
	}
}

func TestSetEmptyKeyReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", []string{"x"}); err == nil {
		t.Fatal("expected error for empty chain key")
	}
}

func init() {
	// ensure test binary can locate a writable temp dir
	_ = os.MkdirTemp("", "tags-test")
}
