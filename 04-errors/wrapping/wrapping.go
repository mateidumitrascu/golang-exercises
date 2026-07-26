// Package wrapping is the full errors API: sentinels, custom types, %w, Is, As,
// Join, and the two shapes of Unwrap.
//
// The mental model: an error is a tree. A single %w gives you a chain
// (Unwrap() error); errors.Join gives you a node with several children
// (Unwrap() []error). errors.Is and errors.As do a depth-first search over that
// tree.
package wrapping

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrPermission = errors.New("permission denied")
	ErrCorrupt    = errors.New("corrupt data")
)

// PathError is the classic wrapper: what operation, on what, and why.
// Its message is "open /etc/shadow: permission denied" - operation, path,
// then the wrapped error, separated by ": ".
type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string { panic("TODO: implement PathError.Error") }
func (e *PathError) Unwrap() error { panic("TODO: implement PathError.Unwrap") }

// HTTPError carries a status code. It implements a CUSTOM Is so that
// errors.Is(err, HTTPError{Code: 404}) is true for any HTTPError with code 404,
// no matter what message it carries or how deeply it is wrapped.
type HTTPError struct {
	Code int
	Msg  string
}

func (e HTTPError) Error() string { panic("TODO: implement HTTPError.Error") } // "404: not found"

// Is reports whether target is an HTTPError with the same Code.
func (e HTTPError) Is(target error) bool { panic("TODO: implement HTTPError.Is") }

// Open simulates a filesystem:
//
//	"missing"     -> *PathError{Op: "open", Err: ErrNotFound}
//	"secret"      -> *PathError{Op: "open", Err: ErrPermission}
//	"corrupt"     -> *PathError{Op: "read", Err: ErrCorrupt} wrapped once more
//	                 with fmt.Errorf("checksum: %w", ...)
//	anything else -> nil
func Open(path string) error { panic("TODO: implement Open") }

// Chain returns a single error combining all non-nil errors in errs, or nil if
// there are none. One error in, that same error out (not a wrapper around it).
func Chain(errs ...error) error { panic("TODO: implement Chain") }

// RootCause unwraps err as far as it can go and returns the deepest error.
// If it meets a multi-error node (Unwrap() []error), it stops there and returns
// that node - there is no single root below it. RootCause(nil) is nil.
func RootCause(err error) error { panic("TODO: implement RootCause") }

// Flatten returns every error in the tree, depth first, parents before
// children, including err itself. Flatten(nil) is empty.
func Flatten(err error) []error { panic("TODO: implement Flatten") }

// CountLeaves returns the number of errors in the tree that wrap nothing.
func CountLeaves(err error) int { panic("TODO: implement CountLeaves") }

// FirstPathError finds the outermost *PathError in the tree, if any.
func FirstPathError(err error) (*PathError, bool) { panic("TODO: implement FirstPathError") }

// Summarise returns a one-line description of what went wrong, chosen by
// inspecting the tree with errors.Is:
//
//	contains ErrNotFound    -> "missing"
//	contains ErrPermission  -> "denied"
//	contains ErrCorrupt     -> "corrupt"
//	nil                     -> "ok"
//	anything else           -> "unknown"
//
// Check them in that order.
func Summarise(err error) string { panic("TODO: implement Summarise") }
