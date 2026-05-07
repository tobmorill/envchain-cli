package readonly_test

import (
	"path/filepath"
	"testing"

	"github.com/envchain/envchain-cli/internal/readonly"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *readonly.Manager {
	t.Helper()
	st, err := store.New(filepath.Join(t.TempDir(), "store.db"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return readonly.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", true); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !rec.ReadOnly {
		t.Error("expected ReadOnly=true")
	}
	if rec.Project != "myproject" {
		t.Errorf("expected project=myproject, got %q", rec.Project)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	rec, err := m.Get("unknown")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.ReadOnly {
		t.Error("expected ReadOnly=false for unknown project")
	}
}

func TestIsReadOnly(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", true)
	ok, err := m.IsReadOnly("proj")
	if err != nil {
		t.Fatalf("IsReadOnly: %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
}

func TestGuardBlocks(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("locked", true)
	if err := m.Guard("locked"); err == nil {
		t.Error("expected error for read-only project")
	}
}

func TestGuardAllows(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("open", false)
	if err := m.Guard("open"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", true)
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	ok, _ := m.IsReadOnly("proj")
	if ok {
		t.Error("expected ReadOnly=false after delete")
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", true); err == nil {
		t.Error("expected error for empty project")
	}
}
