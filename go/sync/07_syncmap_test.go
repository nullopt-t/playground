package main

import (
	"sync"
	"testing"
)

// USE CASE: sync.Map is a concurrent map WITHOUT a separate mutex. Use ONLY
// when: keys are inserted once and read many times, OR different goroutines
// touch disjoint keys with little contention.
//
// WHAT BREAKS WITHOUT IT:
//   A plain map with concurrent read/write panics:
//   "fatal error: concurrent map reads and writes".
//   You would normally wrap map + sync.RWMutex; sync.Map replaces that.
//
// SIDE EFFECTS / GOTCHAS:
//   - For most workloads, map + RWMutex is FASTER. sync.Map has higher
//     per-operation overhead. It is a specialized tool, not a general map.
//   - No len() or indexing. Use Load/Store/Delete/Range only.
//   - Do not use for read-heavy/write-heavy mixed workloads; RWMutex wins.

func TestSyncMap(t *testing.T) {
	var m sync.Map
	m.Store("a", 1) // insert once

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				v, ok := m.Load("a") // concurrent reads, no writer -> safe
				if !ok || v != 1 {
					t.Errorf("bad load: %v %v", v, ok)
				}
			}
		}()
	}
	wg.Wait()
}
