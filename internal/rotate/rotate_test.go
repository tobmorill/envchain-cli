package rotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envchain-cli/internal/chain"
	"github.com/your-org/envchain-cli/internal/env"
	"github.com/your-org/envchain-cli/internal/rotate"
	"github.com/your-org/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) (*rotate.Manager, *store.Store) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	st, err := store.New(path)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { os.Remove(path) })
	return rotate.New(st), st
}

func seedChain(t *testing.T, st *store.Store, name, passphrase string, entries []env.Entry) {
	t.Helper()
	cm := chain.New(st)
	if err := cm.Save(name, passphrase, entries); err != nil {
		t.Fatalf("seed chain %q: %v", name, err)
	}
}

func TestRotateSuccess(t *testing.T) {
	rm, st := newTempManager(t)
	entries := []env.Entry{{Key: "TOKEN", Value: "abc123"}}
	seedChain(t, st, "myproject", "old-pass", entries)

	if err := rm.Rotate("myproject", "old-pass", "new-pass"); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	cm := chain.New(st)
	got, err := cm.Load("myproject", "new-pass")
	if err != nil {
		t.Fatalf("Load after rotate: %v", err)
	}
	if len(got) != 1 || got[0].Key != "TOKEN" || got[0].Value != "abc123" {
		t.Errorf("unexpected entries after rotate: %v", got)
	}
}

func TestRotateSamePassphraseReturnsError(t *testing.T) {
	rm, st := newTempManager(t)
	seedChain(t, st, "proj", "same", []env.Entry{{Key: "K", Value: "V"}})

	if err := rm.Rotate("proj", "same", "same"); err == nil {
		t.Fatal("expected error when old and new passphrases are identical")
	}
}

func TestRotateWrongOldPassphrase(t *testing.T) {
	rm, st := newTempManager(t)
	seedChain(t, st, "proj", "correct", []env.Entry{{Key: "K", Value: "V"}})

	if err := rm.Rotate("proj", "wrong", "new-pass"); err == nil {
		t.Fatal("expected error for wrong old passphrase")
	}
}

func TestRotateAllSuccess(t *testing.T) {
	rm, st := newTempManager(t)
	names := []string{"alpha", "beta"}
	for _, n := range names {
		seedChain(t, st, n, "shared-old", []env.Entry{{Key: "X", Value: n}})
	}

	if err := rm.RotateAll(names, "shared-old", "shared-new"); err != nil {
		t.Fatalf("RotateAll: %v", err)
	}

	cm := chain.New(st)
	for _, n := range names {
		got, err := cm.Load(n, "shared-new")
		if err != nil {
			t.Errorf("Load %q after RotateAll: %v", n, err)
			continue
		}
		if len(got) == 0 || got[0].Value != n {
			t.Errorf("chain %q: unexpected entries %v", n, got)
		}
	}
}

func TestPreview(t *testing.T) {
	rm, st := newTempManager(t)
	want := []env.Entry{{Key: "SECRET", Value: "s3cr3t"}}
	seedChain(t, st, "preview-proj", "pass", want)

	got, err := rm.Preview("preview-proj", "pass")
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	if len(got) != 1 || got[0].Key != want[0].Key || got[0].Value != want[0].Value {
		t.Errorf("Preview returned %v, want %v", got, want)
	}
}
