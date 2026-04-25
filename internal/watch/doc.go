// Package watch implements lightweight polling-based file watching for
// cargohold bundle files.
//
// It is intentionally simple: rather than relying on OS-specific
// inotify/kqueue APIs, it uses a configurable polling interval so that
// behaviour is predictable across platforms and inside containers.
//
// Usage:
//
//	w := watch.New(5 * time.Second)
//	events := w.Watch("/home/user/.cargohold/staging.enc", "staging")
//	for ev := range events {
//		log.Printf("bundle %s changed at %s", ev.Env, ev.ModTime)
//	}
package watch
