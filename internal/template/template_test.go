package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain-cli/internal/store"
	"github.com/envchain/envchain-cli/internal/template"
)

func newTempManager(t *testing.T) *template.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "store.db")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { os.Remove(dir) })
	return template.New(st)
}

func TestSaveAndLoad(t *testing.T) {
	m := newTempManager(t)
	tmpl := template.Template{Name: "base", Keys: []string{"API_KEY", "DB_URL"}}
	if err := m.Save(tmpl); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := m.Load("base")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != tmpl.Name || len(got.Keys) != len(tmpl.Keys) {
		t.Errorf("got %+v, want %+v", got, tmpl)
	}
}

func TestLoadNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Load("nonexistent")
	if err != template.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteTemplate(t *testing.T) {
	m := newTempManager(t)
	tmpl := template.Template{Name: "todelete", Keys: []string{"X"}}
	if err := m.Save(tmpl); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := m.Delete("todelete"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Load("todelete")
	if err != template.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("ghost"); err != template.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestInvalidName(t *testing.T) {
	m := newTempManager(t)
	badNames := []string{"", "has space", "slash/name", "dot.name"}
	for _, name := range badNames {
		err := m.Save(template.Template{Name: name, Keys: []string{"K"}})
		if err != template.ErrInvalidName {
			t.Errorf("name %q: expected ErrInvalidName, got %v", name, err)
		}
	}
}

func TestIsValidName(t *testing.T) {
	valid := []string{"base", "my-env", "env_123", "A"}
	for _, n := range valid {
		if !template.IsValidName(n) {
			t.Errorf("%q should be valid", n)
		}
	}
	invalid := []string{"", "bad name", "x/y"}
	for _, n := range invalid {
		if template.IsValidName(n) {
			t.Errorf("%q should be invalid", n)
		}
	}
}
