package ttl_test

import (
	"testing"
	"time"

	"github.com/yourorg/envchain-cli/internal/ttl"
)

func TestSetAndGet(t *testing.T) {
	c := ttl.New()
	c.Set("key", "value", 5*time.Second)
	v, err := c.Get("key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "value" {
		t.Fatalf("expected %q, got %q", "value", v)
	}
}

func TestGetNotFound(t *testing.T) {
	c := ttl.New()
	_, err := c.Get("missing")
	if err != ttl.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetExpired(t *testing.T) {
	c := ttl.New()
	c.Set("key", "secret", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_, err := c.Get("key")
	if err != ttl.ErrExpired {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	c := ttl.New()
	c.Set("key", "value", 5*time.Second)
	c.Delete("key")
	_, err := c.Get("key")
	if err != ttl.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNoop(t *testing.T) {
	c := ttl.New()
	// Should not panic
	c.Delete("nonexistent")
}

func TestPurgeRemovesExpired(t *testing.T) {
	c := ttl.New()
	c.Set("short", "a", 1*time.Millisecond)
	c.Set("long", "b", 10*time.Second)
	time.Sleep(5 * time.Millisecond)
	c.Purge()
	if _, err := c.Get("long"); err != nil {
		t.Fatalf("long-lived entry should still be valid: %v", err)
	}
	if _, err := c.Get("short"); err != ttl.ErrNotFound {
		t.Fatalf("expected short entry to be purged, got %v", err)
	}
}

func TestOverwriteEntry(t *testing.T) {
	c := ttl.New()
	c.Set("key", "first", 5*time.Second)
	c.Set("key", "second", 5*time.Second)
	v, err := c.Get("key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "second" {
		t.Fatalf("expected overwritten value %q, got %q", "second", v)
	}
}
