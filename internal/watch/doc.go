// Package watch implements lightweight file-system polling for envchain-cli.
//
// It is used to detect when a chain store file has been modified externally
// (e.g. by another process or a sync service) so that cached plaintext can
// be invalidated and the user notified.
//
// Usage:
//
//	w := watch.New("/path/to/store.db", 2*time.Second, func(p string) {
//		log.Printf("store changed: %s", p)
//	})
//	w.Start(ctx)
//
// The watcher runs in a background goroutine and stops automatically when
// the provided context is cancelled.
package watch
