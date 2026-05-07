package shield_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain-cli/internal/shield"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *shield.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "shield.db")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { os.Remove(dir) })
	return shield.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", []string{"SECRET", "API_KEY"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	keys, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	keys, err := m.Get("ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys != nil {
		t.Fatalf("expected nil, got %v", keys)
	}
}

func TestSetNormalisesAndDeduplicates(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", []string{"api_key", "API_KEY", " secret ", "SECRET"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	keys, _ := m.Get("proj")
	if len(keys) != 2 {
		t.Fatalf("expected 2 deduplicated keys, got %d: %v", len(keys), keys)
	}
}

func TestGuardBlocks(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", []string{"DB_PASSWORD"})
	err := m.Guard("proj", "db_password")
	if !errors.Is(err, shield.ErrShielded) {
		t.Fatalf("expected ErrShielded, got %v", err)
	}
}

func TestGuardAllows(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", []string{"DB_PASSWORD"})
	if err := m.Guard("proj", "SAFE_KEY"); err != nil {
		t.Fatalf("unexpected block: %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", []string{"KEY"})
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	keys, _ := m.Get("proj")
	if len(keys) != 0 {
		t.Fatalf("expected empty after delete, got %v", keys)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", []string{"KEY"}); !errors.Is(err, shield.ErrEmptyProject) {
		t.Fatalf("expected ErrEmptyProject, got %v", err)
	}
}
