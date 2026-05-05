package expiry_test

import (
	"os"
	"testing"
	"time"

	"github.com/envchain/envchain-cli/internal/expiry"
)

func newTempManager(t *testing.T) *expiry.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "expiry-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	m, err := expiry.New(dir)
	if err != nil {
		t.Fatalf("expiry.New: %v", err)
	}
	return m
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	at := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	if err := m.Set("myproject", at, "rotate soon"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Project != "myproject" {
		t.Errorf("project = %q, want %q", rec.Project, "myproject")
	}
	if !rec.ExpiresAt.Equal(at) {
		t.Errorf("expires_at = %v, want %v", rec.ExpiresAt, at)
	}
	if rec.Note != "rotate soon" {
		t.Errorf("note = %q, want %q", rec.Note, "rotate soon")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != expiry.ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestIsExpiredFuture(t *testing.T) {
	m := newTempManager(t)
	at := time.Now().Add(time.Hour)
	if err := m.Set("proj", at, ""); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, _ := m.Get("proj")
	if rec.IsExpired() {
		t.Error("future expiry should not be expired")
	}
}

func TestIsExpiredPast(t *testing.T) {
	m := newTempManager(t)
	at := time.Now().Add(-time.Second)
	if err := m.Set("proj", at, ""); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, _ := m.Get("proj")
	if !rec.IsExpired() {
		t.Error("past expiry should be expired")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", time.Now().Add(time.Hour), "")
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if err != expiry.ErrNotFound {
		t.Errorf("after delete: err = %v, want ErrNotFound", err)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set("", time.Now().Add(time.Hour), "")
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}
