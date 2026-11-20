package main

import (
	"sync"
	"testing"
)

// USE CASE: sync.Cond wakes goroutines when a condition becomes true.
// Use for producer/consumer queues, batch flushes, or any "wait until X"
// scenario where channels are awkward.
//
// WHAT BREAKS WITHOUT IT:
//   Busy-looping (spin-waiting) on a flag burns CPU. Or a goroutine blocks
//   forever because nobody signaled it that the condition changed.
//
// SIDE EFFECTS / GOTCHAS:
//   - Spurious wakeups happen: always Wait inside a `for` loop checking the
//     condition, never `if`.
//   - Easy to misuse (lost wakeups, deadlocks). In idiomatic Go, prefer
//     channels or context for signaling. Use Cond only when you truly need
//     condition-variable semantics.
//   - The Lock passed to NewCond must be held when calling Wait/Signal.

type Queue struct {
	mu    sync.Mutex
	cond  *sync.Cond   // tied to mu; Wait releases mu, sleeps, re-acquires on wake
	items []int
}

func NewQueue() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu) // associate cond with this mutex
	return q
}

func (q *Queue) Push(n int) {
	q.mu.Lock()
	q.items = append(q.items, n)
	q.mu.Unlock()
	q.cond.Signal() // wake ONE waiting Pop
}

func (q *Queue) Pop() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 { // `for`, not `if`: guards against spurious wakeups
		q.cond.Wait()     // sleeps; woken by Push's Signal
	}
	n := q.items[0]
	q.items = q.items[1:]
	return n
}

func TestCondQueue(t *testing.T) {
	q := NewQueue()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if n := q.Pop(); n != 42 { // blocks until Push runs
			t.Errorf("got %d, want 42", n)
		}
	}()
	q.Push(42)
	wg.Wait()
}
