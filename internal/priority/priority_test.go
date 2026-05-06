package priority_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain-cli/internal/priority"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *priority.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "priority-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return priority.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", "API_KEY", priority.High); err != nil {
		t.Fatalf("Set: %v", err)
	}
	lvl, err := m.Get("myproject", "API_KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if lvl != priority.High {
		t.Errorf("expected High, got %v", lvl)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	lvl, err := m.Get("ghost", "MISSING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lvl != priority.Normal {
		t.Errorf("expected Normal default, got %v", lvl)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", "DB_PASS", priority.Low)
	if err := m.Delete("proj", "DB_PASS"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	lvl, _ := m.Get("proj", "DB_PASS")
	if lvl != priority.Normal {
		t.Errorf("expected Normal after delete, got %v", lvl)
	}
}

func TestGetAll(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", "A", priority.High)
	_ = m.Set("proj", "B", priority.Low)
	all, err := m.GetAll("proj")
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	if all["A"] != priority.High {
		t.Errorf("A: expected High, got %v", all["A"])
	}
	if all["B"] != priority.Low {
		t.Errorf("B: expected Low, got %v", all["B"])
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", "KEY", priority.High); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestParseLevelRoundtrip(t *testing.T) {
	for _, tc := range []struct {
		str string
		lvl priority.Level
	}{
		{"low", priority.Low},
		{"normal", priority.Normal},
		{"high", priority.High},
	} {
		got, err := priority.ParseLevel(tc.str)
		if err != nil {
			t.Errorf("ParseLevel(%q): %v", tc.str, err)
		}
		if got != tc.lvl {
			t.Errorf("ParseLevel(%q) = %v, want %v", tc.str, got, tc.lvl)
		}
		if got.String() != tc.str {
			t.Errorf("Level.String() = %q, want %q", got.String(), tc.str)
		}
	}
}

func TestParseLevelInvalid(t *testing.T) {
	_, err := priority.ParseLevel("critical")
	if err == nil {
		t.Error("expected error for unknown level")
	}
}
