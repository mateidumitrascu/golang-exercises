// Package races is different from the other exercises: the code below is
// FINISHED, and WRONG. Six real concurrency bugs, of the kind that pass code
// review, pass tests on your laptop, and fail in production at 3am.
//
// Your job is to fix each one without changing the exported behaviour.
//
// Run the tests with the race detector - always:
//
//	go test -race ./06-concurrency/races/
//
// Some of these fail as a hard crash ("fatal error: concurrent map writes")
// rather than a test failure, and a crash takes the whole package's test binary
// down with it. Fix them one at a time, top to bottom.
package races

import (
	"errors"
	"sync"
)

// BUG 1: data race on the accumulator.
// Several goroutines read-modify-write `total` with no synchronisation. The
// race detector will point straight at it. There are three good fixes -
// a mutex, an atomic, or per-goroutine partial sums combined at the end.
// The third is usually the fastest. Do that one.
func SumConcurrent(nums []int, workers int) int {
	if workers <= 0 {
		workers = 1
	}
	total := 0
	var wg sync.WaitGroup
	chunk := (len(nums) + workers - 1) / workers
	for i := 0; i < len(nums); i += chunk {
		end := min(i+chunk, len(nums))
		part := nums[i:end]
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, n := range part {
				total += n
			}
		}()
	}
	wg.Wait()
	return total
}

// BUG 2: concurrent map writes.
// Go maps are not safe for concurrent writers; this does not merely race, the
// runtime deliberately crashes the process when it notices. Fix it with a
// mutex, or by having each worker build its own map and merging at the end.
func WordCount(texts []string) map[string]int {
	counts := make(map[string]int)
	var wg sync.WaitGroup
	for _, text := range texts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			word := ""
			for _, r := range text + " " {
				if r == ' ' {
					if word != "" {
						counts[word]++
						word = ""
					}
					continue
				}
				word += string(r)
			}
		}()
	}
	wg.Wait()
	return counts
}

// BUG 3: goroutine leak.
// Only the first result is ever read from an UNBUFFERED channel, so every other
// goroutine blocks forever on its send and is never collected - along with
// everything its closure captured. The fix is one character long.
func FirstError(fns []func() error) error {
	if len(fns) == 0 {
		return nil
	}
	results := make(chan error)
	for _, fn := range fns {
		go func() {
			results <- fn()
		}()
	}
	return <-results
}

// BUG 4: the WaitGroup is used wrongly.
// Add is called inside the goroutine, so Wait can return before any of them has
// started. Fix the ordering.
//
// This one is special: `go vet ./06-concurrency/races/` reports it without
// running anything ("WaitGroup.Add called from inside new goroutine"). Until it
// is fixed, vet fails for this package - which is exactly why vet belongs in
// your pre-commit hook. It runs automatically as part of `go test`, too.
func CountUp(n int) int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0
	for range n {
		go func() {
			wg.Add(1)
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	return count
}

// BUG 5: closing a channel twice panics, and so does sending on a closed one.
// Broadcaster.Close may be called from several goroutines. Make Close
// idempotent and make Send safe to call concurrently with Close - it should
// report false once the broadcaster is closed rather than panicking.
type Broadcaster struct {
	mu     sync.Mutex
	subs   []chan string
	closed bool
}

func (b *Broadcaster) Subscribe() <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 16)
	b.subs = append(b.subs, ch)
	return ch
}

func (b *Broadcaster) Send(msg string) bool {
	for _, ch := range b.subs {
		ch <- msg
	}
	return true
}

func (b *Broadcaster) Close() {
	for _, ch := range b.subs {
		close(ch)
	}
	b.closed = true
}

// BUG 6: check-then-act.
// The lock is released between the lookup and the store, so two goroutines can
// both miss, both compute, and both write - and the expensive `compute` runs
// more than once per key. Nothing here trips the race detector; it is a logic
// race. Fix it so compute runs exactly once per key, and explain to yourself
// why holding the lock across compute is a trade-off rather than a fix.
type Cache struct {
	mu sync.Mutex
	m  map[string]int
	// Computed counts how many times compute actually ran.
	Computed int
}

func NewCache() *Cache { return &Cache{m: make(map[string]int)} }

func (c *Cache) Get(key string, compute func() int) int {
	c.mu.Lock()
	v, ok := c.m[key]
	c.mu.Unlock()
	if ok {
		return v
	}
	v = compute()
	c.mu.Lock()
	c.Computed++
	c.m[key] = v
	c.mu.Unlock()
	return v
}

// ErrNotFixed is unused - delete it once you have been through the file.
var ErrNotFixed = errors.New("not fixed yet")
