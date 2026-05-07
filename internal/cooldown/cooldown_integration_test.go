package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/envchain-cli/internal/cooldown"
	"github.com/user/envchain-cli/internal/store"
)

func TestMultipleProjectsAreIsolated(t *testing.T) {
	s, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	m := cooldown.New(s)

	if err := m.Set("alpha", time.Hour); err != nil {
		t.Fatalf("Set alpha: %v", err)
	}
	if err := m.Set("beta", -time.Second); err != nil {
		t.Fatalf("Set beta: %v", err)
	}

	alphaActive, err := m.IsActive("alpha")
	if err != nil {
		t.Fatalf("IsActive alpha: %v", err)
	}
	betaActive, err := m.IsActive("beta")
	if err != nil {
		t.Fatalf("IsActive beta: %v", err)
	}

	if !alphaActive {
		t.Error("alpha should be active")
	}
	if betaActive {
		t.Error("beta should be expired")
	}
}

func TestOverwriteUpdatesCooldown(t *testing.T) {
	s, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	m := cooldown.New(s)

	if err := m.Set("proj", -time.Second); err != nil {
		t.Fatalf("Set (expired): %v", err)
	}
	active, _ := m.IsActive("proj")
	if active {
		t.Fatal("expected inactive before overwrite")
	}

	if err := m.Set("proj", time.Hour); err != nil {
		t.Fatalf("Set (active): %v", err)
	}
	active, err := m.IsActive("proj")
	if err != nil {
		t.Fatalf("IsActive: %v", err)
	}
	if !active {
		t.Error("expected active after overwrite")
	}
}

func TestDeleteThenCheckIsInactive(t *testing.T) {
	s, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	m := cooldown.New(s)

	_ = m.Set("proj", time.Hour)
	_ = m.Delete("proj")

	active, err := m.IsActive("proj")
	if err != nil {
		t.Fatalf("IsActive after delete: %v", err)
	}
	if active {
		t.Error("expected inactive after delete")
	}
}
