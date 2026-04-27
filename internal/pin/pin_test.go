package pin_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain-cli/internal/pin"
)

func newTempManager(t *testing.T) *pin.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "pin-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return pin.New(dir)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	keys := []string{"DB_URL", "API_KEY"}
	if err := m.Set("myproject", keys); err != nil {
		t.Fatal(err)
	}
	got, err := m.Get("myproject")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "DB_URL" || got[1] != "API_KEY" {
		t.Fatalf("unexpected keys: %v", got)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("missing")
	if err != pin.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSetDeduplicatesAndNormalises(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", []string{"db_url", "DB_URL", " db_url "}); err != nil {
		t.Fatal(err)
	}
	got, err := m.Get("proj")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "DB_URL" {
		t.Fatalf("expected single normalised key, got %v", got)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", []string{"KEY"})
	if err := m.Delete("proj"); err != nil {
		t.Fatal(err)
	}
	_, err := m.Get("proj")
	if err != pin.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNoop(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("nonexistent"); err != nil {
		t.Fatalf("delete of missing project should not error: %v", err)
	}
}

func TestSetOverwrites(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", []string{"OLD_KEY"})
	_ = m.Set("proj", []string{"NEW_KEY"})
	got, err := m.Get("proj")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "NEW_KEY" {
		t.Fatalf("expected overwritten value, got %v", got)
	}
}
