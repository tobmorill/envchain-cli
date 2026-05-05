package notes_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain-cli/internal/notes"
	"github.com/envchain/envchain-cli/internal/store"
)

const testPass = "hunter2"

func newTempManager(t *testing.T) *notes.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "notes-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return notes.New(s)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", "remember to rotate keys", testPass); err != nil {
		t.Fatalf("Set: %v", err)
	}
	n, err := m.Get("myproject", testPass)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if n.Body != "remember to rotate keys" {
		t.Errorf("body = %q, want %q", n.Body, "remember to rotate keys")
	}
	if n.Project != "myproject" {
		t.Errorf("project = %q, want %q", n.Project, "myproject")
	}
	if n.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost", testPass)
	if err == nil {
		t.Fatal("expected error for missing note, got nil")
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set("", "body", testPass)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestSetOverwritesPrevious(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", "first", testPass)
	_ = m.Set("proj", "second", testPass)
	n, err := m.Get("proj", testPass)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if n.Body != "second" {
		t.Errorf("body = %q, want %q", n.Body, "second")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", "note", testPass)
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj", testPass)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestGetWrongPassphrase(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", "secret note", testPass)
	_, err := m.Get("proj", "wrongpass")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}
