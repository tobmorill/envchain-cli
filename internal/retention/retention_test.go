package retention_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/envchain-cli/internal/retention"
	"github.com/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *retention.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "retention-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return retention.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	p := retention.Policy{
		Project:     "myapp",
		MaxAge:      7 * 24 * time.Hour,
		MaxVersions: 10,
	}
	if err := m.Set(p); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("myapp")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.MaxAge != p.MaxAge {
		t.Errorf("MaxAge: got %v, want %v", got.MaxAge, p.MaxAge)
	}
	if got.MaxVersions != p.MaxVersions {
		t.Errorf("MaxVersions: got %d, want %d", got.MaxVersions, p.MaxVersions)
	}
	if got.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if !errors.Is(err, retention.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set(retention.Policy{MaxAge: time.Hour})
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set(retention.Policy{Project: "myapp", MaxAge: time.Hour})
	if err := m.Delete("myapp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("myapp")
	if !errors.Is(err, retention.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newTempManager(t)
	err := m.Delete("nonexistent")
	if !errors.Is(err, retention.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestShouldPruneByAge(t *testing.T) {
	p := retention.Policy{Project: "x", MaxAge: 24 * time.Hour}
	old := time.Now().Add(-48 * time.Hour)
	recent := time.Now().Add(-1 * time.Hour)
	if !p.ShouldPrune(old) {
		t.Error("expected old record to be pruned")
	}
	if p.ShouldPrune(recent) {
		t.Error("expected recent record not to be pruned")
	}
}

func TestShouldPruneZeroMaxAge(t *testing.T) {
	p := retention.Policy{Project: "x", MaxAge: 0}
	if p.ShouldPrune(time.Now().Add(-365 * 24 * time.Hour)) {
		t.Error("zero MaxAge should never prune")
	}
}
