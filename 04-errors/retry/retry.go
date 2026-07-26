// Package retry: exponential backoff done properly, with error classification
// and context cancellation.
//
// The tests use testing/synctest, so time is fake: a 10-minute backoff runs
// instantly, and the test can still assert exactly how long it "took". Write
// ordinary code with time.NewTimer and it just works.
package retry

import (
	"context"
	"errors"
	"time"
)

// Policy describes how to space out attempts.
type Policy struct {
	MaxAttempts int                               // total attempts including the first; <= 0 means 1
	Base        time.Duration                     // the delay before attempt 2
	Max         time.Duration                     // cap on any single delay; 0 means no cap
	Multiplier  float64                           // growth factor; <= 1 means 2
	Jitter      func(time.Duration) time.Duration // optional, applied last
}

// Delay returns how long to wait after the given attempt number (1 means
// "after the first attempt"). It is Base * Multiplier^(attempt-1), capped at
// Max, then passed through Jitter if set. Delay(n) for n <= 0 is 0.
//
// Watch out for overflow: cap as you go rather than computing a huge number
// first. A duration is an int64 of nanoseconds and it does wrap around.
func (p Policy) Delay(attempt int) time.Duration {
	panic("TODO: implement Policy.Delay")
}

// Permanent marks an error as not worth retrying.
type Permanent struct{ Err error }

func (p *Permanent) Error() string { panic("TODO: implement Permanent.Error") }
func (p *Permanent) Unwrap() error { panic("TODO: implement Permanent.Unwrap") }

// Stop wraps err so that Do gives up immediately. Stop(nil) is nil.
func Stop(err error) error { panic("TODO: implement Stop") }

// Retryabler lets a caller's own error type opt in or out of retrying.
type Retryabler interface{ Retryable() bool }

// IsRetryable decides whether an error is worth another attempt:
//
//	nil                                        -> false
//	anything wrapping a *Permanent             -> false
//	anything implementing Retryabler           -> whatever it says
//	context.Canceled / DeadlineExceeded        -> false
//	anything else                              -> true
//
// Check them in that order, and search the whole error tree (errors.As).
func IsRetryable(err error) bool { panic("TODO: implement IsRetryable") }

// Do calls fn until it succeeds, until the error is not retryable, or until the
// attempts run out - sleeping Policy.Delay(attempt) in between, and aborting the
// sleep if ctx is cancelled.
//
// Return values:
//
//	success                -> nil
//	non-retryable error    -> that error, unwrapped from any *Permanent
//	attempts exhausted     -> an error wrapping the last one, whose message
//	                          contains the number of attempts made
//	ctx cancelled          -> an error for which BOTH errors.Is(err, ctx.Err())
//	                          and errors.Is(err, lastErr) are true
//
// fn receives ctx so it can be cancelled mid-flight.
func Do(ctx context.Context, p Policy, fn func(context.Context) error) error {
	panic("TODO: implement Do")
}

// DoValue is Do for a function that produces a value.
func DoValue[T any](ctx context.Context, p Policy, fn func(context.Context) (T, error)) (T, error) {
	panic("TODO: implement DoValue")
}

// ErrExhausted is returned (wrapped) when the attempts run out. Use it so
// callers can tell "gave up" apart from "failed permanently".
var ErrExhausted = errors.New("retries exhausted")
