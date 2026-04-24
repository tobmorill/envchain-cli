package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain-cli/internal/config"
)

func newTempManager(t *testing.T) *config.Manager {
	t.Helper()
	dir := t.TempDir()
	mgr, err := config.NewManager(dir)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return mgr
}

func TestLoadMissingReturnsZero(t *testing.T) {
	mgr := newTempManager(t)
	cfg, err := mgr.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultShell != "" || cfg.StorePath != "" {
		t.Errorf("expected zero config, got %+v", cfg)
	}
}

func TestSaveAndLoad(t *testing.T) {
	mgr := newTempManager(t)
	want := config.Config{
		DefaultShell:   "zsh",
		StorePath:      "/tmp/mystore",
		PassphraseHint: "pet name",
	}
	if err := mgr.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := mgr.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	base := t.TempDir()
	nestedDir := filepath.Join(base, "a", "b", "c")
	mgr, err := config.NewManager(nestedDir)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	if err := mgr.Save(config.Config{DefaultShell: "bash"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(mgr.Path()); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}

func TestSaveOverwrites(t *testing.T) {
	mgr := newTempManager(t)
	if err := mgr.Save(config.Config{DefaultShell: "bash"}); err != nil {
		t.Fatalf("first Save: %v", err)
	}
	if err := mgr.Save(config.Config{DefaultShell: "fish"}); err != nil {
		t.Fatalf("second Save: %v", err)
	}
	cfg, err := mgr.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DefaultShell != "fish" {
		t.Errorf("expected fish, got %q", cfg.DefaultShell)
	}
}

func TestFilePermissions(t *testing.T) {
	mgr := newTempManager(t)
	if err := mgr.Save(config.Config{DefaultShell: "zsh"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	info, err := os.Stat(mgr.Path())
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected permissions 0600, got %04o", perm)
	}
}
