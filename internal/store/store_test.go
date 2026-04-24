package store_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/envchain-cli/internal/store"
)

const testPass = "hunter2"

func newTempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	return store.New(dir)
}

func TestPutAndGet(t *testing.T) {
	s := newTempStore(t)
	set := store.EnvSet{
		Name: "dev",
		Vars: map[string]string{"FOO": "bar", "BAZ": "qux"},
	}
	if err := s.Put(set, testPass); err != nil {
		t.Fatalf("Put: %v", err)
	}
	got, err := s.Get("dev", testPass)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Vars["FOO"] != "bar" || got.Vars["BAZ"] != "qux" {
		t.Errorf("unexpected vars: %v", got.Vars)
	}
}

func TestGetNotFound(t *testing.T) {
	s := newTempStore(t)
	_, err := s.Get("missing", testPass)
	if err == nil {
		t.Fatal("expected error for missing env set")
	}
}

func TestDelete(t *testing.T) {
	s := newTempStore(t)
	set := store.EnvSet{Name: "staging", Vars: map[string]string{"KEY": "val"}}
	_ = s.Put(set, testPass)

	if err := s.Delete("staging", testPass); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get("staging", testPass)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestDeleteNotFound(t *testing.T) {
	s := newTempStore(t)
	if err := s.Delete("ghost", testPass); err == nil {
		t.Fatal("expected error deleting non-existent set")
	}
}

func TestList(t *testing.T) {
	s := newTempStore(t)
	for _, name := range []string{"alpha", "beta", "gamma"} {
		_ = s.Put(store.EnvSet{Name: name, Vars: map[string]string{}}, testPass)
	}
	names, err := s.List(testPass)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	for _, want := range []string{"alpha", "beta", "gamma"} {
		if !slices.Contains(names, want) {
			t.Errorf("missing name %q in list", want)
		}
	}
}

func TestWrongPassphrase(t *testing.T) {
	s := newTempStore(t)
	_ = s.Put(store.EnvSet{Name: "x", Vars: map[string]string{}}, testPass)
	_, err := s.Get("x", "wrongpass")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}

func TestFilePermissions(t *testing.T) {
	dir := t.TempDir()
	s := store.New(dir)
	_ = s.Put(store.EnvSet{Name: "perm", Vars: map[string]string{}}, testPass)

	info, err := os.Stat(filepath.Join(dir, "envchain.enc"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
