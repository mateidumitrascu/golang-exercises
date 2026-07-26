// Package syncprim rebuilds the synchronisation primitives you normally import,
// so that you know what they cost and what they guarantee.
//
// Run every test in this package with -race. Without it, a broken
// implementation will pass by luck on most runs.
package syncprim

import "sync"

// Once runs a function exactly once, even when several goroutines race to be
// first. Your own sync.Once.
//
// The guarantee that matters and that people get wrong: Do must not return
// until f HAS COMPLETED. A second caller arriving while the first is still
// inside f must block, not skip ahead - otherwise it can observe a half-built
// object. That is why the fast path is an atomic load and the slow path is a
// mutex, not a bare CAS.
//
// A panic inside f still counts as "done".
type Once struct {
	// TODO: your fields. sync/atomic plus a sync.Mutex.
}

func (o *Once) Do(f func()) { panic("TODO: implement Once.Do") }

// Done reports whether f has run.
func (o *Once) Done() bool { panic("TODO: implement Once.Done") }

// Map is a concurrent map guarded by an RWMutex. The zero value is ready to use.
type Map[K comparable, V any] struct {
	// TODO: your fields
	_ sync.RWMutex
}

func (m *Map[K, V]) Get(k K) (V, bool) { panic("TODO: implement Map.Get") }
func (m *Map[K, V]) Set(k K, v V)      { panic("TODO: implement Map.Set") }
func (m *Map[K, V]) Delete(k K)        { panic("TODO: implement Map.Delete") }
func (m *Map[K, V]) Len() int          { panic("TODO: implement Map.Len") }

// GetOrCompute returns the existing value, or computes and stores one.
//
// The subtlety: compute must not run while the write lock is held if it might
// be slow, but two goroutines must not both compute for the same key... which
// is exactly the problem singleflight below solves. For THIS method, take the
// simple road: hold the lock, and document the trade-off in a comment.
// It must call compute at most once per key.
func (m *Map[K, V]) GetOrCompute(k K, compute func() V) (V, bool) {
	panic("TODO: implement Map.GetOrCompute")
}

// Snapshot returns a copy of the contents, taken atomically.
func (m *Map[K, V]) Snapshot() map[K]V { panic("TODO: implement Map.Snapshot") }

// Counter is a set of named counters using atomics under a lock-free read path
// where possible. Keep it simple: a Map[string, *atomic.Int64] is fine, but
// think about the race between two goroutines creating the same counter.
type Counter struct {
	// TODO: your fields
}

func (c *Counter) Inc(name string)          { panic("TODO: implement Counter.Inc") }
func (c *Counter) Add(name string, n int64) { panic("TODO: implement Counter.Add") }
func (c *Counter) Get(name string) int64    { panic("TODO: implement Counter.Get") }

// Max atomically raises the value if n is greater, and returns the new value.
// This is the compare-and-swap loop: load, decide, CAS, retry if someone beat
// you to it. Write it by hand.
func (c *Counter) Max(name string, n int64) int64 { panic("TODO: implement Counter.Max") }

// Group deduplicates concurrent calls for the same key - the singleflight
// pattern. If ten goroutines ask for "user:42" at once and it is not cached,
// exactly one call to fn happens and all ten get its result.
type Group[K comparable, V any] struct {
	// TODO: your fields. A map from key to an in-flight call, plus a mutex.
	// The in-flight record needs a WaitGroup (or a channel) plus space for the
	// result.
}

// Do returns the result of fn for key. shared reports whether the result was
// shared with another caller. fn's panic must not deadlock the other waiters.
func (g *Group[K, V]) Do(key K, fn func() (V, error)) (v V, shared bool, err error) {
	panic("TODO: implement Group.Do")
}

// Forget drops any memory of key so the next Do calls fn again. (Do itself must
// also forget the key once the call completes - this is deduplication of
// in-flight calls, not a cache.)
func (g *Group[K, V]) Forget(key K) { panic("TODO: implement Group.Forget") }
