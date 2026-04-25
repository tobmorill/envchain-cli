package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/snapshot"
	"github.com/user/envchain-cli/internal/store"
)

const testPassphrase = "hunter2"

func newTempManager(t *testing.T) (*snapshot.Manager, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "snapshot-test-*")
	if err != nil {
		t.Fatalf("mkdtemp: %v", err)
	}
	st, err := store.New(filepath.Join(dir, "store.db"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	cm := chain.New(st)
	mgr := snapshot.New(st, cm)
	return mgr, func() {
		st.Close()
		os.RemoveAll(dir)
	}
}

func seedChain(t *testing.T, dir string, project, passphrase string) *snapshot.Manager {
	t.Helper()
	st, _ := store.New(filepath.Join(dir, "store.db"))
	cm := chain.New(st)
	entries := []string{"FOO=bar", "BAZ=qux"}
	if err := cm.Save(project, passphrase, entries); err != nil {
		t.Fatalf("seed chain: %v", err)
	}
	return snapshot.New(st, cm)
}

func TestSaveAndGet(t *testing.T) {
	dir, _ := os.MkdirTemp("", "snap-*")
	defer os.RemoveAll(dir)

	st, _ := store.New(filepath.Join(dir, "store.db"))
	defer st.Close()
	cm := chain.New(st)
	_ = cm.Save("myproject", testPassphrase, []string{"KEY=val"})

	mgr := snapshot.New(st, cm)
	if err := mgr.Save("myproject", "v1", testPassphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := mgr.Get("myproject", "v1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if snap.Project != "myproject" {
		t.Errorf("project = %q, want %q", snap.Project, "myproject")
	}
	if snap.Label != "v1" {
		t.Errorf("label = %q, want %q", snap.Label, "v1")
	}
	if len(snap.Entries) != 1 || snap.Entries[0] != "KEY=val" {
		t.Errorf("entries = %v, want [KEY=val]", snap.Entries)
	}
	if snap.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestGetNotFound(t *testing.T) {
	mgr, cleanup := newTempManager(t)
	defer cleanup()

	_, err := mgr.Get("ghost", "v0")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestDelete(t *testing.T) {
	dir, _ := os.MkdirTemp("", "snap-del-*")
	defer os.RemoveAll(dir)

	st, _ := store.New(filepath.Join(dir, "store.db"))
	defer st.Close()
	cm := chain.New(st)
	_ = cm.Save("proj", testPassphrase, []string{"A=1"})

	mgr := snapshot.New(st, cm)
	_ = mgr.Save("proj", "snap1", testPassphrase)

	if err := mgr.Delete("proj", "snap1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := mgr.Get("proj", "snap1"); err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestDeleteNotFound(t *testing.T) {
	mgr, cleanup := newTempManager(t)
	defer cleanup()

	if err := mgr.Delete("proj", "nonexistent"); err == nil {
		t.Fatal("expected error deleting missing snapshot")
	}
}
