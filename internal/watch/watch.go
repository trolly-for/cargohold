// Package watch provides file-system watching for secret bundle files,
// notifying callers when a bundle on disk has been modified outside of
// the current process.
package watch

import (
	"os"
	"sync"
	"time"
)

// Event describes a change detected for a watched bundle.
type Event struct {
	Path    string
	Env     string
	ModTime time.Time
}

// Watcher polls a set of bundle paths and emits Events when a file's
// modification time changes.
type Watcher struct {
	interval time.Duration
	stop     chan struct{}
	wg       sync.WaitGroup
}

// New creates a Watcher that checks for changes on the given interval.
func New(interval time.Duration) *Watcher {
	return &Watcher{
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Watch begins polling the file at path, associated with the given env
// label. Events are sent on the returned channel. The channel is closed
// when Stop is called.
func (w *Watcher) Watch(path, env string) <-chan Event {
	ch := make(chan Event, 4)

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer close(ch)

		var lastMod time.Time

		if fi, err := os.Stat(path); err == nil {
			lastMod = fi.ModTime()
		}

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-w.stop:
				return
			case <-ticker.C:
				fi, err := os.Stat(path)
				if err != nil {
					continue
				}
				if fi.ModTime().After(lastMod) {
					lastMod = fi.ModTime()
					ch <- Event{Path: path, Env: env, ModTime: lastMod}
				}
			}
		}
	}()

	return ch
}

// Stop halts all polling goroutines and waits for them to finish.
func (w *Watcher) Stop() {
	close(w.stop)
	w.wg.Wait()
}
