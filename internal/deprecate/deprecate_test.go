package deprecate_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain-cli/internal/deprecate"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *deprecate.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "deprecate-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return deprecate.New(st)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	entries := []deprecate.Entry{
		{Key: "OLD_TOKEN", Replacement: "API_TOKEN", Reason: "renamed"},
	}
	if err := m.Set("myproject", entries); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(rec.Entries) != 1 || rec.Entries[0].Key != "OLD_TOKEN" {
		t.Errorf("unexpected entries: %+v", rec.Entries)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err != deprecate.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", nil); err == nil {
		t.Error("expected error for empty project name")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", []deprecate.Entry{{Key: "LEGACY"}})
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if err != deprecate.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestCheckReturnsMatchingKeys(t *testing.T) {
	m := newTempManager(t)
	entries := []deprecate.Entry{
		{Key: "OLD_HOST", Replacement: "DB_HOST"},
		{Key: "OLD_PORT"},
	}
	_ = m.Set("svc", entries)
	hits, err := m.Check("svc", []string{"OLD_HOST", "NEW_KEY", "ANOTHER"})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if len(hits) != 1 || hits[0].Key != "OLD_HOST" {
		t.Errorf("unexpected hits: %+v", hits)
	}
}

func TestCheckNoRecordReturnsEmpty(t *testing.T) {
	m := newTempManager(t)
	hits, err := m.Check("unknown", []string{"SOME_KEY"})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if len(hits) != 0 {
		t.Errorf("expected no hits, got %+v", hits)
	}
}

func TestCheckIsCaseInsensitive(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("app", []deprecate.Entry{{Key: "Legacy_Key"}})
	hits, err := m.Check("app", []string{"LEGACY_KEY"})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if len(hits) != 1 {
		t.Errorf("expected 1 hit, got %d", len(hits))
	}
}
