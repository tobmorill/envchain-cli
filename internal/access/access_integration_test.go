package access_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/access"
	"github.com/envchain/envchain-cli/internal/store"
)

// TestMultipleProjectsAreIsolated ensures that touching one project does not
// affect the access record of a different project.
func TestMultipleProjectsAreIsolated(t *testing.T) {
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	m := access.New(st)

	projects := []string{"alpha", "beta", "gamma"}
	for i, p := range projects {
		for j := 0; j <= i; j++ {
			if err := m.Touch(p); err != nil {
				t.Fatalf("Touch(%q): %v", p, err)
			}
		}
	}

	for i, p := range projects {
		rec, err := m.Get(p)
		if err != nil {
			t.Fatalf("Get(%q): %v", p, err)
		}
		want := int64(i + 1)
		if rec.Count != want {
			t.Errorf("%q: Count = %d, want %d", p, rec.Count, want)
		}
	}
}

// TestResetThenTouchStartsFresh verifies that after a reset, Touch begins a
// new record with Count == 1.
func TestResetThenTouchStartsFresh(t *testing.T) {
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	m := access.New(st)

	for i := 0; i < 10; i++ {
		_ = m.Touch("proj")
	}
	if err := m.Reset("proj"); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if err := m.Touch("proj"); err != nil {
		t.Fatalf("Touch after reset: %v", err)
	}
	rec, err := m.Get("proj")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Count != 1 {
		t.Errorf("Count = %d after reset+touch, want 1", rec.Count)
	}
}
