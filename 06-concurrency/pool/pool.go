// Package pool: bounded parallelism, ordered results, and first-error
// cancellation - the three things every real "do these N things concurrently"
// helper needs.
//
// The naive version (start a goroutine per item, collect on a channel) breaks
// in production for two reasons: 10,000 items means 10,000 concurrent calls to
// your database, and the results come back in completion order, not input
// order. Fix both.
package pool

import (
	"context"
	"errors"
)

// ParallelMap applies f to every element of in, with at most `workers` calls to
// f running at once, and returns the results IN INPUT ORDER.
//
// Rules:
//   - workers <= 0 panics.
//   - An empty input returns an empty slice, nil error, and starts no goroutines.
//   - The first non-nil error wins: cancel the context passed to the remaining
//     calls to f, stop starting new work, and return (nil, thatError).
//   - If ctx is already cancelled, return ctx.Err() without calling f at all.
//   - f is given a context derived from ctx that is cancelled when the whole
//     operation is done - so a slow f can notice and give up.
//
// Hint: an output slice indexed by position needs no lock. Each worker owns
// out[i] for the i it is working on, and the WaitGroup provides the
// happens-before edge that makes reading them afterwards safe.
func ParallelMap[T, U any](ctx context.Context, in []T, workers int, f func(context.Context, T) (U, error)) ([]U, error) {
	panic("TODO: implement ParallelMap")
}

// ForEach is ParallelMap when you only care about the errors. It must collect
// ALL the errors rather than just the first, joined with errors.Join, and it
// does not cancel on the first failure.
func ForEach[T any](ctx context.Context, in []T, workers int, f func(context.Context, T) error) error {
	panic("TODO: implement ForEach")
}

// ErrPoolClosed is returned by Submit after Close.
var ErrPoolClosed = errors.New("pool: closed")

// Pool is a fixed set of goroutines consuming a task queue. Unlike ParallelMap
// it is long-lived: you create it once and feed it work over time.
//
// Design notes to respect:
//   - New starts exactly `workers` goroutines, no more, ever.
//   - Submit blocks when the queue is full (that is backpressure - it is a
//     feature, not a bug).
//   - Close stops accepting work; Wait blocks until everything already accepted
//     has run. Close must be safe to call more than once and from several
//     goroutines.
//   - A panic in a task must not kill the pool; it is recovered and counted.
type Pool struct {
	// TODO: your fields
}

// New creates a pool with the given number of workers and queue capacity.
// It panics if workers <= 0 or queue < 0.
func New(workers, queue int) *Pool { panic("TODO: implement New") }

// Submit queues a task, blocking if the queue is full. It returns
// ErrPoolClosed if the pool has been closed.
func (p *Pool) Submit(task func()) error { panic("TODO: implement Submit") }

// TrySubmit never blocks: it returns false if the queue is full.
func (p *Pool) TrySubmit(task func()) (bool, error) { panic("TODO: implement TrySubmit") }

// Close stops accepting new work. Idempotent and concurrency-safe.
func (p *Pool) Close() { panic("TODO: implement Close") }

// Wait blocks until all accepted tasks have finished. It must be called after
// Close (calling it before is allowed but it may return early).
func (p *Pool) Wait() { panic("TODO: implement Wait") }

// Panics returns how many tasks panicked.
func (p *Pool) Panics() int { panic("TODO: implement Panics") }
