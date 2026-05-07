package immutable_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/immutable"
	"github.com/envchain/envchain-cli/internal/store"
	"os"
)

// TestMultipleProjectsAreIsolated verifies that immutable sets stored under
// different project names do not bleed into each other.
func TestMultipleProjectsAreIsolated(t *testing.T) {
	dir, err := os.MkdirTemp("", "immutable-integ-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	m := immutable.New(st)

	_ = m.Set("alpha", []string{"ALPHA_KEY"}, testPass)
	_ = m.Set("beta", []string{"BETA_KEY"}, testPass)

	alpha, err := m.Get("alpha", testPass)
	if err != nil {
		t.Fatalf("Get alpha: %v", err)
	}
	beta, err := m.Get("beta", testPass)
	if err != nil {
		t.Fatalf("Get beta: %v", err)
	}

	if len(alpha) != 1 || alpha[0] != "ALPHA_KEY" {
		t.Errorf("alpha: unexpected keys %v", alpha)
	}
	if len(beta) != 1 || beta[0] != "BETA_KEY" {
		t.Errorf("beta: unexpected keys %v", beta)
	}
}

// TestOverwriteReplacesEntireSet verifies that calling Set a second time fully
// replaces the previous key set rather than appending.
func TestOverwriteReplacesEntireSet(t *testing.T) {
	dir, err := os.MkdirTemp("", "immutable-integ-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	m := immutable.New(st)

	_ = m.Set("proj", []string{"OLD_KEY", "ANOTHER"}, testPass)
	_ = m.Set("proj", []string{"NEW_KEY"}, testPass)

	keys, err := m.Get("proj", testPass)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(keys) != 1 || keys[0] != "NEW_KEY" {
		t.Errorf("expected [NEW_KEY], got %v", keys)
	}
}
