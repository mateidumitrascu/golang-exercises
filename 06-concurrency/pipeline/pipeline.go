// Package pipeline is the classic Go concurrency pattern: stages connected by
// channels, each stage a goroutine that owns its output channel.
//
// The contract that makes it work:
//   - The sender closes the channel. Never the receiver.
//   - Every stage must exit when its input closes OR when ctx is cancelled,
//     whichever comes first. A stage blocked on `out <- v` with nobody reading
//     is a leaked goroutine, and it holds everything upstream alive with it.
//   - Leaks are silent. That is why every test here counts goroutines.
package pipeline

import "context"

// Generate sends the values on a channel it owns and closes it. It must return
// as soon as ctx is done, even if nobody is reading.
func Generate[T any](ctx context.Context, vals ...T) <-chan T {
	panic("TODO: implement Generate")
}

// Transform is one stage: read from in, apply f, write to out, close out when
// in closes or ctx is done.
func Transform[T, U any](ctx context.Context, in <-chan T, f func(T) U) <-chan U {
	panic("TODO: implement Transform")
}

// FilterChan drops the values that do not match.
func FilterChan[T any](ctx context.Context, in <-chan T, keep func(T) bool) <-chan T {
	panic("TODO: implement FilterChan")
}

// FanOut starts n goroutines all reading from in, and returns their n output
// channels. This is how you parallelise one slow stage.
func FanOut[T, U any](ctx context.Context, in <-chan T, n int, f func(T) U) []<-chan U {
	panic("TODO: implement FanOut")
}

// FanIn merges several channels into one, which is closed once every input is
// drained or ctx is done. The classic sync.WaitGroup + one goroutine per input.
func FanIn[T any](ctx context.Context, ins ...<-chan T) <-chan T {
	panic("TODO: implement FanIn")
}

// OrDone wraps a channel so that ranging over it stops when ctx is cancelled.
// It exists because `for v := range c` has no way to also select on ctx.Done().
func OrDone[T any](ctx context.Context, in <-chan T) <-chan T {
	panic("TODO: implement OrDone")
}

// Tee duplicates every value to two output channels. Both must be written to,
// so a slow reader on one side slows the other - that is inherent, but neither
// side may deadlock, and both must close when in does.
//
// Hint: inside the loop, use local copies of the two channels and set each to
// nil once you have sent on it. A send on a nil channel blocks forever, which
// removes that case from the select.
func Tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {
	panic("TODO: implement Tee")
}

// Bridge flattens a channel of channels into a single stream, in order.
func Bridge[T any](ctx context.Context, chans <-chan <-chan T) <-chan T {
	panic("TODO: implement Bridge")
}

// Take passes through at most n values, then stops reading (and closes its
// output). Cancelling the pipeline afterwards is the caller's job.
func Take[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	panic("TODO: implement Take")
}
