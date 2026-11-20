package main

import (
	"sync"
	"testing"
)

// USE CASE: sync.Mutex protects shared state so only one goroutine
// can read or write it at a time.
// Use when multiple goroutines read AND write the same variable.
//
// WHAT BREAKS WITHOUT IT:
//   Data race. Two goroutines read x, both compute x+1, both write back.
//   One update is lost. The race is non-deterministic: it passes on your
//   laptop and crashes in production. Go calls this undefined behavior.
//
// SIDE EFFECTS / GOTCHAS:
//   - Never copy a Mutex (a copied lock duplicates internal state -> chaos).
//   - Locking the same mutex twice without unlocking = deadlock.
//   - Hold the lock for the shortest time possible; don't wrap whole HTTP
//     handlers or DB calls in one lock.

type Counter struct {
	mu sync.Mutex // guards n
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()          // exclusive access: no other goroutine can be in Inc/Value
	defer c.mu.Unlock()  // always unlock, even if code panics later
	c.n++                // safe: only one goroutine touches n here
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n           // read is synchronized too; without lock, reader sees torn write
}

func TestMutexCounter(t *testing.T) {
	var c Counter
	const workers = 10
	const perWorker = 1000

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				c.Inc() // without Mutex, final value < workers*perWorker (lost updates)
			}
		}()
	}
	wg.Wait()

	if got := c.Value(); got != workers*perWorker {
		t.Fatalf("counter = %d, want %d", got, workers*perWorker)
	}
}
