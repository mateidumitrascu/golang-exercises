// Package shutdown is about ending things properly: a task group that
// cancels its siblings on the first failure, and an ordered, time-boxed
// shutdown sequence.
//
// You are re-implementing golang.org/x/sync/errgroup. Read its docs afterwards
// and compare - especially the SetLimit semantics.
package shutdown

import (
	"context"
	"time"
)

// Group runs a set of goroutines and collects the first error.
//
// The zero Group is valid and has no limit and no context.
type Group struct {
	// TODO: your fields. You need a WaitGroup, a place for the first error
	// (guarded by a mutex or a sync.Once), an optional cancel func, and an
	// optional limit implemented with a buffered channel.
}

// WithContext returns a Group and a derived context that is cancelled the
// moment the first goroutine returns an error, or when Wait returns.
func WithContext(ctx context.Context) (*Group, context.Context) {
	panic("TODO: implement WithContext")
}

// Go runs f in a new goroutine. If a limit is set, Go BLOCKS until a slot is
// free. A panic inside f must be recovered and turned into the group's error
// rather than killing the process.
func (g *Group) Go(f func() error) { panic("TODO: implement Go") }

// TryGo starts f only if a slot is free, and reports whether it did.
func (g *Group) TryGo(f func() error) bool { panic("TODO: implement TryGo") }

// SetLimit caps the number of goroutines running at once. n < 0 means no limit.
// It panics if called while any goroutine is active.
func (g *Group) SetLimit(n int) { panic("TODO: implement SetLimit") }

// Wait blocks until every goroutine has finished and returns the first non-nil
// error. It also cancels the context from WithContext.
func (g *Group) Wait() error { panic("TODO: implement Wait") }

// Closer is one step of a shutdown sequence.
type Closer struct {
	Name  string
	Close func(context.Context) error
}

// Shutdown runs the closers in REVERSE order - like defer, and for the same
// reason: you bring things down in the opposite order to how you brought them
// up (stop accepting requests, then drain, then close the database).
//
// The whole sequence shares one budget of `timeout`. Each closer gets a context
// carrying the remaining time. If a closer overruns, its error is recorded and
// the sequence continues with whatever budget is left - the rest still deserve
// a chance to flush.
//
// The returned error joins every failure, and each one must name its closer.
// It returns nil if everything closed cleanly.
func Shutdown(ctx context.Context, timeout time.Duration, closers ...Closer) error {
	panic("TODO: implement Shutdown")
}

// Worker is a long-running loop for the drain test below.
type Worker struct {
	// TODO
}

// NewWorker returns a worker that reads jobs from in and calls handle for each.
func NewWorker(in <-chan int, handle func(int)) *Worker { panic("TODO: implement NewWorker") }

// Run processes jobs until ctx is cancelled or in is closed. On cancellation it
// must DRAIN whatever is already buffered in the channel before returning -
// that is what "graceful" means here - and it must return ctx.Err() if it was
// cancelled, or nil if the input simply ended.
func (w *Worker) Run(ctx context.Context) error { panic("TODO: implement Run") }

// Processed reports how many jobs were handled.
func (w *Worker) Processed() int { panic("TODO: implement Processed") }
