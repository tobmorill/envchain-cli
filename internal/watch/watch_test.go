package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/envchain-cli/internal/watch"
)

func writeTempFile(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "store.db")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestWatcherDetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "initial content")

	var calls atomic.Int64
	w := watch.New(path, 50*time.Millisecond, func(p string) {
		if p == path {
			calls.Add(1)
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)

	time.Sleep(120 * time.Millisecond)
	if err := os.WriteFile(path, []byte("changed content"), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}
	time.Sleep(120 * time.Millisecond)

	if calls.Load() == 0 {
		t.Error("expected at least one change notification, got none")
	}
}

func TestWatcherNoSpuriousFire(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "stable")

	var calls atomic.Int64
	w := watch.New(path, 40*time.Millisecond, func(_ string) {
		calls.Add(1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)

	time.Sleep(200 * time.Millisecond)

	if n := calls.Load(); n != 0 {
		t.Errorf("expected 0 spurious calls, got %d", n)
	}
}

func TestWatcherStopsOnContextCancel(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "data")

	var calls atomic.Int64
	w := watch.New(path, 30*time.Millisecond, func(_ string) {
		calls.Add(1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	cancel()

	time.Sleep(100 * time.Millisecond)
	if err := os.WriteFile(path, []byte("after cancel"), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	if n := calls.Load(); n != 0 {
		t.Errorf("watcher fired %d time(s) after cancel", n)
	}
}

func TestWatcherDefaultInterval(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "x")
	// interval <= 0 should not panic
	w := watch.New(path, 0, func(_ string) {})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)
}
