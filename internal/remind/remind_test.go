package remind_test

import (
	"os"
	"testing"
	"time"

	"github.com/envchain-cli/internal/remind"
)

func newTempManager(t *testing.T) *remind.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "remind-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return remind.New(dir)
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	r := remind.Reminder{
		Project:  "myapp",
		Interval: 24 * time.Hour,
		Message:  "rotate your secrets",
	}
	if err := m.Set(r); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("myapp")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Project != r.Project {
		t.Errorf("project: got %q, want %q", got.Project, r.Project)
	}
	if got.Message != r.Message {
		t.Errorf("message: got %q, want %q", got.Message, r.Message)
	}
	if got.Interval != r.Interval {
		t.Errorf("interval: got %v, want %v", got.Interval, r.Interval)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err != remind.ErrNoReminder {
		t.Fatalf("expected ErrNoReminder, got %v", err)
	}
}

func TestIsDue(t *testing.T) {
	r := remind.Reminder{
		Project:   "myapp",
		Interval:  1 * time.Millisecond,
		LastReset: time.Now().Add(-1 * time.Hour),
	}
	if !r.IsDue() {
		t.Error("expected reminder to be due")
	}
	r.LastReset = time.Now().Add(1 * time.Hour)
	if r.IsDue() {
		t.Error("expected reminder to not be due")
	}
}

func TestReset(t *testing.T) {
	m := newTempManager(t)
	r := remind.Reminder{Project: "myapp", Interval: time.Hour}
	if err := m.Set(r); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := m.Reset("myapp"); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	got, _ := m.Get("myapp")
	if time.Since(got.LastReset) > 2*time.Second {
		t.Error("LastReset not updated by Reset")
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	r := remind.Reminder{Project: "myapp", Interval: time.Hour}
	_ = m.Set(r)
	if err := m.Delete("myapp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("myapp")
	if err != remind.ErrNoReminder {
		t.Errorf("expected ErrNoReminder after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("ghost"); err != remind.ErrNoReminder {
		t.Errorf("expected ErrNoReminder, got %v", err)
	}
}
