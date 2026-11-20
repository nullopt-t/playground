package main

import (
	"sync"
	"testing"
)

// USE CASE: sync.RWMutex lets many goroutines READ concurrently while a
// WRITER waits for exclusive access. Use when reads vastly outnumber writes.
//
// WHAT BREAKS WITHOUT IT:
//   - Plain Mutex: every reader serializes; read throughput collapses under load.
//   - Plain map with no lock: concurrent map read/write panics
//     ("fatal error: concurrent map reads and writes").
//
// SIDE EFFECTS / GOTCHAS:
//   - RWMutex is slower than Mutex when uncontended. Don't use it unless reads
//     truly dominate writes.
//   - A writer can starve if readers keep arriving (Go gives writers priority,
//     but new readers may still block the writer).
//   - Never hold a read lock while trying to acquire the write lock -> deadlock.

type ConfigCache struct {
	mu   sync.RWMutex // RLock for readers, Lock for writers
	data map[string]string
}

func NewConfigCache() *ConfigCache {
	return &ConfigCache{data: make(map[string]string)}
}

func (c *ConfigCache) Get(key string) (string, bool) {
	c.mu.RLock()          // many goroutines can be here at once
	defer c.mu.RUnlock()
	v, ok := c.data[key]  // safe read: writer cannot mutate while we hold RLock
	return v, ok
}

func (c *ConfigCache) Set(key, value string) {
	c.mu.Lock()           // exclusive: blocks until all readers release RLock
	defer c.mu.Unlock()
	c.data[key] = value   // safe write: no reader sees a half-written map
}

func TestRWMutexCache(t *testing.T) {
	c := NewConfigCache()
	c.Set("app", "playground")

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if _, ok := c.Get("app"); !ok {
					t.Errorf("missing key")
				}
			}
		}()
	}
	wg.Wait()
}
