package cache_test

import (
	"testing"
	"time"

	"github.com/yourorg/envchain-cli/internal/cache"
	"github.com/yourorg/envchain-cli/internal/ttl"
)

func TestStoreAndRetrieve(t *testing.T) {
	m := cache.New(5 * time.Second)
	m.Store("myproject", "s3cr3t")
	v, err := m.Retrieve("myproject")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "s3cr3t" {
		t.Fatalf("expected %q, got %q", "s3cr3t", v)
	}
}

func TestRetrieveNotFound(t *testing.T) {
	m := cache.New(5 * time.Second)
	_, err := m.Retrieve("unknown")
	if err != ttl.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRetrieveExpired(t *testing.T) {
	m := cache.New(1 * time.Millisecond)
	m.Store("proj", "pass")
	time.Sleep(5 * time.Millisecond)
	_, err := m.Retrieve("proj")
	if err != ttl.ErrExpired {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestInvalidate(t *testing.T) {
	m := cache.New(5 * time.Second)
	m.Store("proj", "pass")
	m.Invalidate("proj")
	_, err := m.Retrieve("proj")
	if err != ttl.ErrNotFound {
		t.Fatalf("expected ErrNotFound after invalidate, got %v", err)
	}
}

func TestDefaultTTLUsedWhenZero(t *testing.T) {
	m := cache.New(0)
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
	m.Store("proj", "pass")
	v, err := m.Retrieve("proj")
	if err != nil {
		t.Fatalf("unexpected error with default TTL: %v", err)
	}
	if v != "pass" {
		t.Fatalf("expected %q, got %q", "pass", v)
	}
}

// TestStoreOverwritesExistingEntry verifies that storing a new value for an
// existing key replaces the previous value and resets the TTL.
func TestStoreOverwritesExistingEntry(t *testing.T) {
	m := cache.New(5 * time.Second)
	m.Store("proj", "old-pass")
	m.Store("proj", "new-pass")
	v, err := m.Retrieve("proj")
	if err != nil {
		t.Fatalf("unexpected error after overwrite: %v", err)
	}
	if v != "new-pass" {
		t.Fatalf("expected %q after overwrite, got %q", "new-pass", v)
	}
}
