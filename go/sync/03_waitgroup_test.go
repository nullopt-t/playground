package main

import (
	"sync"
	"testing"
)

// USE CASE: sync.WaitGroup waits for a fixed group of goroutines to finish.
// Use when you fan out work to N goroutines and must block until ALL are done.
//
// WHAT BREAKS WITHOUT IT:
//   main() returns before workers complete. Results are missing, or the
//   program panics (write to closed channel, nil map, etc.) because a worker
//   is still running after main exits.
//
// SIDE EFFECTS / GOTCHAS:
//   - Call Add(n) BEFORE launching the goroutine. Calling Add after `go` can
//     race with Wait() returning early.
//   - Always call defer wg.Done() inside the goroutine; if you forget, Wait()
//     blocks forever (goroutine leak).
//   - WaitGroup is not reusable: after counter reaches 0, do not Add again.

func TestWaitGroupFanout(t *testing.T) {
	urls := []string{"a", "b", "c"}
	results := make(map[string]string, len(urls))
	var mu sync.Mutex       // protects results map from concurrent writes
	var wg sync.WaitGroup   // counts outstanding fetch goroutines

	for _, u := range urls {
		wg.Add(1)           // must happen before go, or Wait may return too early
		go func(url string) {
			defer wg.Done() // signal completion even if fetch panics
			res := "fetched:" + url // pretend network call
			mu.Lock()
			results[url] = res      // concurrent map write needs its own lock
			mu.Unlock()
		}(u)
	}
	wg.Wait()               // blocks until counter hits 0

	if len(results) != len(urls) {
		t.Fatalf("got %d results, want %d", len(results), len(urls))
	}
}
