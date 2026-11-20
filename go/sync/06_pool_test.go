package main

import (
	"bytes"
	"sync"
	"testing"
)

// USE CASE: sync.Pool reuses temporary objects to reduce garbage-collection
// pressure. Use for short-lived objects allocated many times (bytes.Buffer,
// JSON encoders, request-scratch buffers).
//
// WHAT BREAKS WITHOUT IT:
//   Allocating a new bytes.Buffer per HTTP request forces GC to reclaim
//   thousands of short-lived objects per second. Under load, GC pause times
//   spike and latency tail worsens.
//
// SIDE EFFECTS / GOTCHAS:
//   - The pool can drop objects at any time (GC may shrink it). Get() may
//     return nil (use New func to avoid).
//   - NEVER keep a pointer to a pooled object after calling Put(); the pool
//     may hand it to another goroutine, causing data races on reused memory.
//   - Not a connection pool. For DB connections use database/sql or a
//     dedicated pool, not sync.Pool.

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) }, // called when pool is empty
}

func TestPoolReuse(t *testing.T) {
	buf := bufPool.Get().(*bytes.Buffer) // borrow from pool (or New)
	buf.Reset()                          // clear previous user's data
	buf.WriteString("hello")             // use it

	if buf.String() != "hello" {
		t.Fatalf("got %q", buf.String())
	}

	bufPool.Put(buf) // return for another goroutine to reuse
}
