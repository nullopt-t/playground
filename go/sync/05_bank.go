package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// USE CASE: sync.Mutex on a struct field. Every method that touches Balance
// locks first, so concurrent deposits and withdrawals are serialized.
//
// WHAT BREAKS WITHOUT IT:
//   Data race on Balance. Deposit reads Balance, adds, writes back.
//   Two concurrent deposits both read the same starting value; one write is
//   lost. Balance ends up wrong, non-deterministically.
//
// SIDE EFFECTS / GOTCHAS:
//   - String() also locks: reading Balance without lock while another
//     goroutine writes is itself a race.
//   - Do not copy Account after first use (copies the mutex state).
//   - Fields ID and Balance are exported but only safe to read when mu is
//     held; prefer methods over direct field access.

type Account struct {
	ID      string
	Balance int64
	mu      sync.Mutex // guards Balance
}

func NewAccount(balance int64) *Account {
	return &Account{
		ID:      fmt.Sprintf("acc-%d", rand.Int()),
		Balance: balance,
	}
}

func (a *Account) Deposit(amount int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance += amount // safe: only one goroutine here at a time
}

func (a *Account) Withdraw(amount int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance -= amount
}

func (a *Account) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return fmt.Sprintf("ID: %s, Balance:%d", a.ID, a.Balance)
}
