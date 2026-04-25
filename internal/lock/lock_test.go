package lock_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/envchain-cli/internal/lock"
)

func newTempManager(t *testing.T) *lock.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "lock-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return lock.NewManager(dir)
}

func TestUnlockAndIsUnlocked(t *testing.T) {
	m := newTempManager(t)
	if err := m.Unlock("mychain", 5*time.Minute); err != nil {
		t.Fatalf("Unlock: %v", err)
	}
	ok, err := m.IsUnlocked("mychain")
	if err != nil {
		t.Fatalf("IsUnlocked: %v", err)
	}
	if !ok {
		t.Error("expected chain to be unlocked")
	}
}

func TestLockRemovesSession(t *testing.T) {
	m := newTempManager(t)
	_ = m.Unlock("mychain", 5*time.Minute)
	if err := m.Lock("mychain"); err != nil {
		t.Fatalf("Lock: %v", err)
	}
	ok, err := m.IsUnlocked("mychain")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected chain to be locked after Lock()")
	}
}

func TestLockNotUnlockedReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Lock("missing"); err != lock.ErrNotLocked {
		t.Errorf("expected ErrNotLocked, got %v", err)
	}
}

func TestExpiredSessionReportsLocked(t *testing.T) {
	m := newTempManager(t)
	if err := m.Unlock("mychain", -1*time.Second); err != nil {
		t.Fatalf("Unlock: %v", err)
	}
	ok, err := m.IsUnlocked("mychain")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected expired session to report locked")
	}
}

func TestIsUnlockedUnknownChain(t *testing.T) {
	m := newTempManager(t)
	ok, err := m.IsUnlocked("never-unlocked")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected unknown chain to be locked")
	}
}

func TestUnlockOverwritesExistingSession(t *testing.T) {
	m := newTempManager(t)
	_ = m.Unlock("mychain", 1*time.Minute)
	if err := m.Unlock("mychain", 10*time.Minute); err != nil {
		t.Fatalf("second Unlock: %v", err)
	}
	ok, err := m.IsUnlocked("mychain")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected chain to remain unlocked after re-unlock")
	}
}
