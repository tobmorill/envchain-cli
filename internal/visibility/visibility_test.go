package visibility_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain-cli/internal/store"
	"github.com/envchain/envchain-cli/internal/visibility"
)

func newTempManager(t *testing.T) *visibility.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "visibility-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return visibility.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", "SECRET_KEY", visibility.LevelHidden); err != nil {
		t.Fatalf("Set: %v", err)
	}
	lvl, err := m.Get("myproject", "SECRET_KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if lvl != visibility.LevelHidden {
		t.Errorf("expected hidden, got %q", lvl)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	lvl, err := m.Get("ghost", "MISSING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lvl != visibility.LevelVisible {
		t.Errorf("expected default visible, got %q", lvl)
	}
}

func TestSetNormalisesCase(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", "api_token", visibility.LevelHidden); err != nil {
		t.Fatalf("Set: %v", err)
	}
	lvl, err := m.Get("proj", "API_TOKEN")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if lvl != visibility.LevelHidden {
		t.Errorf("expected hidden after case normalisation, got %q", lvl)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", "DB_PASS", visibility.LevelHidden)
	if err := m.Delete("proj", "DB_PASS"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	lvl, _ := m.Get("proj", "DB_PASS")
	if lvl != visibility.LevelVisible {
		t.Errorf("expected visible after delete, got %q", lvl)
	}
}

func TestGetAll(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", "KEY_A", visibility.LevelHidden)
	_ = m.Set("proj", "KEY_B", visibility.LevelVisible)
	all, err := m.GetAll("proj")
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 settings, got %d", len(all))
	}
	if all["KEY_A"] != visibility.LevelHidden {
		t.Errorf("KEY_A should be hidden")
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", "KEY", visibility.LevelHidden); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestSetEmptyKeyReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", "", visibility.LevelHidden); err == nil {
		t.Error("expected error for empty key")
	}
}
