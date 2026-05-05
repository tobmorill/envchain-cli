package policy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envchain-cli/internal/policy"
)

func newTempManager(t *testing.T) *policy.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "policies")
	return policy.New(dir)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	r := policy.Rule{AllowKeys: []string{"DB_URL", "API_KEY"}, DenyKeys: []string{"SECRET"}}
	if err := m.Set("myproject", r); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got.AllowKeys) != 2 || got.AllowKeys[0] != "DB_URL" {
		t.Errorf("unexpected AllowKeys: %v", got.AllowKeys)
	}
	if len(got.DenyKeys) != 1 || got.DenyKeys[0] != "SECRET" {
		t.Errorf("unexpected DenyKeys: %v", got.DenyKeys)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err != policy.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", policy.Rule{})
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := m.Get("proj"); err != policy.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("missing"); err != policy.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAllowedDenyKeys(t *testing.T) {
	r := policy.Rule{DenyKeys: []string{"SECRET"}}
	ok, err := policy.Allowed(r, "SECRET")
	if err != nil || ok {
		t.Errorf("expected denied; ok=%v err=%v", ok, err)
	}
	ok, err = policy.Allowed(r, "DB_URL")
	if err != nil || !ok {
		t.Errorf("expected allowed; ok=%v err=%v", ok, err)
	}
}

func TestAllowedAllowKeys(t *testing.T) {
	r := policy.Rule{AllowKeys: []string{"DB_URL"}}
	ok, _ := policy.Allowed(r, "DB_URL")
	if !ok {
		t.Error("expected DB_URL to be allowed")
	}
	ok, _ = policy.Allowed(r, "OTHER")
	if ok {
		t.Error("expected OTHER to be denied")
	}
}

func TestAllowedPattern(t *testing.T) {
	r := policy.Rule{AllowPattern: `^APP_`}
	ok, err := policy.Allowed(r, "APP_PORT")
	if err != nil || !ok {
		t.Errorf("expected APP_PORT allowed; ok=%v err=%v", ok, err)
	}
	ok, err = policy.Allowed(r, "DB_PASS")
	if err != nil || ok {
		t.Errorf("expected DB_PASS denied; ok=%v err=%v", ok, err)
	}
}

func TestAllowedInvalidPattern(t *testing.T) {
	r := policy.Rule{AllowPattern: `[invalid`}
	_, err := policy.Allowed(r, "KEY")
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestSetCreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "deep", "policies")
	m := policy.New(dir)
	if err := m.Set("p", policy.Rule{}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
}
