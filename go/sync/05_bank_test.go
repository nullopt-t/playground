package main

import (
	"sync"
	"testing"
)

// USE CASE: same Mutex pattern as bank.go, tested under concurrency.
//
// WHAT BREAKS WITHOUT IT:
//   If bank.go's lock were removed, Balance would not return to 100 after
//   20 goroutines each deposit and withdraw 500 times. Lost updates.
//
// SIDE EFFECTS / GOTCHAS:
//   - Add(2) per loop iteration: one deposit goroutine, one withdraw
//     goroutine. Must happen before `go`.
//   - wg.Wait() blocks until all 40 goroutines call Done.

func TestAccountConcurrent(t *testing.T) {
	a := NewAccount(100)
	const workers = 20
	const perWorker = 500

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(2)           // 2 goroutines per iteration
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				a.Deposit(1)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				a.Withdraw(1)
			}
		}()
	}
	wg.Wait()

	if got := a.Balance; got != 100 {
		t.Fatalf("balance = %d, want 100", got)
	}
}
