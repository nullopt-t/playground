# Go `sync` package study

One test file per primitive, ordered by learning curve. Every file is
heavily commented with **use case**, **what breaks without it**, and
**side effects / gotchas**.

Run everything with the race detector:

```bash
go test -race -count=1 ./...
```

## Order (read top to bottom)

| # | File | Primitive | When to use |
|---|------|-----------|-------------|
| 01 | `01_mutex_test.go` | `sync.Mutex` | Shared state, exclusive access (reads + writes) |
| 02 | `02_rwmutex_test.go` | `sync.RWMutex` | Reads ≫ writes; concurrent readers, exclusive writer |
| 03 | `03_waitgroup_test.go` | `sync.WaitGroup` | Wait for N goroutines to finish |
| 04 | `04_once_test.go` | `sync.Once` | Run initialization exactly once (lazy singleton) |
| 05 | `05_bank.go` / `05_bank_test.go` | `sync.Mutex` | Real-world example: concurrent bank account |
| 06 | `06_pool_test.go` | `sync.Pool` | Reuse temporary objects to reduce GC |
| 07 | `07_syncmap_test.go` | `sync.Map` | Insert-once / read-many map (specialized) |
| 08 | `08_cond_test.go` | `sync.Cond` | Wait until a condition becomes true |

## Rules of thumb

1. If two goroutines touch the same variable and one writes → you must synchronize.
2. Prefer `map + sync.RWMutex` over `sync.Map` unless you have the read-heavy/stable-key profile.
3. Prefer channels over `sync.Cond` unless you truly need condition-variable semantics.
4. `sync.Pool` objects can be dropped at any time — never keep a pointer after `Put`.
5. `sync.Once` is not reusable; `sync.WaitGroup` is not reusable after counter hits 0.
6. Never copy a `Mutex` or `RWMutex`.
7. Always run concurrent code with `go test -race`.

## Files

- `main.go` — original snippet using `sync.Once` (left untouched; its `x++` is a deliberate race example)
- `go.mod` — module `synk`, Go 1.27.3

