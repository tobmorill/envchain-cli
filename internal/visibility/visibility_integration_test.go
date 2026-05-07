package visibility_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/visibility"
)

func TestMultipleProjectsAreIsolated(t *testing.T) {
	m := newTempManager(t)

	_ = m.Set("alpha", "TOKEN", visibility.LevelHidden)
	_ = m.Set("beta", "TOKEN", visibility.LevelVisible)

	alpha, err := m.Get("alpha", "TOKEN")
	if err != nil {
		t.Fatalf("Get alpha: %v", err)
	}
	beta, err := m.Get("beta", "TOKEN")
	if err != nil {
		t.Fatalf("Get beta: %v", err)
	}

	if alpha != visibility.LevelHidden {
		t.Errorf("alpha TOKEN should be hidden, got %q", alpha)
	}
	if beta != visibility.LevelVisible {
		t.Errorf("beta TOKEN should be visible, got %q", beta)
	}
}

func TestOverwriteUpdatesLevel(t *testing.T) {
	m := newTempManager(t)

	_ = m.Set("proj", "DB_PASS", visibility.LevelHidden)
	_ = m.Set("proj", "DB_PASS", visibility.LevelVisible)

	lvl, err := m.Get("proj", "DB_PASS")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if lvl != visibility.LevelVisible {
		t.Errorf("expected visible after overwrite, got %q", lvl)
	}
}

func TestDeleteDoesNotAffectOtherKeys(t *testing.T) {
	m := newTempManager(t)

	_ = m.Set("proj", "KEY_A", visibility.LevelHidden)
	_ = m.Set("proj", "KEY_B", visibility.LevelHidden)
	_ = m.Delete("proj", "KEY_A")

	lvlA, _ := m.Get("proj", "KEY_A")
	lvlB, _ := m.Get("proj", "KEY_B")

	if lvlA != visibility.LevelVisible {
		t.Errorf("KEY_A should default to visible after delete, got %q", lvlA)
	}
	if lvlB != visibility.LevelHidden {
		t.Errorf("KEY_B should remain hidden, got %q", lvlB)
	}
}
