package timeout_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/envchain-cli/internal/timeout"
)

func newTempManager(t *testing.T) *timeout.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "timeout-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return timeout.New(dir)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	rule := timeout.Rule{Project: "myapp", Duration: 30 * time.Minute, Enabled: true}
	if err := m.Set(rule); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("myapp")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Duration != rule.Duration || !got.Enabled {
		t.Errorf("got %+v, want %+v", got, rule)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err != timeout.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set(timeout.Rule{Project: "proj", Duration: time.Hour, Enabled: true})
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if err != timeout.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("nope"); err != timeout.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestIsDueWhenExpired(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set(timeout.Rule{Project: "svc", Duration: time.Minute, Enabled: true})
	last := time.Now().Add(-2 * time.Minute)
	due, err := m.IsDue("svc", last)
	if err != nil {
		t.Fatalf("IsDue: %v", err)
	}
	if !due {
		t.Error("expected IsDue=true for expired session")
	}
}

func TestIsDueWhenFresh(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set(timeout.Rule{Project: "svc", Duration: time.Hour, Enabled: true})
	due, err := m.IsDue("svc", time.Now())
	if err != nil {
		t.Fatalf("IsDue: %v", err)
	}
	if due {
		t.Error("expected IsDue=false for fresh session")
	}
}

func TestIsDueDisabledRule(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set(timeout.Rule{Project: "svc", Duration: time.Second, Enabled: false})
	due, _ := m.IsDue("svc", time.Now().Add(-time.Hour))
	if due {
		t.Error("expected IsDue=false when rule is disabled")
	}
}

func TestIsDueMissingRule(t *testing.T) {
	m := newTempManager(t)
	due, err := m.IsDue("unknown", time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if due {
		t.Error("expected IsDue=false when no rule exists")
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set(timeout.Rule{Duration: time.Minute, Enabled: true}); err == nil {
		t.Error("expected error for empty project name")
	}
}
