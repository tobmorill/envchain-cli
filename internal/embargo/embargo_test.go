package embargo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain-cli/internal/embargo"
	"github.com/user/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *embargo.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "embargo-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(filepath.Join(dir, "store.db"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return embargo.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	w := embargo.Window{StartHour: 9, StartMin: 0, EndHour: 17, EndMin: 30}
	if err := m.Set("myproject", w); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != w {
		t.Errorf("got %+v, want %+v", got, w)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set("", embargo.Window{})
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestSetInvalidHourReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set("proj", embargo.Window{StartHour: 25})
	if err == nil {
		t.Fatal("expected error for invalid hour")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	w := embargo.Window{StartHour: 8, EndHour: 18}
	if err := m.Set("proj", w); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestCheckNoWindowPermits(t *testing.T) {
	m := newTempManager(t)
	// No window set — should always permit.
	if err := m.Check("proj"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCheckMidnightSpan(t *testing.T) {
	// Window 22:00–06:00 spans midnight; just ensure no panic and logic runs.
	m := newTempManager(t)
	w := embargo.Window{StartHour: 22, StartMin: 0, EndHour: 6, EndMin: 0}
	if err := m.Set("proj", w); err != nil {
		t.Fatalf("Set: %v", err)
	}
	// We can't control the clock, so just ensure Check returns a typed error or nil.
	err := m.Check("proj")
	if err != nil && err != embargo.ErrEmbargoActive {
		t.Errorf("unexpected error type: %v", err)
	}
}
