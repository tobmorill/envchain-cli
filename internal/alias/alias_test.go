package alias_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain-cli/internal/alias"
)

func newTempManager(t *testing.T) *alias.Manager {
	t.Helper()
	dir := t.TempDir()
	return alias.New(filepath.Join(dir, "aliases.json"))
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("prod", "myapp-production"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "myapp-production" {
		t.Errorf("got %q, want %q", got, "myapp-production")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("missing")
	if !errors.Is(err, alias.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSetNormalizesCase(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("PROD", "myapp-production"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "myapp-production" {
		t.Errorf("got %q, want %q", got, "myapp-production")
	}
}

func TestInvalidAliasName(t *testing.T) {
	m := newTempManager(t)
	invalid := []string{"has space", "dot.name", "", "has/slash"}
	for _, name := range invalid {
		if err := m.Set(name, "chain"); !errors.Is(err, alias.ErrInvalidName) {
			t.Errorf("Set(%q): expected ErrInvalidName, got %v", name, err)
		}
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("staging", "myapp-staging")
	if err := m.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("staging")
	if !errors.Is(err, alias.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("nope"); !errors.Is(err, alias.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("prod", "myapp-production")
	_ = m.Set("dev", "myapp-dev")
	all, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(all))
	}
}

func TestPersistenceAcrossInstances(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.json")
	m1 := alias.New(path)
	_ = m1.Set("ci", "myapp-ci")
	m2 := alias.New(path)
	got, err := m2.Get("ci")
	if err != nil {
		t.Fatalf("Get from second manager: %v", err)
	}
	if got != "myapp-ci" {
		t.Errorf("got %q, want %q", got, "myapp-ci")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "aliases.json")
	m := alias.New(path)
	if err := m.Set("x", "chain-x"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
