package version_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain-cli/internal/store"
	"github.com/envchain/envchain-cli/internal/version"
)

func newTempManager(t *testing.T) *version.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "store")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return version.New(st)
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("myproject")
	if !errors.Is(err, version.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestBumpStartsAtOne(t *testing.T) {
	m := newTempManager(t)
	r, err := m.Bump("alpha")
	if err != nil {
		t.Fatalf("Bump: %v", err)
	}
	if r.Version != 1 {
		t.Fatalf("expected version 1, got %d", r.Version)
	}
	if r.Project != "alpha" {
		t.Fatalf("expected project alpha, got %q", r.Project)
	}
}

func TestBumpIncrementsMonotonically(t *testing.T) {
	m := newTempManager(t)
	for i := uint64(1); i <= 5; i++ {
		r, err := m.Bump("beta")
		if err != nil {
			t.Fatalf("Bump %d: %v", i, err)
		}
		if r.Version != i {
			t.Fatalf("expected version %d, got %d", i, r.Version)
		}
	}
}

func TestGetAfterBump(t *testing.T) {
	m := newTempManager(t)
	if _, err := m.Bump("gamma"); err != nil {
		t.Fatalf("Bump: %v", err)
	}
	r, err := m.Get("gamma")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if r.Version != 1 {
		t.Fatalf("expected version 1, got %d", r.Version)
	}
}

func TestProjectsAreIsolated(t *testing.T) {
	m := newTempManager(t)
	for i := 0; i < 3; i++ {
		if _, err := m.Bump("proj-a"); err != nil {
			t.Fatalf("Bump proj-a: %v", err)
		}
	}
	if _, err := m.Bump("proj-b"); err != nil {
		t.Fatalf("Bump proj-b: %v", err)
	}
	ra, _ := m.Get("proj-a")
	rb, _ := m.Get("proj-b")
	if ra.Version != 3 {
		t.Fatalf("proj-a: expected 3, got %d", ra.Version)
	}
	if rb.Version != 1 {
		t.Fatalf("proj-b: expected 1, got %d", rb.Version)
	}
}

func TestResetRemovesRecord(t *testing.T) {
	m := newTempManager(t)
	if _, err := m.Bump("delta"); err != nil {
		t.Fatalf("Bump: %v", err)
	}
	if err := m.Reset("delta"); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	_, err := m.Get("delta")
	if !errors.Is(err, version.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after reset, got %v", err)
	}
}

func TestResetNonExistentIsNoop(t *testing.T) {
	m := newTempManager(t)
	if err := m.Reset("ghost"); err != nil {
		t.Fatalf("Reset on missing project should not error: %v", err)
	}
}
