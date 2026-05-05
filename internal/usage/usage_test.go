package usage_test

import (
	"os"
	"testing"
	"time"

	"github.com/envchain/envchain-cli/internal/store"
	"github.com/envchain/envchain-cli/internal/usage"
)

func newTempManager(t *testing.T) *usage.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "usage-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return usage.New(s)
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != nil {
		t.Fatalf("expected nil record, got %+v", rec)
	}
}

func TestTouchCreatesRecord(t *testing.T) {
	m := newTempManager(t)
	before := time.Now().UTC().Add(-time.Second)

	if err := m.Touch("alpha"); err != nil {
		t.Fatalf("Touch: %v", err)
	}

	rec, err := m.Get("alpha")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec == nil {
		t.Fatal("expected record, got nil")
	}
	if rec.Count != 1 {
		t.Errorf("count: want 1, got %d", rec.Count)
	}
	if rec.FirstUsed.Before(before) {
		t.Errorf("FirstUsed %v is before test start %v", rec.FirstUsed, before)
	}
}

func TestTouchIncrementsCount(t *testing.T) {
	m := newTempManager(t)

	for i := 0; i < 5; i++ {
		if err := m.Touch("beta"); err != nil {
			t.Fatalf("Touch #%d: %v", i, err)
		}
	}

	rec, err := m.Get("beta")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Count != 5 {
		t.Errorf("count: want 5, got %d", rec.Count)
	}
}

func TestTouchUpdatesLastUsed(t *testing.T) {
	m := newTempManager(t)

	if err := m.Touch("gamma"); err != nil {
		t.Fatalf("first Touch: %v", err)
	}
	rec1, _ := m.Get("gamma")

	time.Sleep(2 * time.Millisecond)

	if err := m.Touch("gamma"); err != nil {
		t.Fatalf("second Touch: %v", err)
	}
	rec2, _ := m.Get("gamma")

	if !rec2.LastUsed.After(rec1.LastUsed) {
		t.Errorf("LastUsed not updated: first=%v second=%v", rec1.LastUsed, rec2.LastUsed)
	}
	if !rec2.FirstUsed.Equal(rec1.FirstUsed) {
		t.Errorf("FirstUsed changed unexpectedly")
	}
}

func TestResetRemovesRecord(t *testing.T) {
	m := newTempManager(t)

	if err := m.Touch("delta"); err != nil {
		t.Fatalf("Touch: %v", err)
	}
	if err := m.Reset("delta"); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	rec, err := m.Get("delta")
	if err != nil {
		t.Fatalf("Get after reset: %v", err)
	}
	if rec != nil {
		t.Errorf("expected nil after reset, got %+v", rec)
	}
}

func TestResetNoopOnMissing(t *testing.T) {
	m := newTempManager(t)
	if err := m.Reset("nonexistent"); err != nil {
		t.Errorf("Reset on missing should not error, got: %v", err)
	}
}

func TestTouchEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Touch(""); err == nil {
		t.Error("expected error for empty project name")
	}
}
