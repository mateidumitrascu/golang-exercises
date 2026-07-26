// Package recovery is about panic, recover, and defer - and about the rules
// that make them less useful than exceptions, on purpose.
//
// Rules to keep in mind:
//   - recover() only does anything when called DIRECTLY by a deferred function
//     of the panicking frame.
//   - a deferred closure can modify NAMED return values; that is how you turn a
//     panic into an error return.
//   - deferred calls run last-in-first-out, and they run even when panicking.
//   - a panic in one goroutine cannot be recovered by another. It kills the
//     whole process. Every goroutine you start needs its own recover if it can
//     panic.
package recovery

// PanicError wraps a recovered panic value.
type PanicError struct {
	Value any    // whatever was passed to panic()
	Stack []byte // the stack trace, captured at recover time
}

// Error is "panic: <value>" using %v for the value.
func (e *PanicError) Error() string { panic("TODO: implement PanicError.Error") }

// Unwrap returns Value if it happens to be an error, so that errors.Is and
// errors.As can see through a panicked error. Otherwise it returns nil.
func (e *PanicError) Unwrap() error { panic("TODO: implement PanicError.Unwrap") }

// Safely runs fn and converts a panic into a *PanicError. It returns nil if fn
// returns normally. If fn panics with nil... note that since Go 1.21 that is
// itself a runtime panic, so you will get a runtime error value.
//
// Capture the stack with runtime.Stack while you are still inside the deferred
// function - by the time you return, the frames are gone.
func Safely(fn func()) (err error) {
	panic("TODO: implement Safely")
}

// SafelyValue is Safely for a function that returns a value. On panic it
// returns the zero value of T and the error.
func SafelyValue[T any](fn func() T) (v T, err error) {
	panic("TODO: implement SafelyValue")
}

// Wrap turns a panicking function into an error-returning one, preserving the
// error fn returns when it does not panic.
func Wrap(fn func() error) (err error) {
	panic("TODO: implement Wrap")
}

// Go runs fn in a new goroutine and delivers either its error or a *PanicError
// on the returned channel, which is closed afterwards. Nothing may escape and
// crash the process.
func Go(fn func() error) <-chan error {
	panic("TODO: implement Go")
}

// CleanupOrder appends the name of each cleanup to a log as it runs, then
// returns the log. Register the cleanups for "a", "b" and "c" in that order
// with defer, panic in the middle, recover, and return the order they actually
// ran in. (Work out the answer before you run the test.)
func CleanupOrder() (order []string) {
	panic("TODO: implement CleanupOrder")
}

// MustPositive panics with a PanicError-friendly value if n <= 0, else returns n.
// Use panic here deliberately: it is a programmer error, not a runtime failure.
func MustPositive(n int) int {
	panic("TODO: implement MustPositive")
}

// Repanic runs fn, and if it panics with a value that is NOT an error, it
// re-panics with the same value (preserving the original stack is impossible,
// but preserving the value is not). If the panic value IS an error, it is
// returned. This models "handle what you understand, let the rest crash".
func Repanic(fn func()) (err error) {
	panic("TODO: implement Repanic")
}
