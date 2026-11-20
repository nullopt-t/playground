package main

import (
	"sync"
	"testing"
)

// USE CASE: sync.Once runs a function exactly once across any number of
// callers. Use for lazy, one-time initialization (config load, DB connection,
// compiled regex, metric registration).
//
// WHAT BREAKS WITHOUT IT:
//   Two goroutines both initialize: two DB connections opened, two config
//   files written, two regexes compiled. Duplicate work, resource leaks, or
//   file corruption.
//
// SIDE EFFECTS / GOTCHAS:
//   - If the function panics, Once considers it "done" and never retries.
//   - Once is NOT reusable: after one Do(), the value is spent.
//   - Do blocks subsequent callers until the first call completes.

type lazy struct {
	value string
	once  sync.Once // guarantees the body runs exactly once
}

func (l *lazy) Get() string {
	l.once.Do(func() {
		l.value = "computed" // runs once, no matter how many goroutines call Get
	})
	return l.value
}

func TestOnce(t *testing.T) {
	var l lazy
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if v := l.Get(); v != "computed" {
				t.Errorf("got %q", v)
			}
		}()
	}
	wg.Wait()
}
