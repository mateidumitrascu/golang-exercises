// Package nilcheck is about the single most-reported Go "gotcha": an interface
// value holding a nil pointer is not a nil interface.
//
// An interface value is a pair (type, value). It is nil only when BOTH halves
// are nil. Assigning a nil *MyError to an error variable gives you
// (*MyError, nil), which is not nil, which means `if err != nil` fires for a
// function that thought it was reporting success.
package nilcheck

import "errors"

// MyError is a typical custom error type with a pointer receiver.
type MyError struct{ Code int }

func (e *MyError) Error() string { panic("TODO: implement MyError.Error") }

// BuggyValidate is the trap, written out so you can see it. Leave the body of
// this one as it is described: declare `var e *MyError`, set it only when there
// is a real problem, and return it directly. It is SUPPOSED to be broken - a
// test asserts that it is, so you can feel the failure.
func BuggyValidate(n int) error {
	panic("TODO: write the buggy version - declare `var e *MyError`, assign it " +
		"only when n < 0, and `return e` at the end")
}

// Validate is the same logic done right: it must return a truly nil error when
// n >= 0. There is more than one fix; pick one and be able to name the others.
func Validate(n int) error {
	panic("TODO: implement Validate")
}

// IsNil reports whether v is nil in the way a human means it: an untyped nil
// interface, OR an interface holding a nil pointer, map, slice, channel, func
// or unsafe pointer. Everything else is not nil.
//
// You will need reflect for this. reflect.Value.IsNil panics for kinds that
// cannot be nil, so check the Kind first.
func IsNil(v any) bool {
	panic("TODO: implement IsNil")
}

// SafeEqual compares two interface values with ==, without ever panicking.
// Comparing interfaces whose dynamic type is uncomparable (a slice, a map, a
// func) panics at runtime, so guard against it and report comparability.
//
//	SafeEqual(1, 1)               -> true, true
//	SafeEqual([]int{1}, []int{1}) -> false, false
//	SafeEqual(nil, nil)           -> true, true
func SafeEqual(a, b any) (equal, comparable bool) {
	panic("TODO: implement SafeEqual")
}

// FirstNonNil returns the first non-nil pointer in ptrs, or nil.
func FirstNonNil[T any](ptrs ...*T) *T {
	panic("TODO: implement FirstNonNil")
}

// Coalesce returns the first argument that is not the zero value of T, or the
// zero value if there is none.
func Coalesce[T comparable](vals ...T) T {
	panic("TODO: implement Coalesce")
}

// ErrNotFound is a sentinel used by Lookup.
var ErrNotFound = errors.New("not found")

// Lookup returns the value for key, or a wrapped ErrNotFound. The returned
// error must satisfy errors.Is(err, ErrNotFound) and its message must contain
// the key. On success the error must be exactly nil.
func Lookup(m map[string]int, key string) (int, error) {
	panic("TODO: implement Lookup")
}
