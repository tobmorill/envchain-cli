package schedule_test

import (
	"testing"
	"time"

	"github.com/envchain/envchain-cli/internal/schedule"
)

func newTempManager(t *testing.T) *schedule.Manager {
	t.Helper()
	m, err := schedule.New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("myproject", 24*time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Project != "myproject" {
		t.Errorf("project = %q, want %q", rec.Project, "myproject")
	}
	if rec.Interval != 24*time.Hour {
		t.Errorf("interval = %v, want %v", rec.Interval, 24*time.Hour)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != schedule.ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", time.Hour); err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestSetNonPositiveIntervalReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", 0); err == nil {
		t.Fatal("expected error for zero interval")
	}
	if err := m.Set("proj", -time.Second); err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestIsDue(t *testing.T) {
	m := newTempManager(t)
	// Set a very short interval so it becomes due immediately.
	if err := m.Set("proj", time.Nanosecond); err != nil {
		t.Fatalf("Set: %v", err)
	}
	time.Sleep(2 * time.Millisecond)
	rec, err := m.Get("proj")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !rec.IsDue() {
		t.Error("expected IsDue() == true")
	}
}

func TestTouchResetsWindow(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("proj", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := m.Touch("proj"); err != nil {
		t.Fatalf("Touch: %v", err)
	}
	rec, _ := m.Get("proj")
	if rec.IsDue() {
		t.Error("expected IsDue() == false after Touch")
	}
}

func TestDeleteSchedule(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", time.Hour)
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if err != schedule.ErrNotFound {
		t.Errorf("after delete: err = %v, want ErrNotFound", err)
	}
}

func TestSetPreservesCreatedAt(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", time.Hour)
	first, _ := m.Get("proj")
	time.Sleep(2 * time.Millisecond)
	_ = m.Set("proj", 2*time.Hour)
	second, _ := m.Get("proj")
	if !second.CreatedAt.Equal(first.CreatedAt) {
		t.Errorf("CreatedAt changed on overwrite: %v -> %v", first.CreatedAt, second.CreatedAt)
	}
}
