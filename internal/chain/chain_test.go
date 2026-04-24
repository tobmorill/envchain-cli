package chain_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/env"
	"github.com/user/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *chain.Manager {
	t.Helper()
	dir := t.TempDir()
	s, err := store.New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return chain.New(s)
}

func TestSaveAndLoad(t *testing.T) {
	m := newTempManager(t)
	entries := []env.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	if err := m.Save("/my/project", "default", "secret", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := m.Load("/my/project", "default", "secret")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
	for i, e := range entries {
		if got[i].Key != e.Key || got[i].Value != e.Value {
			t.Errorf("entry %d: expected %v, got %v", i, e, got[i])
		}
	}
}

func TestLoadNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Load("/my/project", "missing", "secret")
	if !errors.Is(err, chain.ErrChainNotFound) {
		t.Errorf("expected ErrChainNotFound, got %v", err)
	}
}

func TestLoadWrongPassphrase(t *testing.T) {
	m := newTempManager(t)
	entries := []env.Entry{{Key: "K", Value: "V"}}
	if err := m.Save("/proj", "dev", "correct", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	_, err := m.Load("/proj", "dev", "wrong")
	if err == nil {
		t.Fatal("expected error with wrong passphrase, got nil")
	}
}

func TestDeleteChain(t *testing.T) {
	m := newTempManager(t)
	entries := []env.Entry{{Key: "X", Value: "1"}}
	if err := m.Save("/proj", "staging", "pass", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := m.Delete("/proj", "staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Load("/proj", "staging", "pass")
	if !errors.Is(err, chain.ErrChainNotFound) {
		t.Errorf("expected ErrChainNotFound after delete, got %v", err)
	}
}

func TestEmptyChainName(t *testing.T) {
	m := newTempManager(t)
	if err := m.Save("/proj", "", "pass", nil); err == nil {
		t.Error("expected error for empty chain name")
	}
	if _, err := m.Load("/proj", "", "pass"); err == nil {
		t.Error("expected error for empty chain name on load")
	}
}
