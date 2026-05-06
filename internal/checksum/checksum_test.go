package checksum_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/checksum"
	"github.com/envchain/envchain-cli/internal/env"
	"github.com/envchain/envchain-cli/internal/store"
)

func newTempManager(t *testing.T) *checksum.Manager {
	t.Helper()
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return checksum.New(st)
}

func sampleEntries() []env.Entry {
	return []env.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
}

func TestSaveAndVerify(t *testing.T) {
	m := newTempManager(t)
	entries := sampleEntries()

	if err := m.Save("myproject", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := m.Verify("myproject", entries); err != nil {
		t.Fatalf("Verify: %v", err)
	}
}

func TestVerifyMismatch(t *testing.T) {
	m := newTempManager(t)
	entries := sampleEntries()

	if err := m.Save("myproject", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}

	modified := append(sampleEntries(), env.Entry{Key: "EXTRA", Value: "val"})
	err := m.Verify("myproject", modified)
	if err == nil {
		t.Fatal("expected mismatch error, got nil")
	}
	if err != checksum.ErrMismatch {
		t.Fatalf("expected ErrMismatch, got %v", err)
	}
}

func TestVerifyNotFound(t *testing.T) {
	m := newTempManager(t)
	err := m.Verify("ghost", sampleEntries())
	if err == nil {
		t.Fatal("expected error for missing checksum, got nil")
	}
}

func TestSaveIsDeterministic(t *testing.T) {
	// Saving twice with the same entries should succeed and verify cleanly.
	m := newTempManager(t)
	entries := sampleEntries()

	if err := m.Save("proj", entries); err != nil {
		t.Fatalf("first Save: %v", err)
	}
	if err := m.Save("proj", entries); err != nil {
		t.Fatalf("second Save: %v", err)
	}
	if err := m.Verify("proj", entries); err != nil {
		t.Fatalf("Verify after double save: %v", err)
	}
}

func TestOrderIndependentHash(t *testing.T) {
	m := newTempManager(t)

	a := []env.Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	b := []env.Entry{{Key: "B", Value: "2"}, {Key: "A", Value: "1"}}

	if err := m.Save("proj", a); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := m.Verify("proj", b); err != nil {
		t.Fatalf("Verify with reordered entries: %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	entries := sampleEntries()

	if err := m.Save("proj", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := m.Verify("proj", entries); err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestSaveEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Save("", sampleEntries()); err == nil {
		t.Fatal("expected error for empty project name")
	}
}
