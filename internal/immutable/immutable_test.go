package immutable_test

import (
	"errors"
	"os"
	"testing"

	"github.com/envchain/envchain-cli/internal/immutable"
	"github.com/envchain/envchain-cli/internal/store"
)

const testPass = "hunter2"

func newTempManager(t *testing.T) *immutable.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "immutable-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return immutable.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	keys := []string{"SECRET", "API_KEY", "TOKEN"}
	if err := m.Set("myproject", keys, testPass); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("myproject", testPass)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	// Result must be sorted.
	want := []string{"API_KEY", "SECRET", "TOKEN"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost", testPass)
	if err == nil {
		t.Fatal("expected error for missing project, got nil")
	}
}

func TestSetDeduplicates(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", []string{"A", "B", "A", "B", "C"}, testPass); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("proj", testPass)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 unique keys, got %d: %v", len(got), got)
	}
}

func TestIsImmutable(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", []string{"LOCKED", "ALSO_LOCKED"}, testPass)

	ok, err := m.IsImmutable("proj", "LOCKED", testPass)
	if err != nil {
		t.Fatalf("IsImmutable: %v", err)
	}
	if !ok {
		t.Error("expected LOCKED to be immutable")
	}

	ok, err = m.IsImmutable("proj", "FREE", testPass)
	if err != nil {
		t.Fatalf("IsImmutable: %v", err)
	}
	if ok {
		t.Error("expected FREE NOT to be immutable")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", []string{"KEY"}, testPass)
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj", testPass)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", []string{"KEY"}, testPass); !errors.Is(err, immutable.ErrEmptyProject) {
		t.Errorf("Set: expected ErrEmptyProject, got %v", err)
	}
	if _, err := m.Get("", testPass); !errors.Is(err, immutable.ErrEmptyProject) {
		t.Errorf("Get: expected ErrEmptyProject, got %v", err)
	}
}

func TestEmptyKeyReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", []string{"VALID", ""}, testPass); !errors.Is(err, immutable.ErrEmptyKey) {
		t.Errorf("expected ErrEmptyKey, got %v", err)
	}
}
