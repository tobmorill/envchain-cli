// Package watch provides file-system watching for automatic chain reload
// when the underlying store file changes on disk.
package watch

import (
	"context"
	"crypto/sha256"
	"io"
	"os"
	"sync"
	"time"
)

// ChangeFunc is called whenever the watched file changes.
type ChangeFunc func(path string)

// Watcher polls a file for changes and notifies via a callback.
type Watcher struct {
	path     string
	interval time.Duration
	onChange ChangeFunc
	lastHash [sha256.Size]byte
	mu       sync.Mutex
}

// New creates a Watcher for the given file path.
// interval controls how often the file is polled.
func New(path string, interval time.Duration, fn ChangeFunc) *Watcher {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	return &Watcher{path: path, interval: interval, onChange: fn}
}

// Start begins polling in a background goroutine until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) {
	go w.loop(ctx)
}

func (w *Watcher) loop(ctx context.Context) {
	// Seed initial hash so the first tick does not fire a spurious change.
	h, _ := w.hashFile()
	w.mu.Lock()
	w.lastHash = h
	w.mu.Unlock()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.check()
		}
	}
}

func (w *Watcher) check() {
	h, err := w.hashFile()
	if err != nil {
		return
	}
	w.mu.Lock()
	changed := h != w.lastHash
	w.lastHash = h
	w.mu.Unlock()
	if changed {
		w.onChange(w.path)
	}
}

func (w *Watcher) hashFile() ([sha256.Size]byte, error) {
	f, err := os.Open(w.path)
	if err != nil {
		return [sha256.Size]byte{}, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return [sha256.Size]byte{}, err
	}
	var out [sha256.Size]byte
	copy(out[:], h.Sum(nil))
	return out, nil
}
