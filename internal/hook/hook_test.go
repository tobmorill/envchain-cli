package hook_test

import (
	"os"
	"testing"

	"github.com/yourorg/envchain-cli/internal/hook"
)

func newTempManager(t *testing.T) *hook.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "hook-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return hook.New(dir)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	h := hook.Hook{Project: "myapp", Phase: hook.PhasePre, Command: "echo pre"}
	if err := m.Set(h); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, ok, err := m.Get("myapp", hook.PhasePre)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !ok {
		t.Fatal("expected hook to exist")
	}
	if got.Command != h.Command {
		t.Errorf("command: got %q, want %q", got.Command, h.Command)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, ok, err := m.Get("noproject", hook.PhasePost)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected not found")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	h := hook.Hook{Project: "myapp", Phase: hook.PhasePost, Command: "echo post"}
	_ = m.Set(h)
	if err := m.Delete("myapp", hook.PhasePost); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok, _ := m.Get("myapp", hook.PhasePost)
	if ok {
		t.Fatal("expected hook to be deleted")
	}
}

func TestDeleteNoop(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("ghost", hook.PhasePre); err != nil {
		t.Errorf("Delete on missing hook should not error: %v", err)
	}
}

func TestList(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set(hook.Hook{Project: "app", Phase: hook.PhasePre, Command: "make pre"})
	_ = m.Set(hook.Hook{Project: "app", Phase: hook.PhasePost, Command: "make post"})
	hooks, err := m.List("app")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(hooks) != 2 {
		t.Errorf("expected 2 hooks, got %d", len(hooks))
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set(hook.Hook{Project: "", Phase: hook.PhasePre, Command: "echo hi"})
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestSetUnknownPhaseReturnsError(t *testing.T) {
	m := newTempManager(t)
	err := m.Set(hook.Hook{Project: "app", Phase: "during", Command: "echo hi"})
	if err == nil {
		t.Fatal("expected error for unknown phase")
	}
}
